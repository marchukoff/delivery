package ports

import (
	"context"
)

type UnitOfWorkFactory interface {
	New(ctx context.Context) (UnitOfWork, error)
}

type UnitOfWork interface {
	// DB spesific
	Begin(ctx context.Context)
	Commit(ctx context.Context) error
	RollbackUnlessCommitted(ctx context.Context)
	// Domain specific
	OrderRepository() OrderRepository
	CourierRepository() CourierRepository
}
