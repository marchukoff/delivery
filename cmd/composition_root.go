package cmd

import (
	"io"
	"log"
	"sync"

	kafkain "delivery/internal/adapters/in/kafka"
	grpcout "delivery/internal/adapters/out/grpc"
	kafkaout "delivery/internal/adapters/out/kafka"
	"delivery/internal/adapters/out/postgres"
	"delivery/internal/core/application/eventhandlers"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/ddd"

	"gorm.io/gorm"
)

type CompositionRoot struct {
	config    Config
	db        *gorm.DB
	geoClient ports.GeoClient
	onceGeo   sync.Once
	//
	closers []io.Closer
}

func NewCompositionRoot(cfg Config, db *gorm.DB) *CompositionRoot {
	return &CompositionRoot{config: cfg, db: db}
}

func (c *CompositionRoot) NewOrderDispatcherService() services.OrderDispatcher {
	return services.NewOrderDispatcher()
}

func (c *CompositionRoot) NewUnitOfWorkFactory() ports.UnitOfWorkFactory {
	factory, err := postgres.NewUnitOfWorkFactory(c.db, c.NewMediatrWithSubscriptions())
	if err != nil {
		log.Fatalf("ERROR: could not create unit of work factory: %v", err)
	}
	return factory
}

func (c *CompositionRoot) NewCreateOrderCommandHandler() commands.CreateOrderCommandHandler {
	h, err := commands.NewCreateOrderCommandHandler(c.NewUnitOfWorkFactory(), c.NewGeoClient())
	if err != nil {
		log.Fatalf("ERROR: cannot create CreateOrderCommandHandler: %v", err)
	}
	return h
}

func (c *CompositionRoot) NewCreateCourierCommandHandler() commands.CreateCourierCommandHandler {
	h, err := commands.NewCreateCourierCommandHandler(c.NewUnitOfWorkFactory())
	if err != nil {
		log.Fatalf("ERROR: cannot create CreateCourierCommandHandler: %v", err)
	}
	return h
}

func (c *CompositionRoot) NewGetAllCouriersQueryHandler() queries.GetAllCouriersQueryHandler {
	h, err := queries.NewGetAllCouriersQueryHandler(c.db)
	if err != nil {
		log.Fatalf("ERROR: cannot create GetAllCouriersQueryHandler: %v", err)
	}
	return h
}

func (c *CompositionRoot) NewGetIncompletedOrdersQueryHandler() queries.GetIncompleteOrdersQueryHandler {
	h, err := queries.NewGetIncompleteOrdersHandler(c.db)
	if err != nil {
		log.Fatalf("ERROR: cannot create GetIncompletedOrdersQueryHandler: %v", err)
	}
	return h
}

func (cr *CompositionRoot) NewGeoClient() ports.GeoClient {
	cr.onceGeo.Do(func() {
		client, err := grpcout.NewClient(cr.config.GeoServiceGrpcHost)
		if err != nil {
			log.Fatalf("ERROR: create GeoClient: %v", err)
		}
		cr.RegisterCloser(client)
		cr.geoClient = client
	})
	return cr.geoClient
}

func (cr *CompositionRoot) NewBasketConsumer() kafkain.BasketConsumer {
	consumer, err := kafkain.NewBasketConsumer(
		[]string{cr.config.KafkaHost},
		cr.config.KafkaConsumerGroup,
		cr.config.KafkaBasketConfirmedTopic,
		cr.NewCreateOrderCommandHandler(),
	)
	if err != nil {
		log.Fatalf("ERROR: create BasketConsumer: %v", err)
	}

	cr.RegisterCloser(consumer)
	return consumer
}

func (cr *CompositionRoot) NewOrderProducer() ports.OrderProducer {
	producer, err := kafkaout.NewOrderProducer([]string{cr.config.KafkaHost}, cr.config.KafkaOrderChangedTopic)
	if err != nil {
		log.Fatalf("ERROR: cannot create OrderProducer: %v", err)
	}

	cr.RegisterCloser(producer)
	return producer
}

func (cr *CompositionRoot) NewOrderCompletedDomainEventHandler() ddd.EventHandler {
	producer := cr.NewOrderProducer()
	handler, err := eventhandlers.NewOrderCompletedDomainEventHandler(producer)
	if err != nil {
		log.Fatalf("ERROR: cannot create OrderCompletedDomainEventHandler: %v", err)
	}

	return handler
}

func (cr *CompositionRoot) NewMediatrWithSubscriptions() ddd.Mediatr {
	mediatr := ddd.NewMediatr()
	mediatr.Subscribe(cr.NewOrderCompletedDomainEventHandler(), order.NewEmptyCompletedDomainEvent())
	return mediatr
}

func (cr *CompositionRoot) RegisterCloser(c io.Closer) {
	cr.closers = append(cr.closers, c)
}

func (cr *CompositionRoot) CloseAll() {
	for _, closer := range cr.closers {
		if err := closer.Close(); err != nil {
			log.Printf("ERROR: closing resource: %v", err)
		}
	}
}
