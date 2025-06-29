package courier_test

import (
	"math"
	"testing"

	. "delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewCourier(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		name     string
		speed    int
		location kernel.Location
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "good",
			args: args{
				name:     "R2D2",
				speed:    1,
				location: kernel.NewRandomLocation(),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCourier(tt.args.name, tt.args.speed, tt.args.location)
			if err != nil {
				assert.ErrorIs(err, tt.want)
				assert.Empty(got.StoragePlaces())
				assert.Empty(got.Name())
				assert.Empty(got.Speed())
				assert.Equal(kernel.Location{}, got.Location())
			} else {
				assert.Equal(tt.args.name, got.Name())
				assert.Equal(tt.args.speed, got.Speed())
				assert.Equal(tt.args.location, got.Location())
				assert.NotEmpty(got.StoragePlaces())
			}
		})
	}
}

func TestCourier_AddStoragePlace(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		name   string
		volume int
	}

	tests := []struct {
		name    string
		args    args
		courier *Courier
		want    error
	}{
		{
			name: "good",
			args: args{name: "Pocket", volume: 1},
			courier: func() *Courier {
				c, err := NewCourier("test", 1, kernel.NewRandomLocation())
				assert.NoError(err)
				return c
			}(),
			want: nil,
		},
		{
			name: "bad storage place",
			args: args{volume: 1},
			courier: func() *Courier {
				c, err := NewCourier("test", 1, kernel.NewRandomLocation())
				assert.NoError(err)
				return c
			}(),
			want: errs.ErrValueIsRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.courier.AddStoragePlace(tt.args.name, tt.args.volume); err != nil {
				assert.ErrorIs(err, tt.want)
			} else {
				assert.NotEmpty(tt.courier.StoragePlaces())
				assert.Greater(len(tt.courier.StoragePlaces()), 1)
			}
		})
	}
}

func TestCourier_CanTakeOrder(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		courier *Courier
		order   *order.Order
		want    bool
		wantErr error
	}{
		{
			name: "good",
			courier: func() *Courier {
				c, err := NewCourier("test", 1, kernel.NewRandomLocation())
				assert.NoError(err)
				return c
			}(),
			order: func() *order.Order {
				x, err := order.NewOrder(uuid.New(), kernel.NewRandomLocation(), 1)
				assert.NoError(err)
				return x
			}(),
			want:    true,
			wantErr: nil,
		},
		{
			name: "good overweigth",
			courier: func() *Courier {
				c, err := NewCourier("test", 1, kernel.NewRandomLocation())
				assert.NoError(err)
				return c
			}(),
			order: func() *order.Order {
				x, err := order.NewOrder(uuid.New(), kernel.NewRandomLocation(), 100)
				assert.NoError(err)
				return x
			}(),
			want:    false,
			wantErr: nil,
		},
		{
			name:    "bad courier not initialized",
			courier: new(Courier),
			order: func() *order.Order {
				x, err := order.NewOrder(uuid.New(), kernel.NewRandomLocation(), 100)
				assert.NoError(err)
				return x
			}(),
			want:    false,
			wantErr: nil,
		},
		{
			name: "bad nil order",
			courier: func() *Courier {
				c, err := NewCourier("test", 1, kernel.NewRandomLocation())
				assert.NoError(err)
				return c
			}(),
			want:    false,
			wantErr: errs.ErrValueIsRequired,
		},
		{
			name: "bad order not initialized",
			courier: func() *Courier {
				c, err := NewCourier("test", 1, kernel.NewRandomLocation())
				assert.NoError(err)
				return c
			}(),
			order:   new(order.Order),
			want:    false,
			wantErr: errs.ErrValueIsInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.courier.CanTakeOrder(tt.order)
			if err != nil {
				assert.False(got)
				assert.ErrorIs(err, tt.wantErr)
			} else {
				assert.Equal(tt.want, got)
			}
		})
	}
}

func TestCourier_TakeOrder(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name    string
		courier *Courier
		order   *order.Order
		want    error
	}{
		{
			name: "good",
			courier: func() *Courier {
				c, err := NewCourier("test", 1, kernel.NewRandomLocation())
				assert.NoError(err)
				return c
			}(),
			order: func() *order.Order {
				x, err := order.NewOrder(uuid.New(), kernel.NewRandomLocation(), 1)
				assert.NoError(err)
				return x
			}(),
			want: nil,
		},
		{
			name: "good too big order",
			courier: func() *Courier {
				c, err := NewCourier("test", 1, kernel.NewRandomLocation())
				assert.NoError(err)
				return c
			}(),
			order: func() *order.Order {
				x, err := order.NewOrder(uuid.New(), kernel.NewRandomLocation(), math.MaxInt)
				assert.NoError(err)
				return x
			}(),
			want: ErrNoSuitableStoragePlace,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.courier.TakeOrder(tt.order); err != nil {
				assert.ErrorIs(err, tt.want)
			} else {
				assert.Contains(func() []uuid.UUID {
					var ids []uuid.UUID
					for _, sp := range tt.courier.StoragePlaces() {
						if sp.IsOccupied() {
							ids = append(ids, *sp.OrderID())
						}
					}
					return ids
				}(), tt.order.ID())
			}
		})
	}
}

