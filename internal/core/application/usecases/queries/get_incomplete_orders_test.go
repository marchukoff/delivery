package queries

import (
	"context"
	"testing"

	"delivery/internal/adapters/out/postgres"
	"delivery/internal/adapters/out/postgres/courierrepo"
	"delivery/internal/adapters/out/postgres/orderrepo"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/testcnts"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	postgresgorm "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_GetIncompleteOrdersQuery(t *testing.T) {
	assert := assert.New(t)
	// Инициализируем окружение
	ctx, db, err := setupTest(t)
	assert.NoError(err)

	uow, err := postgres.NewUnitOfWork(db)
	assert.NoError(err)

	order, err := order.NewOrder(uuid.New(), kernel.NewRandomLocation(), 1)
	assert.NoError(err)

	repo, err := orderrepo.NewRepository(uow)
	assert.NoError(err)

	err = repo.Add(ctx, order)
	assert.NoError(err)

	query, err := NewGetIncompleteOrdersQuery()
	assert.NoError(err)

	handler, err := NewGetIncompleteOrdersHandler(db)
	assert.NoError(err)

	res, err := handler.Handle(ctx, query)
	assert.NoError(err)

	assert.Len(res.Orders, 1)
}

func setupTest(t *testing.T) (context.Context, *gorm.DB, error) {
	assert := assert.New(t)
	ctx := context.Background()
	postgresContainer, dsn, err := testcnts.StartPostgresContainer(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Подключаемся к БД через Gorm
	db, err := gorm.Open(postgresgorm.Open(dsn), &gorm.Config{})
	assert.NoError(err)

	// Авто миграция (создаём таблицу)
	err = db.AutoMigrate(&courierrepo.StoragePlaceDTO{})
	assert.NoError(err)
	err = db.AutoMigrate(&courierrepo.CourierDTO{})
	assert.NoError(err)
	err = db.AutoMigrate(&orderrepo.OrderDTO{})
	assert.NoError(err)

	// Очистка выполняется после завершения теста
	t.Cleanup(func() {
		err := postgresContainer.Terminate(ctx)
		assert.NoError(err)
	})

	return ctx, db, nil
}
