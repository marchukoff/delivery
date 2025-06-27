package postgres

import (
	"context"
	"errors"

	"delivery/internal/adapters/out/postgres/courierrepo"
	"delivery/internal/adapters/out/postgres/orderrepo"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

var _ ports.UnitOfWork = &UnitOfWork{}

type UnitOfWork struct {
	tx                *gorm.DB
	db                *gorm.DB
	trackedAggregates []ddd.AggregateRoot
	courierRepository ports.CourierRepository
	orderRepository   ports.OrderRepository
}

func NewUnitOfWork(db *gorm.DB) (ports.UnitOfWork, error) {
	var (
		uow = &UnitOfWork{db: db}
		err error
	)

	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}

	uow.courierRepository, err = courierrepo.NewRepository(uow)
	if err != nil {
		return nil, err
	}

	uow.orderRepository, err = orderrepo.NewRepository(uow)
	if err != nil {
		return nil, err
	}

	return uow, nil
}

func (u *UnitOfWork) Tx() *gorm.DB { return u.tx }

func (u *UnitOfWork) Db() *gorm.DB { return u.db }

func (u *UnitOfWork) InTx() bool { return u.tx != nil }

func (u *UnitOfWork) Track(agg ddd.AggregateRoot) {
	u.trackedAggregates = append(u.trackedAggregates, agg)
}

func (u *UnitOfWork) CourierRepository() ports.CourierRepository {
	return u.courierRepository
}

func (u *UnitOfWork) OrderRepository() ports.OrderRepository {
	return u.orderRepository
}

func (u *UnitOfWork) Begin(ctx context.Context) {
	u.tx = u.db.WithContext(ctx).Begin()
}

func (u *UnitOfWork) Commit(ctx context.Context) (err error) {
	if u.tx == nil {
		return errs.NewValueIsRequiredError("cannot commit without transaction")
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, u.tx.WithContext(ctx).Rollback().Error)
			if err != nil && !errors.Is(err, gorm.ErrInvalidTransaction) {
				log.Error(err)
			}
		}

		u.clearTx()
	}()

	if err = u.tx.WithContext(ctx).Commit().Error; err != nil {
		return err
	}

	return nil
}

func (u *UnitOfWork) clearTx() { u.tx, u.trackedAggregates = nil, nil }
