package commands

import (
	"context"

	"delivery/internal/pkg/errs"
)

type MoveCouriersCommandHandler interface {
	Handle(context.Context, MoveCouriersCommand) error
}

type moveCouriersCommandHandler struct{}

func NewMoveCouriersCommandHandler() (*moveCouriersCommandHandler, error) {
	return &moveCouriersCommandHandler{}, nil
}

func (h *moveCouriersCommandHandler) Handle(_ context.Context, command MoveCouriersCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("command")
	}
	return nil
}
