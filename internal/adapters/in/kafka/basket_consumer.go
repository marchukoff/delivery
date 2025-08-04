package kafka

import (
	"context"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/generated/queues/basketconfirmedpb"
	"delivery/internal/pkg/errs"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type BasketConsumer interface {
	Consume() error
	Close() error
}

var _ BasketConsumer = &basketConsumer{}

type basketConsumer struct {
	topic                     string
	consumerGroup             sarama.ConsumerGroup
	createOrderCommandHandler commands.CreateOrderCommandHandler
	cancel                    chan struct{}
}

func NewBasketConsumer(
	brokers []string,
	group string,
	topic string,
	createOrderCommandHandler commands.CreateOrderCommandHandler,
) (BasketConsumer, error) {
	if len(brokers) == 0 {
		return nil, errs.NewValueIsRequiredError("brokers")
	}

	if group == "" {
		return nil, errs.NewValueIsRequiredError("group")
	}

	if topic == "" {
		return nil, errs.NewValueIsRequiredError("topic")
	}

	if createOrderCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("createOrderCommandHandler")
	}

	saramaCfg := sarama.NewConfig()
	saramaCfg.Version = sarama.V3_4_0_0
	saramaCfg.Consumer.Return.Errors = true
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, group, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("create consumer group: %w", err)
	}

	return &basketConsumer{
		topic:                     topic,
		consumerGroup:             consumerGroup,
		createOrderCommandHandler: createOrderCommandHandler,
		cancel:                    make(chan struct{}),
	}, nil
}

func (c *basketConsumer) Close() error {
	close(c.cancel)
	return c.consumerGroup.Close()
}

func (c *basketConsumer) Consume() error {
	ctx := context.Background()

	for {
		select {
		case <-c.cancel:
			return nil
		default:
			err := c.consumerGroup.Consume(ctx, []string{c.topic}, c)
			if err != nil {
				log.Printf("ERROR: consumer got an error: %v", err)
				return err
			}
		}
	}
}

var _ sarama.ConsumerGroupHandler = &basketConsumer{}

func (c *basketConsumer) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *basketConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *basketConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		ctx := context.Background()
		log.Printf("INFO: received: topic = %s, partition = %d, offset = %d, key = %s, value = %s\n",
			message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))

		var event basketconfirmedpb.BasketConfirmedIntegrationEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			log.Printf("ERROR: failed to unmarshal message: %v", err)
			session.MarkMessage(message, "")
			continue
		}

		orderID, err := uuid.Parse(event.BasketId)
		if err != nil {
			log.Printf("ERROR: failed to parse order id: %v", err)
			continue
		}

		cmd, err := commands.NewCreateOrderCommand(orderID, event.Address.Street, int(event.Volume))
		if err != nil {
			log.Printf("ERROR: failed to create createOrder command: %v", err)
			session.MarkMessage(message, "")
			continue
		}

		if err := c.createOrderCommandHandler.Handle(ctx, cmd); err != nil {
			log.Printf("ERROR: failed to handle createOrder command: %v", err)
			continue
		}

		session.MarkMessage(message, "")
	}

	return nil
}
