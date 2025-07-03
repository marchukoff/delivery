package commands

import (
	"context"
	"testing"

	"delivery/internal/adapters/out/postgres"
	"delivery/internal/adapters/out/postgres/courierrepo"
	"delivery/internal/adapters/out/postgres/orderrepo"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/testcnts"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	postgresgorm "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_MoveCouriersCommand(t *testing.T) {
	assert := assert.New(t)
	// Инициализируем окружение
	ctx, db, err := setupTest(t)
	assert.NoError(err)

	// Создаем UnitOfWork
	factory := postgres.NewUnitOfWorkFactory(db)

	// Вызываем Add
	name, speed := "test", 5
	loc1, err := kernel.NewLocation(1, 1)
	assert.NoError(err)
	loc2, err := kernel.NewLocation(9, 9)
	assert.NoError(err)

	courier, err := courier.NewCourier(name, speed, loc1)
	assert.NoError(err)

	order, err := order.NewOrder(uuid.New(), loc2, 1)
	assert.NoError(err)

	assert.NoError(order.Assign(courier.ID()))
	assert.NoError(courier.TakeOrder(order))
	// save
	uow, err := factory()
	assert.NoError(err)
	err = uow.CourierRepository().Add(ctx, courier)
	assert.NoError(err)

	err = uow.OrderRepository().Add(ctx, order)
	assert.NoError(err)
	assert.NoError(uow.Commit(ctx))

	// change
	command, err := NewMoveCouriersCommand()
	assert.NoError(err)

	handler, err := NewMoveCouriersCommandHandler(factory)
	assert.NoError(err)
	err = handler.Handle(ctx, command)
	assert.NoError(err)

	// retrieve
	// FIXME: check location
	// uow, err = factory()
	// assert.NoError(err)
	// courier2, err := uow.CourierRepository().Get(ctx, courier.ID())
	// assert.NoError(err)
	// assert.NoError(uow.Commit(ctx))
	// assert.False(loc1.Equals(courier2.Location()))
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
