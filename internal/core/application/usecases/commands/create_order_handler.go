package commands

import (
	"context"

	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type CreateOrderCommandHandler interface {
	Handle(context.Context, CreateOrderCommand) error
}

type createOrderCommandHandler struct {
	factory ports.UnitOfWorkFactory
}

func NewCreateOrderCommandHandler(factory ports.UnitOfWorkFactory) (*createOrderCommandHandler, error) {
	if factory == nil {
		return nil, errs.NewValueIsRequiredError("factory")
	}
	return &createOrderCommandHandler{factory: factory}, nil
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

	order, err := order.NewOrder(command.OrderID(), kernel.NewRandomLocation(), command.Volume())
	if err != nil {
		return err
	}

	if err = uow.OrderRepository().Add(ctx, order); err != nil {
		return err
	}

	return uow.Commit(ctx)
}
