package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"delivery/cmd"
	httpin "delivery/internal/adapters/in/http"
	"delivery/internal/adapters/out/postgres/courierrepo"
	"delivery/internal/adapters/out/postgres/orderrepo"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/core/domain/services"
	"delivery/internal/generated/servers"
	"delivery/internal/pkg/errs"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	oam "github.com/oapi-codegen/echo-middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var loadOnce sync.Once

func main() {
	cfg := getConfigs()

	dsn, err := makeConnectionString(cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbSslMode)
	if err != nil {
		log.Fatal(err.Error())
	}

	crateDbIfNotExists(cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbSslMode)
	db := mustGormOpen(dsn)
	mustAutoMigrate(db)

	cr := cmd.NewCompositionRoot(cfg, db)
	startJobs(cr, context.TODO())
	startWebServer(cr, cfg.HttpPort)
}

func getConfigs() cmd.Config {
	config := cmd.Config{
		HttpPort:                  goDotEnvVariable("HTTP_PORT"),
		DbHost:                    goDotEnvVariable("DB_HOST"),
		DbPort:                    goDotEnvVariable("DB_PORT"),
		DbUser:                    goDotEnvVariable("DB_USER"),
		DbPassword:                goDotEnvVariable("DB_PASSWORD"),
		DbName:                    goDotEnvVariable("DB_NAME"),
		DbSslMode:                 goDotEnvVariable("DB_SSLMODE"),
		GeoServiceGrpcHost:        goDotEnvVariable("GEO_SERVICE_GRPC_HOST"),
		KafkaHost:                 goDotEnvVariable("KAFKA_HOST"),
		KafkaConsumerGroup:        goDotEnvVariable("KAFKA_CONSUMER_GROUP"),
		KafkaBasketConfirmedTopic: goDotEnvVariable("KAFKA_BASKET_CONFIRMED_TOPIC"),
		KafkaOrderChangedTopic:    goDotEnvVariable("KAFKA_ORDER_CHANGED_TOPIC"),
	}
	return config
}

func goDotEnvVariable(key string) string {
	loadOnce.Do(func() {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatalf("ERROR: loading .env file")
		}
	})

	return os.Getenv(key)
}

func startWebServer(compositionRoot *cmd.CompositionRoot, port string) {
	handlers, err := httpin.New(
		compositionRoot.NewCreateOrderCommandHandler(),
		compositionRoot.NewCreateCourierCommandHandler(),
		compositionRoot.NewGetAllCouriersQueryHandler(),
		compositionRoot.NewGetIncompletedOrdersQueryHandler(),
	)
	if err != nil {
		log.Fatalf("ERROR: init HTTP Server: %v", err)
	}

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
	}))

	spec, err := servers.GetSwagger()
	if err != nil {
		log.Fatalf("ERROR: reading OpenAPI spec: %v", err)
	}
	_ = oam.OapiRequestValidator(spec)
	// e.Use(oam.OapiRequestValidator(spec)) // bug: it will break //docs ang /openapi.json
	e.Pre(middleware.RemoveTrailingSlash())
	registerSwaggerOpenApi(e)
	registerSwaggerUi(e)
	servers.RegisterHandlers(e, handlers)
	e.Logger.Fatal(e.Start("0.0.0.0:" + port))
}

func registerSwaggerOpenApi(e *echo.Echo) {
	e.GET("/openapi.json", func(c echo.Context) error {
		swagger, err := servers.GetSwagger()
		if err != nil {
			return c.String(http.StatusInternalServerError, "failed to load swagger: "+err.Error())
		}

		data, err := swagger.MarshalJSON()
		if err != nil {
			return c.String(http.StatusInternalServerError, "failed to marshal swagger: "+err.Error())
		}

		return c.Blob(http.StatusOK, "application/json", data)
	})
}

