package postgres

import (
	"context"
	"testing"

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
	"gorm.io/gorm/clause"
)

func Test_CourierRepositoryShouldCanAddCourier(t *testing.T) {
	assert := assert.New(t)
	// Инициализируем окружение
	ctx, db, err := setupTest(t)
	assert.NoError(err)

	// Создаем UnitOfWork
	factory := NewUnitOfWorkFactory(db)
	assert.NotNil(factory)

	// Вызываем Add
	name := "test"
	speed := 5
	loc := kernel.NewRandomLocation()
	courier, err := courier.NewCourier(name, speed, loc)
	assert.NoError(err)

	uow, err := factory()
	assert.NoError(err)
	err = uow.CourierRepository().Add(ctx, courier)
	assert.NoError(err)
	err = uow.Commit(ctx)
	assert.NoError(err)

	// Считываем данные из БД
	var dto courierrepo.CourierDTO
	err = db.Preload(clause.Associations).
		First(&dto, "id = ?", courier.ID()).
		Error
	assert.NoError(err)

	// Проверяем эквивалентность
	assert.Equal(courier.ID(), dto.ID)
	assert.Equal(courier.Name(), dto.Name)
	assert.Equal(courier.Speed(), dto.Speed)
	assert.Equal(courier.Location().X(), dto.Location.X)
	assert.Equal(courier.Location().Y(), dto.Location.Y)
	assert.Equal(len(courier.StoragePlaces()), len(dto.StoragePlaces))
}

func Test_OrderRepositoryShouldCanAddOrder(t *testing.T) {
	assert := assert.New(t)
	// Инициализируем окружение
	ctx, db, err := setupTest(t)
	assert.NoError(err)

	// Создаем UnitOfWork
	factory := NewUnitOfWorkFactory(db)
	assert.NotNil(factory())
	uow, err := factory()
	assert.NoError(err)

	// Вызываем Add

	id := uuid.New()
	volume := 5
	loc := kernel.NewRandomLocation()
	order, err := order.NewOrder(id, loc, volume)
	assert.NoError(err)
	err = uow.OrderRepository().Add(ctx, order)
	assert.NoError(err)
	err = uow.Commit(ctx)
	assert.NoError(err)

	// Считываем данные из БД
	var dto orderrepo.OrderDTO
	err = db.First(&dto, "id = ?", order.ID()).Error
	assert.NoError(err)

	// Проверяем эквивалентность
	assert.Equal(order.ID(), dto.ID)
	assert.Equal(order.CourierID(), dto.CourierID)
	assert.Equal(order.Location().X(), dto.Location.X)
	assert.Equal(order.Location().Y(), dto.Location.Y)
	assert.Equal(order.Volume(), dto.Volume)
	assert.Equal(order.Status(), dto.Status)
}

func setupTest(t *testing.T) (context.Context, *gorm.DB, error) {
	ctx := context.Background()
	postgresContainer, dsn, err := testcnts.StartPostgresContainer(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Подключаемся к БД через Gorm
	db, err := gorm.Open(postgresgorm.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	// Авто миграция (создаём таблицу)
	err = db.AutoMigrate(&courierrepo.StoragePlaceDTO{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&courierrepo.CourierDTO{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&orderrepo.OrderDTO{})
	assert.NoError(t, err)

	// Очистка выполняется после завершения теста
	t.Cleanup(func() {
		err := postgresContainer.Terminate(ctx)
		assert.NoError(t, err)
	})

	return ctx, db, nil
}
