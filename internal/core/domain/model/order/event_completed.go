package order

import (
	"delivery/internal/pkg/ddd"
	"reflect"

	"github.com/google/uuid"
)

type CompletedDomainEvent struct {
	ID    uuid.UUID
	Name  string
	valid bool

	OrderID     uuid.UUID
	OrderStatus string
}

func NewCompletedDomainEvent(aggregate *Order) ddd.DomainEvent {
	completedDomainEvent := CompletedDomainEvent{
		ID: uuid.New(),

		OrderID:     aggregate.ID(),
		OrderStatus: aggregate.Status().String(),

		valid: true,
	}

	completedDomainEvent.Name = reflect.TypeOf(completedDomainEvent).Name()
	return &completedDomainEvent
}

func NewEmptyCompletedDomainEvent() ddd.DomainEvent {
	completedDomainEvent := CompletedDomainEvent{}
	completedDomainEvent.Name = reflect.TypeOf(completedDomainEvent).Name()
	return &completedDomainEvent
}

func (e CompletedDomainEvent) GetID() uuid.UUID { return e.ID }

func (e CompletedDomainEvent) GetName() string {
	return e.Name
}

func (e CompletedDomainEvent) Empty() bool {
	return !e.valid
}
