package queries

import (
	"context"

	"delivery/internal/pkg/errs"

	"gorm.io/gorm"
)

type GetAllCouriersQueryHandler interface {
	Handle(context.Context, GetAllCouriersQuery) (GetAllCouriersResponse, error)
}

func NewGetAllCouriersQueryHandler(db *gorm.DB) (*getAllCouriersQueryHandler, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}
	return &getAllCouriersQueryHandler{db: db}, nil
}

type getAllCouriersQueryHandler struct {
	db *gorm.DB
}

func (h *getAllCouriersQueryHandler) Handle(ctx context.Context, query GetAllCouriersQuery) (GetAllCouriersResponse, error) {
	if !query.IsValid() {
		return GetAllCouriersResponse{}, errs.NewValueIsRequiredError("query")
	}

	var couriers []Courier

	err := h.db.WithContext(ctx).
		Raw("SELECT id,name, location_x, location_y FROM couriers").
		Scan(&couriers).
		Error
	if err != nil {
		return GetAllCouriersResponse{}, err
	}

	return GetAllCouriersResponse{Couriers: couriers}, nil
}