func registerSwaggerUi(e *echo.Echo) {
	e.GET("/docs", func(c echo.Context) error {
		html := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
		  <meta charset="UTF-8">
		  <title>Swagger UI</title>
		  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css">
		</head>
		<body>
		  <div id="swagger-ui"></div>
		  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
		  <script>
			window.onload = () => {
			  SwaggerUIBundle({
				url: "/openapi.json",
				dom_id: "#swagger-ui",
			  });
			};
		  </script>
		</body>
		</html>`
		return c.HTML(http.StatusOK, html)
	})
}

func makeConnectionString(host string, port string, user string,
	password string, dbName string, sslMode string,
) (string, error) {
	if host == "" {
		return "", errs.NewValueIsRequiredError(host)
	}
	if port == "" {
		return "", errs.NewValueIsRequiredError(port)
	}
	if user == "" {
		return "", errs.NewValueIsRequiredError(user)
	}
	if password == "" {
		return "", errs.NewValueIsRequiredError(password)
	}
	if dbName == "" {
		return "", errs.NewValueIsRequiredError(dbName)
	}
	if sslMode == "" {
		return "", errs.NewValueIsRequiredError(sslMode)
	}
	return fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		host, port, user, password, dbName, sslMode), nil
}

func crateDbIfNotExists(host string, port string, user string,
	password string, dbName string, sslMode string,
) {
	dsn, err := makeConnectionString(host, port, user, password, "postgres", sslMode)
	if err != nil {
		log.Fatalf("ERROR: make DSN for PostgreSQL: %v", err)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("ERROR: connect to PostgreSQL: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("ERROR: close db: %v", err)
		}
	}()

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		log.Printf("WARNING: cannot create DB (may be DB already exists): %v", err)
	}
}

func mustGormOpen(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.New(
		postgres.Config{DSN: dsn, PreferSimpleProtocol: true},
	), &gorm.Config{})
	if err != nil {
		log.Fatalf("ERROR: connect to postgres through gorm\n: %s", err)
	}
	return db
}

func mustAutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&courierrepo.CourierDTO{})
	if err != nil {
		log.Fatalf("ERROR: automigrate courier %v", err)
	}

	err = db.AutoMigrate(&courierrepo.StoragePlaceDTO{})
	if err != nil {
		log.Fatalf("ERROR: automigrate storage place:%v", err)
	}

	err = db.AutoMigrate(&orderrepo.OrderDTO{})
	if err != nil {
		log.Fatalf("ERROR: automigrate order: %v", err)
	}
}

func startJobs(cr *cmd.CompositionRoot, ctx context.Context) {
	assignOrdersCommandHandler, err := commands.NewAssignOrderCommandHandler(
		cr.NewUnitOfWorkFactory(),
		services.NewOrderDispatcher(),
	)
	if err != nil {
		log.Fatalf("ERROR: create assignOrdersCommandHandler: %v", err)
	}

	moveCouriersCommandHandler, err := commands.NewMoveCouriersCommandHandler(cr.NewUnitOfWorkFactory())
	if err != nil {
		log.Fatalf("ERROR: create assignOrdersCommandHandler: %v", err)
	}

	ch := time.Tick(time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ch:
				// assign orders
				assignOrderCommand, err := commands.NewAssignOrderCommand()
				if err != nil {
					log.Printf("ERROR: make command AssignOrderCommand: %v", err)
				}
				err = assignOrdersCommandHandler.Handle(ctx, assignOrderCommand)
				if err != nil {
					log.Printf("ERROR: handle command AssignOrderCommand: %v", err)
				}
				// move couriers
				moveCouriersCommand, err := commands.NewMoveCouriersCommand()
				if err != nil {
					log.Printf("ERROR: make command MoveCouriersCommand: %v", err)
				}
				err = moveCouriersCommandHandler.Handle(ctx, moveCouriersCommand)
				if err != nil {
					log.Printf("ERROR: handle command MoveCouriersCommand: %v", err)
				}
			}
		}
	}()
}
