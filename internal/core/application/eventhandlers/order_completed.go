package eventhandlers

import (
	"context"

	"delivery/internal/core/ports"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
)

type orderCompletedDomainEventHandler struct {
	producer ports.OrderProducer
}

func NewOrderCompletedDomainEventHandler(orderProducer ports.OrderProducer) (ddd.EventHandler, error) {
	if orderProducer == nil {
		return nil, errs.NewValueIsRequiredError("orderProducer")
	}

	return &orderCompletedDomainEventHandler{producer: orderProducer}, nil
}

func (h *orderCompletedDomainEventHandler) Handle(ctx context.Context, event ddd.DomainEvent) error {
	return h.producer.Publish(ctx, event)
}
