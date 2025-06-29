package commands

import (
	"context"

	"delivery/internal/pkg/errs"
)

type CreateCourierCommandHandler interface {
	Handle(context.Context, CreateCourierCommand) error
}

type createCourierCommandHandler struct{}

func NewAddCourierCommandHandler() (*createCourierCommandHandler, error) {
	return &createCourierCommandHandler{}, nil
}

func (h *createCourierCommandHandler) Handle(_ context.Context, command CreateCourierCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("command")
	}
	return nil
}
