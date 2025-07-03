package commands

import (
	"context"

	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type AddStoragePlaceCommandHandler interface {
	Handle(context.Context, AddStoragePlaceCommand) error
}

type addStoragePlaceCommandHandler struct {
	factory ports.UnitOfWorkFactory
}

func NewAddStoragePlaceCommandHandler(factory ports.UnitOfWorkFactory) (*addStoragePlaceCommandHandler, error) {
	if factory == nil {
		return nil, errs.NewValueIsRequiredError("factory")
	}
	return &addStoragePlaceCommandHandler{factory: factory}, nil
}

func (h *addStoragePlaceCommandHandler) Handle(ctx context.Context, command AddStoragePlaceCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("command")
	}

	uow, err := h.factory()
	if err != nil {
		return err
	}
	defer uow.Rollback(ctx)

	courier, err := uow.CourierRepository().Get(ctx, command.CourierID())
	if err != nil {
		return err
	}

	if err = courier.AddStoragePlace(command.Name(), command.TotalVolume()); err != nil {
		return err
	}

	if err = uow.CourierRepository().Update(ctx, courier); err != nil {
		return err
	}

	return uow.Commit(ctx)
}
