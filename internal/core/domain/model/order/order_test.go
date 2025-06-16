package order_test

import (
	"testing"

	"delivery/internal/core/domain/kernel"
	. "delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewOrder(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		orderID  uuid.UUID
		location kernel.Location
		volume   int
	}

	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "good",
			args: args{
				orderID:  uuid.New(),
				location: kernel.NewRandomLocation(),
				volume:   1,
			},
			want: nil,
		},
		{
			name: "bad order id",
			args: args{
				orderID:  uuid.Nil,
				location: kernel.NewRandomLocation(),
				volume:   1,
			},
			want: errs.ErrValueIsRequired,
		},
		{
			name: "bad location",
			args: args{
				orderID:  uuid.New(),
				location: kernel.Location{},
				volume:   1,
			},
			want: errs.ErrValueIsRequired,
		},
		{
			name: "bad volume",
			args: args{
				orderID:  uuid.New(),
				location: kernel.NewRandomLocation(),
				volume:   -1,
			},
			want: errs.ErrValueIsRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOrder(tt.args.orderID, tt.args.location, tt.args.volume)
			if err != nil {
				assert.ErrorIs(err, tt.want)
				assert.Equal(uuid.Nil, got.ID())
				assert.Equal(kernel.Location{}, got.Location())
				assert.Equal(0, got.Volume())
			} else {
				assert.Equal(tt.args.orderID, got.ID())
				assert.Equal(tt.args.location, got.Location())
				assert.Equal(tt.args.volume, got.Volume())
			}
		})
	}
}

func TestOrder_Assign(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name      string
		courierID uuid.UUID
		order     *Order
		want      error
	}{
		{
			name:      "good",
			courierID: uuid.New(),
			order: func() *Order {
				order, err := NewOrder(uuid.New(), kernel.NewRandomLocation(), 1)
				assert.NoError(err)
				return order
			}(),
			want: nil,
		},
		{
			name: "bad courierID",
			order: func() *Order {
				order, err := NewOrder(uuid.New(), kernel.NewRandomLocation(), 1)
				assert.NoError(err)
				return order
			}(),
			want: errs.ErrValueIsRequired,
		},
		{
			name:      "bad reassign",
			courierID: uuid.New(),
			order: func() *Order {
				order, err := NewOrder(uuid.New(), kernel.NewRandomLocation(), 1)
				assert.NoError(err)
				err = order.Assign(uuid.New())
				assert.NoError(err)
				return order
			}(),
			want: errs.ErrExpectationFailed,
		},
		{
			name:      "bad not initialized order",
			courierID: uuid.New(),
			want:      ErrOrderNotInitialized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.order.Assign(tt.courierID); err != nil {
				assert.ErrorIs(err, tt.want)
			} else {
				assert.Equal(StatusAssigned, tt.order.Status())
			}
		})
	}
}

func TestOrder_Complete(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name  string
		order *Order
		want  error
	}{
		{
			name: "good",
			order: func() *Order {
				order, err := NewOrder(uuid.New(), kernel.NewRandomLocation(), 1)
				assert.NoError(err)
				err = order.Assign(uuid.New())
				assert.NoError(err)
				return order
			}(),
			want: nil,
		},
		{
			name: "bad not assigned",
			order: func() *Order {
				order, err := NewOrder(uuid.New(), kernel.NewRandomLocation(), 1)
				assert.NoError(err)
				return order
			}(),
			want: errs.ErrExpectationFailed,
		},
		{
			name: "bad order not initialized",
			want: ErrOrderNotInitialized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.order.Complete(); err != nil {
				assert.ErrorIs(err, tt.want)
			} else {
				assert.Equal(StatusCompleted, tt.order.Status())
			}
		})
	}
}
