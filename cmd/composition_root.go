package cmd

import (
	"io"
	"log"
	"sync"

	grpcout "delivery/internal/adapters/out/grpc"
	"delivery/internal/adapters/out/postgres"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"

	"gorm.io/gorm"
)

type CompositionRoot struct {
	config    Config
	db        *gorm.DB
	geoClient ports.GeoClient
	onceGeo   sync.Once
	//
	closers []io.Closer
}

func NewCompositionRoot(cfg Config, db *gorm.DB) *CompositionRoot {
	return &CompositionRoot{config: cfg, db: db}
}

func (c *CompositionRoot) NewOrderDispatcherService() services.OrderDispatcher {
	return services.NewOrderDispatcher()
}

func (c *CompositionRoot) NewUnitOfWorkFactory() ports.UnitOfWorkFactory {
	factory, err := postgres.NewUnitOfWorkFactory(c.db)
	if err != nil {
		log.Fatalf("new unit of work factory: %v", err)
	}
	return factory
}

func (c *CompositionRoot) NewCreateOrderCommandHandler() commands.CreateOrderCommandHandler {
	h, err := commands.NewCreateOrderCommandHandler(c.NewUnitOfWorkFactory(), c.NewGeoClient())
	if err != nil {
		log.Fatalf("ERROR: cannot create CreateOrderCommandHandler: %v", err)
	}
	return h
}

func (c *CompositionRoot) NewCreateCourierCommandHandler() commands.CreateCourierCommandHandler {
	h, err := commands.NewCreateCourierCommandHandler(c.NewUnitOfWorkFactory())
	if err != nil {
		log.Fatalf("ERROR: cannot create CreateCourierCommandHandler: %v", err)
	}
	return h
}

func (c *CompositionRoot) NewGetAllCouriersQueryHandler() queries.GetAllCouriersQueryHandler {
	h, err := queries.NewGetAllCouriersQueryHandler(c.db)
	if err != nil {
		log.Fatalf("ERROR: cannot create GetAllCouriersQueryHandler: %v", err)
	}
	return h
}

func (c *CompositionRoot) NewGetIncompletedOrdersQueryHandler() queries.GetIncompleteOrdersQueryHandler {
	h, err := queries.NewGetIncompleteOrdersHandler(c.db)
	if err != nil {
		log.Fatalf("ERROR: cannot create GetIncompletedOrdersQueryHandler: %v", err)
	}
	return h
}

func (cr *CompositionRoot) NewGeoClient() ports.GeoClient {
	cr.onceGeo.Do(func() {
		client, err := grpcout.NewClient(cr.config.GeoServiceGrpcHost)
		if err != nil {
			log.Fatalf("ERROR: create GeoClient: %v", err)
		}
		cr.RegisterCloser(client)
		cr.geoClient = client
	})
	return cr.geoClient
}

func (cr *CompositionRoot) RegisterCloser(c io.Closer) {
	cr.closers = append(cr.closers, c)
}

func (cr *CompositionRoot) CloseAll() {
	for _, closer := range cr.closers {
		if err := closer.Close(); err != nil {
			log.Printf("ERROR: closing resource: %v", err)
		}
	}
}
