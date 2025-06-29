package queries

import "context"

type GetAllCouriersQueryHandler interface {
	Handle(context.Context, GetAllCouriersQuery) (GetAllCouriersResponse, error)
}

func NewGetAllCouriersQueryHandler() (*getAllCouriersQueryHandler, error) {
	return &getAllCouriersQueryHandler{}, nil
}

type getAllCouriersQueryHandler struct{}

func (h *getAllCouriersQueryHandler) Handle(context.Context, GetAllCouriersQuery) (GetAllCouriersResponse, error) {
	return GetAllCouriersResponse{}, nil
}
