package cmd

import (
	"log"

	"delivery/internal/adapters/out/postgres"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"

	"gorm.io/gorm"
)

type CompositionRoot struct {
	cfg Config
	db  *gorm.DB
}

func NewCompositionRoot(cfg Config, db *gorm.DB) *CompositionRoot {
	return &CompositionRoot{cfg: cfg, db: db}
}

func (c *CompositionRoot) NewOrderDispatcherService() services.OrderDispatcher {
	return services.NewOrderDispatcher()
}

func (c *CompositionRoot) NewUnitOfWork() ports.UnitOfWorkFactory {
	return postgres.NewUnitOfWorkFactory((*gorm.DB)(nil))
}

func (c *CompositionRoot) NewCreateOrderCommandHandler() commands.CreateOrderCommandHandler {
	h, err := commands.NewCreateOrderCommandHandler(c.NewUnitOfWorkFactory())
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

func (c *CompositionRoot) NewUnitOfWorkFactory() ports.UnitOfWorkFactory {
	return postgres.NewUnitOfWorkFactory(c.db)
}
