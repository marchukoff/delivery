package outbox

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
)

type EventRegistry interface {
	RegisterDomainEvent(eventType reflect.Type) error
	DecodeDomainEvent(event *Message) (ddd.DomainEvent, error)
}

var _ EventRegistry = &eventRegistry{}

type eventRegistry struct {
	eventRegistry map[string]reflect.Type
}

func NewEventRegistry() (EventRegistry, error) {
	return &eventRegistry{
		eventRegistry: make(map[string]reflect.Type),
	}, nil
}

func (r *eventRegistry) RegisterDomainEvent(eventType reflect.Type) error {
	if eventType == nil {
		return errs.NewValueIsRequiredError("eventType")
	}

	eventName := eventType.Name()
	r.eventRegistry[eventName] = eventType
	return nil
}

func EncodeDomainEvent(domainEvent ddd.DomainEvent) (Message, error) {
	payload, err := json.Marshal(domainEvent)
	if err != nil {
		return Message{}, fmt.Errorf("marshal event: %w", err)
	}

	return Message{
		ID:             domainEvent.GetID(),
		Name:           domainEvent.GetName(),
		Payload:        payload,
		OccurredAtUtc:  time.Now().UTC(),
		ProcessedAtUtc: nil,
	}, nil
}

func EncodeDomainEvents(domainEvent []ddd.DomainEvent) ([]Message, error) {
	outboxMessages := make([]Message, 0, len(domainEvent))
	for _, event := range domainEvent {
		event, err := EncodeDomainEvent(event)
		if err != nil {
			return nil, err
		}
		outboxMessages = append(outboxMessages, event)
	}
	return outboxMessages, nil
}

func (r *eventRegistry) DecodeDomainEvent(outboxMessage *Message) (domainEvent ddd.DomainEvent, err error) {
	defer func() {
		// we have a lot reflect magic and we want cover our ass
		if fail := recover(); fail != nil {
			err = errors.Join(err, fmt.Errorf("recovered in DecodeDomainEvent: %v", fail))
		}
	}()

	t, ok := r.eventRegistry[outboxMessage.Name]
	if !ok {
		return nil, fmt.Errorf("unknown outboxMessage type: %s", outboxMessage.Name)
	}

	eventPtr := reflect.New(t).Interface()
	if err := json.Unmarshal(outboxMessage.Payload, eventPtr); err != nil {
		return nil, fmt.Errorf("decode payload: %w", err)
	}

	domainEvent, ok = eventPtr.(ddd.DomainEvent)
	if !ok {
		return nil, fmt.Errorf("cast outboxMessage to DomainEvent: %w",
			errs.NewExpectationFailedError("outboxMessage", outboxMessage.Name, "ddd.DomainEvent"))
	}

	return domainEvent, nil
}
