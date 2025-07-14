package commands

import (
	"context"

	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type CreateOrderCommandHandler interface {
	Handle(context.Context, CreateOrderCommand) error
}

type createOrderCommandHandler struct {
	factory   ports.UnitOfWorkFactory
	geoClient ports.GeoClient
}

func NewCreateOrderCommandHandler(factory ports.UnitOfWorkFactory, geoClient ports.GeoClient) (*createOrderCommandHandler, error) {
	if factory == nil {
		return nil, errs.NewValueIsRequiredError("factory")
	}
	if geoClient == nil {
		return nil, errs.NewValueIsRequiredError("geoClient")
	}
	return &createOrderCommandHandler{factory: factory, geoClient: geoClient}, nil
}

func (h *createOrderCommandHandler) Handle(ctx context.Context, command CreateOrderCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("command")
	}

	uow, err := h.factory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.RollbackUnlessCommitted(ctx)

	location, err := h.geoClient.GetGeolocation(ctx, command.Street())
	if err != nil {
		return err
	}

	order, err := order.NewOrder(command.OrderID(), location, command.Volume())
	if err != nil {
		return err
	}

	if err = uow.OrderRepository().Add(ctx, order); err != nil {
		return err
	}

	return uow.Commit(ctx)
}
