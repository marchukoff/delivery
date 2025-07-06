package http

import (
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/generated/servers"
	"delivery/internal/pkg/errs"
	"github.com/labstack/echo/v4"
)

var _ servers.ServerInterface = (*Server)(nil)

type Server struct {
	createOrder          commands.CreateOrderCommandHandler
	createCourier        commands.CreateCourierCommandHandler
	getAllCouriers       queries.GetAllCouriersQueryHandler
	getIncompletedOrders queries.GetIncompleteOrdersQueryHandler
}

// CreateCourier implements servers.ServerInterface.
func (s *Server) CreateCourier(ctx echo.Context) error {
	panic("unimplemented")
}

// CreateOrder implements servers.ServerInterface.
func (s *Server) CreateOrder(ctx echo.Context) error {
	panic("unimplemented")
}

// GetCouriers implements servers.ServerInterface.
func (s *Server) GetCouriers(ctx echo.Context) error {
	panic("unimplemented")
}

// GetOrders implements servers.ServerInterface.
func (s *Server) GetOrders(ctx echo.Context) error {
	panic("unimplemented")
}

func New(
	createOrder commands.CreateOrderCommandHandler,
	createCourier commands.CreateCourierCommandHandler,
	getAllCouriers queries.GetAllCouriersQueryHandler,
	getIncompletedOrders queries.GetIncompleteOrdersQueryHandler,
) (*Server, error) {
	if createOrder == nil {
		return nil, errs.NewValueIsRequiredError("createOrder")
	}

	if createCourier == nil {
		return nil, errs.NewValueIsRequiredError("createCourier")
	}

	if getAllCouriers == nil {
		return nil, errs.NewValueIsRequiredError("getAllCouriers")
	}

	if getIncompletedOrders == nil {
		return nil, errs.NewValueIsRequiredError("getIncompletedOrders")
	}

	return &Server{
		createOrder:          createOrder,
		createCourier:        createCourier,
		getAllCouriers:       getAllCouriers,
		getIncompletedOrders: getIncompletedOrders,
	}, nil
}
