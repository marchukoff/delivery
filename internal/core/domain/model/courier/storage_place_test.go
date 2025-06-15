package courier_test

import (
	"testing"

	. "delivery/internal/core/domain/model/courier"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStoragePlace(t *testing.T) {
	tests := []struct {
		name        string
		place       string
		totalVolume int
		want        error
	}{
		{name: "good", place: "Bag", totalVolume: 10, want: nil},
		{name: "bad name", place: "", totalVolume: 10, want: errs.ErrValueIsRequired},
		{name: "bad volume", place: "Bag", totalVolume: 0, want: errs.ErrValueIsInvalid},
	}

	assert := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp, err := NewStoragePlace(tt.place, tt.totalVolume)
			if err != nil {
				assert.ErrorIs(err, tt.want)
			}

			assert.Empty(sp.OrderID())

			if err != nil {
				assert.Empty(sp.ID())
				assert.Equal(uuid.Nil, sp.ID())
				assert.Empty(sp.Name())
				assert.Empty(sp.TotalVolume())
			} else {
				assert.NotEqual(uuid.Nil, sp.ID())
				assert.NotEmpty(sp.Name())
				assert.NotEmpty(sp.TotalVolume())
			}
		})
	}
}

func TestStoragePlace_Store(t *testing.T) {
	tests := []struct {
		name     string
		place    *StoragePlace
		orderID uuid.UUID
		volume   int
		want     error
	}{
		{
			name:     "good",
			place:    func() *StoragePlace { x, _ := NewStoragePlace("bag", 10); return x }(),
			orderID: uuid.New(),
			volume:   1,
			want:     nil,
		},
		{
			name:     "bad nil StoragePlace",
			orderID: uuid.New(),
			volume:   1,
			want:     ErrStoragePlaceNotInitialized,
		},
		{
			name:     "bad overweight",
			place:    func() *StoragePlace { x, _ := NewStoragePlace("bag", 1); return x }(),
			orderID: uuid.New(),
			volume:   10,
			want:     errs.ErrValueIsOutOfRange,
		},
		{
			name:     "bad order uuid",
			place:    func() *StoragePlace { x, _ := NewStoragePlace("bag", 10); return x }(),
			orderID: uuid.Nil,
			volume:   1,
			want:     errs.ErrValueIsInvalid,
		},
		{
			name:     "bad occupied",
			place:    func() *StoragePlace { x, _ := NewStoragePlace("bag", 10); x.Store(uuid.New(), 1); return x }(),
			orderID: uuid.New(),
			volume:   10,
			want:     ErrStoragePlaceIsOccupied,
		},
	}

	assert := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.place.Store(tt.orderID, tt.volume); err != nil {
				assert.ErrorIs(err, tt.want)
			} else {
				assert.NotEmpty(tt.place.OrderID())
			}
		})
	}
}

func TestStoragePlace_Equals(t *testing.T) {
	require := require.New(t)
	sp1, err := NewStoragePlace("first", 10)
	require.NoError(err)

	sp2, err := NewStoragePlace("second", 10)
	require.NoError(err)

	tests := []struct {
		name     string
		sp1, sp2 *StoragePlace
		want     bool
	}{
		{name: "sp1 = sp1", sp1: sp1, sp2: sp1, want: true},
		{name: "sp1 = sp2", sp1: sp1, sp2: sp2, want: false},
		{name: "nil = sp1", sp1: nil, sp2: sp1, want: false},
		{name: "sp2 = nil", sp1: sp2, sp2: nil, want: false},
	}

	assert := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(tt.want, tt.sp1.Equals(tt.sp2))
		})
	}
}

func TestStoragePlace_CanStore(t *testing.T) {
	tests := []struct {
		name    string
		place   *StoragePlace
		volume  int
		want    bool
		wantErr error
	}{
		{
			name:    "good ok",
			place:   func() *StoragePlace { x, _ := NewStoragePlace("bag", 10); return x }(),
			volume:  1,
			want:    true,
			wantErr: nil,
		},
		{
			name:    "good not ok",
			place:   func() *StoragePlace { x, _ := NewStoragePlace("bag", 1); return x }(),
			volume:  10,
			want:    false,
			wantErr: nil,
		},
		{
			name:    "good occupied",
			place:   func() *StoragePlace { x, _ := NewStoragePlace("bag", 10); x.Store(uuid.New(), 1); return x }(),
			volume:  1,
			want:    false,
			wantErr: nil,
		},
		{
			name:    "bad volume invalid",
			place:   func() *StoragePlace { x, _ := NewStoragePlace("bag", 10); return x }(),
			volume:  -1,
			want:    false,
			wantErr: errs.ErrValueIsInvalid,
		},
		{
			name:    "bad nil StoragePlace",
			place:   nil,
			volume:  1,
			want:    false,
			wantErr: ErrStoragePlaceNotInitialized,
		},
	}

	assert := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.place.CanStore(tt.volume)
			if err != nil {
				assert.ErrorIs(err, tt.wantErr)
			}
			assert.Equal(tt.want, got)
		})
	}
}

func TestStoragePlace_Clear(t *testing.T) {
	id := uuid.New()
	tests := []struct {
		name     string
		place    *StoragePlace
		orderID uuid.UUID
		want     error
	}{
		{
			name:     "good",
			place:    func() *StoragePlace { x, _ := NewStoragePlace("bag", 10); x.Store(id, 1); return x }(),
			orderID: id,
			want:     nil,
		},
		{
			name:     "bad not found",
			place:    func() *StoragePlace { x, _ := NewStoragePlace("bag", 10); x.Store(id, 1); return x }(),
			orderID: uuid.New(),
			want:     errs.ErrObjectNotFound,
		},
		{
			name:     "bad place",
			place:    new(StoragePlace),
			orderID: id,
			want:     errs.ErrObjectNotFound,
		},
		{
			name:     "bad nil place",
			place:    nil,
			orderID: id,
			want:     ErrStoragePlaceNotInitialized,
		},
	}

	assert := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.place.Clear(tt.orderID); err != nil {
				assert.ErrorIs(err, tt.want)
			} else {
				assert.Empty(tt.place.OrderID())
				assert.NoError(tt.place.Store(tt.orderID, 1))
			}
		})
	}
}
