package commands

import (
	"context"

	"delivery/internal/pkg/errs"
)

type CreateOrderCommandHandler interface {
	Handle(context.Context, CreateOrderCommand) error
}

type createOrderCommandHandler struct{}

func NewCreateOrderCommandHandler() (*createOrderCommandHandler, error) {
	return &createOrderCommandHandler{}, nil
}

func (h *createOrderCommandHandler) Handle(_ context.Context, command CreateOrderCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("command")
	}
	return nil
}
