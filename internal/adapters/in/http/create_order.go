package http

import (
	"errors"
	"net/http"

	"delivery/internal/adapters/in/http/problems"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (s *Server) CreateOrder(c echo.Context) error {
	cmd, err := commands.NewCreateOrderCommand(uuid.New(), "TestOrder", 5)
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	err = s.createOrder.Handle(c.Request().Context(), cmd)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return problems.NewNotFound(err.Error())
		}
		return problems.NewConflict(err.Error(), "/")
	}

	return c.JSON(http.StatusOK, nil)
}
