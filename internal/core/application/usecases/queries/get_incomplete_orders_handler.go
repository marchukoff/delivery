package queries

import (
	"context"

	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"

	"gorm.io/gorm"
)

type GetIncompleteOrdersQueryHandler interface {
	Handle(context.Context, GetIncompleteOrdersQuery) (GetIncompleteOrdersResponse, error)
}

func NewGetIncompleteOrdersHandler(db *gorm.DB) (*getIncompleteOrdersHandler, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}
	return &getIncompleteOrdersHandler{db: db}, nil
}

type getIncompleteOrdersHandler struct {
	db *gorm.DB
}

func (h *getIncompleteOrdersHandler) Handle(ctx context.Context, query GetIncompleteOrdersQuery) (GetIncompleteOrdersResponse, error) {
	if !query.IsValid() {
		return GetIncompleteOrdersResponse{}, errs.NewValueIsRequiredError("query")
	}

	var orders []Order
	err := h.db.WithContext(ctx).
		Raw("SELECT id, location_x, location_y FROM orders where status!=?", order.StatusCompleted).
		Scan(&orders).
		Error
	if err != nil {
		return GetIncompleteOrdersResponse{}, err
	}

	return GetIncompleteOrdersResponse{Orders: orders}, nil
}
