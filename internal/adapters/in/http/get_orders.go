package http

import (
	"errors"
	"net/http"

	"delivery/internal/adapters/in/http/problems"
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/generated/servers"
	"delivery/internal/pkg/errs"

	"github.com/labstack/echo/v4"
)

func (s *Server) GetOrders(c echo.Context) error {
	query, err := queries.NewGetIncompleteOrdersQuery()
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	queryResponse, err := s.getIncompletedOrders.Handle(c.Request().Context(), query)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return problems.NewNotFound(err.Error())
		}
	}

	httpResponse := make([]servers.Order, 0, len(queryResponse.Orders))
	for _, courier := range queryResponse.Orders {
		location := servers.Location{
			X: courier.Location.X,
			Y: courier.Location.Y,
		}

		courier := servers.Order{
			Id:       courier.ID,
			Location: location,
		}
		httpResponse = append(httpResponse, courier)
	}
	return c.JSON(http.StatusOK, httpResponse)
}
