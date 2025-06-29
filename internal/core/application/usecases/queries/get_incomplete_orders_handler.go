package queries

import "context"

type GetIncompleteOrdersQueryHandler interface {
	Handle(context.Context, GetIncompleteOrdersQuery) (GetIncompleteOrdersResponse, error)
}

func NewGetIncompleteOrdersHandler() (*getIncompleteOrdersHandler, error) {
	return &getIncompleteOrdersHandler{}, nil
}

type getIncompleteOrdersHandler struct{}

func (h *getIncompleteOrdersHandler) Handle(context.Context, GetIncompleteOrdersQuery) (GetIncompleteOrdersResponse, error) {
	return GetIncompleteOrdersResponse{}, nil
}
