package queries

import (
	"testing"

	"delivery/internal/adapters/out/postgres"
	"delivery/internal/adapters/out/postgres/courierrepo"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/kernel"

	"github.com/stretchr/testify/assert"
)

func Test_GetAllCouriersQuery(t *testing.T) {
	assert := assert.New(t)
	// Инициализируем окружение
	ctx, db, err := setupTest(t)
	assert.NoError(err)

	uow, err := postgres.NewUnitOfWork(db)
	assert.NoError(err)

	courier, err := courier.NewCourier("test", 5, kernel.NewRandomLocation())
	assert.NoError(err)

	repo, err := courierrepo.NewRepository(uow)
	assert.NoError(err)

	err = repo.Add(ctx, courier)
	assert.NoError(err)

	query, err := NewGetAllCouriersQuery()
	assert.NoError(err)

	handler, err := NewGetAllCouriersQueryHandler(db)
	assert.NoError(err)

	res, err := handler.Handle(ctx, query)
	assert.NoError(err)

	assert.Len(res.Couriers, 1)
}
