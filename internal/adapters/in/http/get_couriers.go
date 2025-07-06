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

func (s *Server) GetCouriers(c echo.Context) error {
	query, err := queries.NewGetAllCouriersQuery()
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	queryResponse, err := s.getAllCouriers.Handle(c.Request().Context(), query)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return problems.NewNotFound(err.Error())
		}
	}

	httpResponse := make([]servers.Courier, 0, len(queryResponse.Couriers))
	for _, courier := range queryResponse.Couriers {
		location := servers.Location{
			X: courier.Location.X,
			Y: courier.Location.Y,
		}

		courier := servers.Courier{
			Id:       courier.ID,
			Name:     courier.Name,
			Location: location,
		}
		httpResponse = append(httpResponse, courier)
	}
	return c.JSON(http.StatusOK, httpResponse)
}
