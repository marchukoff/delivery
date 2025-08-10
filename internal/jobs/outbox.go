package jobs

import (
	"context"
	"log"
	"time"

	"delivery/internal/adapters/out/postgres/outboxrepo"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"delivery/internal/pkg/outbox"

	"github.com/robfig/cron/v3"
)

var _ cron.Job = (*OutboxJob)(nil)

type OutboxJob struct {
	outboxRepository outboxrepo.OutboxRepository
	registry         outbox.EventRegistry
	mediatr          ddd.Mediatr
}

func NewOutboxJob(
	outboxRepository outboxrepo.OutboxRepository,
	registry outbox.EventRegistry,
	mediatr ddd.Mediatr,
) (*OutboxJob, error) {
	if outboxRepository == nil {
		return nil, errs.NewValueIsRequiredError("outboxRepository")
	}
	if registry == nil {
		return nil, errs.NewValueIsRequiredError("registry")
	}
	if mediatr == nil {
		return nil, errs.NewValueIsRequiredError("mediatr")
	}

	return &OutboxJob{
		outboxRepository: outboxRepository,
		registry:         registry,
		mediatr:          mediatr,
	}, nil
}

func (j *OutboxJob) Run() {
	ctx := context.Background()
	outboxMessages, err := j.outboxRepository.GetNotPublishedMessages()
	if err != nil {
		log.Printf("ERROR: get not published messages: %v", err)
	}

	for _, outboxMessage := range outboxMessages {
		domainEvent, err := j.registry.DecodeDomainEvent(outboxMessage)
		if err != nil {
			log.Printf("ERROR: decode domain event: %v", err)
			continue
		}
		log.Printf("INFO: got domain event: %+v", domainEvent)

		err = j.mediatr.Publish(ctx, domainEvent)
		if err != nil {
			log.Printf("ERROR: publish domain event: %v", err)
			continue
		}

		now := time.Now().UTC()
		outboxMessage.ProcessedAtUtc = &now
		err = j.outboxRepository.Update(ctx, outboxMessage)
		if err != nil {
			log.Printf("ERROR: update domain event: %v", err)
			continue
		}
	}
}
