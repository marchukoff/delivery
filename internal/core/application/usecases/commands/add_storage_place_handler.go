package commands

import (
	"context"
	"delivery/internal/pkg/errs"
)

type AddStoragePlaceCommandHandler interface {
	Handle(context.Context, AddStoragePlaceCommand) error
}

type addStoragePlaceCommandHandler struct{}

func NewAddStoragePlaceCommandHandler() (*addStoragePlaceCommandHandler, error) {
	return &addStoragePlaceCommandHandler{}, nil
}

func (h *addStoragePlaceCommandHandler) Handle(_ context.Context, command AddStoragePlaceCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("command")
	}
	return nil
}
