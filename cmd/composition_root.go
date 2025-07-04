package cmd

import (
	"delivery/internal/adapters/out/postgres"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"

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

func (cr *CompositionRoot) NewUnitOfWork() ports.UnitOfWorkFactory {
	return postgres.NewUnitOfWorkFactory((*gorm.DB)(nil))
}
