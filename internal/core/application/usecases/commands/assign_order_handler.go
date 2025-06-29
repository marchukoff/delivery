package commands

import (
	"context"

	"delivery/internal/pkg/errs"
)

type AssignOrderCommandCommandHandler interface {
	Handle(context.Context, AssignOrderCommand) error
}

type assignOrderCommandHandler struct{}

func NewAssignOrderCommandHandler() (*assignOrderCommandHandler, error) {
	return &assignOrderCommandHandler{}, nil
}

func (h *assignOrderCommandHandler) Handle(_ context.Context, command AssignOrderCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("command")
	}
	return nil
}
