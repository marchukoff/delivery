package ports

import (
	"context"
	"delivery/internal/pkg/ddd"
	"io"
)

type OrderProducer interface {
	Publish(ctx context.Context, domainEvent ddd.DomainEvent) error
	io.Closer
}
