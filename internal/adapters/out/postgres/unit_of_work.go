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

func NewUnitOfWorkFactory(db *gorm.DB) (ports.UnitOfWorkFactory, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}
	return &unitOfWorkFactory{db: db}, nil
}

type unitOfWorkFactory struct {
	db *gorm.DB
}

func (f *unitOfWorkFactory) New(ctx context.Context) (ports.UnitOfWork, error) {
	return NewUnitOfWork(f.db.WithContext(ctx))
}

type UnitOfWork struct {
	tx                *gorm.DB
	db                *gorm.DB
	committed         bool
	trackedAggregates []ddd.AggregateRoot
	//
	orderRepository   ports.OrderRepository
	courierRepository ports.CourierRepository
}

func NewUnitOfWork(db *gorm.DB) (ports.UnitOfWork, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}

	uow := &UnitOfWork{db: db}

	orderRepo, err := orderrepo.NewRepository(uow)
	if err != nil {
		return nil, err
	}
	uow.orderRepository = orderRepo

	courierRepo, err := courierrepo.NewRepository(uow)
	if err != nil {
		return nil, err
	}
	uow.courierRepository = courierRepo

	return uow, nil
}

//

func (u *UnitOfWork) CourierRepository() ports.CourierRepository {
	return u.courierRepository
}

func (u *UnitOfWork) OrderRepository() ports.OrderRepository {
	return u.orderRepository
}

func (u *UnitOfWork) Tx() *gorm.DB {
	return u.tx
}

func (u *UnitOfWork) Db() *gorm.DB {
	return u.db
}

func (u *UnitOfWork) InTx() bool {
	return u.tx != nil
}

func (u *UnitOfWork) Track(agg ddd.AggregateRoot) {
	u.trackedAggregates = append(u.trackedAggregates, agg)
}

func (u *UnitOfWork) Begin(ctx context.Context) {
	u.tx = u.db.WithContext(ctx).Begin()
	u.committed = false
}

func (u *UnitOfWork) Commit(ctx context.Context) error {
	if u.tx == nil {
		return errs.NewValueIsRequiredError("cannot commit without transaction")
	}

	if err := u.tx.WithContext(ctx).Commit().Error; err != nil {
		return err
	}

	u.committed = true
	u.clearTx()
	return nil
}

func (u *UnitOfWork) RollbackUnlessCommitted(ctx context.Context) {
	if u.tx != nil && !u.committed {
		if err := u.tx.WithContext(ctx).Rollback().Error; err != nil && !errors.Is(err, gorm.ErrInvalidTransaction) {
			log.Error(err)
		}
		u.clearTx()
	}
}

func (u *UnitOfWork) clearTx() {
	u.tx = nil
	u.trackedAggregates = nil
	u.committed = false
}
