package postgres

import (
	"context"

	"delivery/internal/adapters/out/postgres/courierrepo"
	"delivery/internal/adapters/out/postgres/orderrepo"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"

	"gorm.io/gorm"
)

var _ ports.UnitOfWork = (*unitOfWork)(nil)

func NewUnitOfWorkFactory(db *gorm.DB) ports.UnitOfWorkFactory {
	return func() (ports.UnitOfWork, error) {
		if db == nil {
			return nil, errs.NewValueIsRequiredError("db")
		}

		return &unitOfWork{db: db.Begin()}, nil
	}
}

type unitOfWork struct {
	db                *gorm.DB
	orderRepository   ports.OrderRepository
	courierRepository ports.CourierRepository
}

func (u *unitOfWork) CourierRepository() ports.CourierRepository {
	if u.courierRepository == nil {
		u.courierRepository = courierrepo.NewRepository(u.db)
	}
	return u.courierRepository
}

func (u *unitOfWork) OrderRepository() ports.OrderRepository {
	if u.orderRepository == nil {
		u.orderRepository = orderrepo.NewRepository(u.db)
	}
	return u.orderRepository
}

func (u *unitOfWork) Commit(_ context.Context) error {
	return u.db.Commit().Error
}

func (u *unitOfWork) Rollback(_ context.Context) {
	u.db.Rollback()
}