func TestCourier_CompleteOrder(t *testing.T) {
	assert := assert.New(t)

	one, err := order.NewOrder(uuid.New(), kernel.NewRandomLocation(), 1)
	assert.NoError(err)

	tests := []struct {
		name    string
		courier *Courier
		order   *order.Order
		want    error
	}{
		{
			name: "good",
			courier: func() *Courier {
				c, err := NewCourier("test", 1, kernel.NewRandomLocation())
				assert.NoError(err)
				err = c.TakeOrder(one)
				assert.NoError(err)
				err = one.Assign(c.ID())
				assert.NoError(err)
				return c
			}(),
			order: one,
			want:  nil,
		},
		{
			name: "good other order",
			courier: func() *Courier {
				c, err := NewCourier("test", 1, kernel.NewRandomLocation())
				assert.NoError(err)
				return c
			}(),
			order: func() *order.Order {
				x, err := order.NewOrder(uuid.New(), kernel.NewRandomLocation(), 1)
				assert.NoError(err)
				return x
			}(),
			want: ErrStoragePlaceNotInitialized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.courier.CompleteOrder(tt.order); err != nil {
				assert.ErrorIs(err, tt.want)
			} else {
				assert.NotContains(func() []uuid.UUID {
					var ids []uuid.UUID
					for _, sp := range tt.courier.StoragePlaces() {
						if sp.IsOccupied() {
							ids = append(ids, *sp.OrderID())
						}
					}
					return ids
				}(), tt.order.ID())
			}
		})
	}
}

func TestCourier_CalculateTimeToLocation(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name    string
		courier *Courier
		target  kernel.Location
		want    float64
		wantErr error
	}{
		{
			name: "good",
			courier: func() *Courier {
				loc, err := kernel.NewLocation(1, 1)
				assert.NoError(err)
				c, err := NewCourier("test", 2, loc)
				assert.NoError(err)
				return c
			}(),
			target: func() kernel.Location {
				loc, err := kernel.NewLocation(5, 5)
				assert.NoError(err)
				return loc
			}(),
			want:    4,
			wantErr: nil,
		},
		{
			name: "good target achieved",
			courier: func() *Courier {
				loc, err := kernel.NewLocation(1, 1)
				assert.NoError(err)
				c, err := NewCourier("test", 2, loc)
				assert.NoError(err)
				return c
			}(),
			target: func() kernel.Location {
				loc, err := kernel.NewLocation(1, 1)
				assert.NoError(err)
				return loc
			}(),
			want:    0,
			wantErr: nil,
		},
		{
			name: "bad target",
			courier: func() *Courier {
				loc, err := kernel.NewLocation(1, 1)
				assert.NoError(err)
				cur, err := NewCourier("test", 2, loc)
				assert.NoError(err)
				return cur
			}(),
			want:    0,
			wantErr: errs.ErrValueIsRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.courier.CalculateTimeToLocation(tt.target)
			if err != nil {
				assert.ErrorIs(err, tt.wantErr)
			} else {
				assert.Equal(tt.want, got)
			}
		})
	}
}

func TestCourier_Move(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name    string
		courier *Courier
		target  kernel.Location
		want    error
	}{
		{
			name: "good",
			courier: func() *Courier {
				loc, err := kernel.NewLocation(1, 1)
				assert.NoError(err)
				cur, err := NewCourier("test", 2, loc)
				assert.NoError(err)
				return cur
			}(),
			target: func() kernel.Location {
				loc, err := kernel.NewLocation(5, 5)
				assert.NoError(err)
				return loc
			}(),
			want: nil,
		},
		{
			name: "bad target",
			courier: func() *Courier {
				loc, err := kernel.NewLocation(1, 1)
				assert.NoError(err)
				cur, err := NewCourier("test", 2, loc)
				assert.NoError(err)
				return cur
			}(),
			want: errs.ErrValueIsInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dt1, err := tt.courier.Location().DistanceTo(tt.target)
			if err != nil {
				assert.ErrorIs(err, tt.want)
				return
			}

			if err := tt.courier.Move(tt.target); err != nil {
				assert.ErrorIs(err, tt.want)
			} else {
				dt2, err := tt.courier.Location().DistanceTo(tt.target)
				assert.NoError(err)
				assert.LessOrEqual(dt1-dt2, tt.courier.Speed())
			}
		})
	}
}
