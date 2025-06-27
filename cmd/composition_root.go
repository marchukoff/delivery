package cmd

import (
	"delivery/internal/adapters/out/postgres"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"
	"log"

	"gorm.io/gorm"
)

type CompositionRoot struct{}

func NewCompositionRoot(_ Config) CompositionRoot {
	app := CompositionRoot{}
	return app
}

func (c *CompositionRoot) NewOrderDispatcherService() services.OrderDispatcher {
	return services.NewOrderDispatcher()
}

func (cr *CompositionRoot) NewUnitOfWork() ports.UnitOfWork {
	unitOfWork, err := postgres.NewUnitOfWork((*gorm.DB)(nil))
	if err != nil {
		log.Fatalf("cannot create UnitOfWork: %v", err)
	}
	return unitOfWork
}
