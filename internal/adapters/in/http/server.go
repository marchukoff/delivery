package http

import (
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/generated/servers"
	"delivery/internal/pkg/errs"
)

var _ servers.ServerInterface = (*Server)(nil)

type Server struct {
	createOrder          commands.CreateOrderCommandHandler
	createCourier        commands.CreateCourierCommandHandler
	getAllCouriers       queries.GetAllCouriersQueryHandler
	getIncompletedOrders queries.GetIncompleteOrdersQueryHandler
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
