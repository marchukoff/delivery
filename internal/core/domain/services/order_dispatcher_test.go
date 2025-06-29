package services_test

import (
	"testing"

	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/domain/services"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_orderDispatcher_Dispatch(t *testing.T) {
	assert := assert.New(t)
	strong := func() *courier.Courier {
		cur, err := courier.NewCourier("strong", 5, kernel.NewRandomLocation())
		assert.NoError(err)
		assert.NoError(cur.AddStoragePlace("BigBag", 20))
		return cur
	}()

	dispatcher := services.NewOrderDispatcher()

	tests := []struct {
		name     string
		order    *order.Order
		couriers []*courier.Courier
		want     *courier.Courier
		wantErr  error
	}{
		{
			name: "good",
			order: func() *order.Order {
				order, err := order.NewOrder(uuid.New(), kernel.NewRandomLocation(), 15)
				assert.NoError(err)
				return order
			}(),
			couriers: []*courier.Courier{strong, func() *courier.Courier {
				cur, err := courier.NewCourier("weak", 1, kernel.NewRandomLocation())
				assert.NoError(err)
				return cur
			}()},
			want:    strong,
			wantErr: nil,
		},
		{
			name:     "bad no order",
			order:    nil,
			couriers: []*courier.Courier{strong},
			want:     nil,
			wantErr:  errs.ErrValueIsRequired,
		},
		{
			name: "bad no courier",
			order: func() *order.Order {
				order, err := order.NewOrder(uuid.New(), kernel.NewRandomLocation(), 15)
				assert.NoError(err)
				return order
			}(),
			couriers: nil,
			want:     nil,
			wantErr:  services.ErrNoRightCourier,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dispatcher.Dispatch(tt.order, tt.couriers)
			if err != nil {
				assert.ErrorIs(err, tt.wantErr)
			} else {
				assert.True(tt.want.Equal(got))
				assert.Equal(order.StatusAssigned, tt.order.Status())
				assert.Equal(got.ID(), *tt.order.CourierID())
			}
		})
	}
}
