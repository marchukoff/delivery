package commands

import (
	"context"

	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type CreateCourierCommandHandler interface {
	Handle(context.Context, CreateCourierCommand) error
}

type createCourierCommandHandler struct {
	factory ports.UnitOfWorkFactory
}

func NewCreateCourierCommandHandler(factory ports.UnitOfWorkFactory) (*createCourierCommandHandler, error) {
	if factory == nil {
		return nil, errs.NewValueIsRequiredError("factory")
	}
	return &createCourierCommandHandler{factory: factory}, nil
}

func (h *createCourierCommandHandler) Handle(ctx context.Context, command CreateCourierCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("command")
	}

	uow, err := h.factory()
	if err != nil {
		return err
	}
	defer uow.Rollback(ctx)

	courier, err := courier.NewCourier(command.Name(), command.Speed(), kernel.NewRandomLocation())
	if err != nil {
		return err
	}

	// Сохранили

	if err = uow.CourierRepository().Add(ctx, courier); err != nil {
		return err
	}

	return uow.Commit(ctx)
}
