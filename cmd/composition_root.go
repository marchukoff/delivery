package cmd

import "delivery/internal/core/domain/services"

type CompositionRoot struct{}

func NewCompositionRoot(_ Config) CompositionRoot {
	app := CompositionRoot{}
	return app
}

func (c *CompositionRoot) NewOrderDispatcherService() services.OrderDispatcher {
	return services.NewOrderDispatcher()
}
