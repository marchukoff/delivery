package commands

import (
	"context"

	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type AssignOrderCommandCommandHandler interface {
	Handle(context.Context, AssignOrderCommand) error
}

type assignOrderCommandHandler struct {
	factory    ports.UnitOfWorkFactory
	dispatcher services.OrderDispatcher
}

func NewAssignOrderCommandHandler(factory ports.UnitOfWorkFactory, dispatcher services.OrderDispatcher) (*assignOrderCommandHandler, error) {
	if factory == nil {
		return nil, errs.NewValueIsRequiredError("factory")
	}

	if dispatcher == nil {
		return nil, errs.NewValueIsRequiredError("dispatcher")
	}

	return &assignOrderCommandHandler{factory: factory, dispatcher: dispatcher}, nil
}

func (h *assignOrderCommandHandler) Handle(ctx context.Context, command AssignOrderCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("command")
	}

	uow, err := h.factory()
	if err != nil {
		return err
	}
	defer uow.Rollback(ctx)

	order, err := uow.OrderRepository().GetFirstInCreatedStatus(ctx)
	if err != nil {
		return err
	}

	couriers, err := uow.CourierRepository().GetAllFree(ctx)
	if err != nil {
		return err
	}

	courier, err := h.dispatcher.Dispatch(order, couriers)
	if err != nil {
		return err
	}

	if err = uow.OrderRepository().Update(ctx, order); err != nil {
		return err
	}

	if err = uow.CourierRepository().Update(ctx, courier); err != nil {
		return err
	}

	return uow.Commit(ctx)
}
