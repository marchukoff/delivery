package commands

import (
	"context"
	"errors"

	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type MoveCouriersCommandHandler interface {
	Handle(context.Context, MoveCouriersCommand) error
}

type moveCouriersCommandHandler struct {
	factory ports.UnitOfWorkFactory
}

func NewMoveCouriersCommandHandler(factory ports.UnitOfWorkFactory) (*moveCouriersCommandHandler, error) {
	if factory == nil {
		errs.NewValueIsRequiredError("uow")
	}
	return &moveCouriersCommandHandler{factory: factory}, nil
}

func (h *moveCouriersCommandHandler) Handle(ctx context.Context, command MoveCouriersCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("command")
	}

	uow, err := h.factory()
	if err != nil {
		return err
	}
	defer uow.Rollback(ctx)

	orders, err := uow.OrderRepository().GetAllInAssignedStatus(ctx)
	if err != nil {
		return err
	}

	for _, order := range orders {
		courier, err := uow.CourierRepository().Get(ctx, *order.CourierID())
		if err != nil {
			if errors.Is(err, errs.ErrObjectNotFound) {
				return nil
			}
			return err
		}

		if err = courier.Move(order.Location()); err != nil {
			return err
		}

		if courier.Location().Equals(order.Location()) {
			if err = order.Complete(); err != nil {
				return err
			}

			if err = courier.CompleteOrder(order); err != nil {
				return err
			}
		}

		if err = uow.OrderRepository().Update(ctx, order); err != nil {
			return err
		}

		if err = uow.CourierRepository().Update(ctx, courier); err != nil {
			return err
		}
	}

	return uow.Commit(ctx)
}
