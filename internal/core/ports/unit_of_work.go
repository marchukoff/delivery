package ports

import (
	"context"
)

type UnitOfWorkFactory func() (UnitOfWork, error)

type UnitOfWork interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context)
	OrderRepository() OrderRepository
	CourierRepository() CourierRepository
}
