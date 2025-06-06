package kernel_test

import (
	"errors"
	"testing"

	. "delivery/internal/core/domain/kernel"
	"delivery/internal/pkg/errs"
)

func TestNewLocation(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name    string
		args    args
		want    Location
		wantErr error
	}{
		{
			name: "good",
			args: args{2, 4},
			want: func() Location {
				loc, _ := NewLocation(2, 4)
				return loc
			}(),
			wantErr: nil,
		},
		{
			name: "bad",
			args: args{0, 4},
			want: func() Location {
				loc, _ := NewLocation(0, 4)
				return loc
			}(),
			wantErr: errs.ErrValueIsOutOfRange,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLocation(tt.args.x, tt.args.y)
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("NewLocation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !got.Equals(tt.want) {
				t.Errorf("NewLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRandomLocation(t *testing.T) {
	got := NewRandomLocation()
	if got.IsEmpty() {
		t.Errorf("NewRandomLocation() = %v, want non empty Location", got)
	}
}

func TestLocation_DistanceTo(t *testing.T) {
	tests := []struct {
		name    string
		src     Location
		dst     Location
		want    int
		wantErr error
	}{
		{
			name: "good",
			src: func() Location {
				loc, _ := NewLocation(2, 6)
				return loc
			}(),
			dst: func() Location {
				loc, _ := NewLocation(4, 9)
				return loc
			}(),
			want:    5,
			wantErr: nil,
		},
		{
			name: "bad source",
			src: func() Location {
				loc, _ := NewLocation(-1, 6)
				return loc
			}(),
			dst: func() Location {
				loc, _ := NewLocation(4, 9)
				return loc
			}(),
			want:    0,
			wantErr: errs.ErrValueIsInvalid,
		},
		{
			name: "bad target",
			src: func() Location {
				loc, _ := NewLocation(2, 6)
				return loc
			}(),
			dst: func() Location {
				loc, _ := NewLocation(4, 90)
				return loc
			}(),
			want:    0,
			wantErr: errs.ErrValueIsInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.src.DistanceTo(tt.dst)
			if got < 0 {
				t.Errorf("Location.DistanceTo() = %v, it should not be less than 0", got)
			}

			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("Location.DistanceTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("Location.DistanceTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func FuzzNewRandomLocation(f *testing.F) {
	for i := range 10 {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		got := NewRandomLocation()
		if got.IsEmpty() {
			t.Errorf("%d: got empty Location", i)
		}
	})
}

func TestLocation_Equals(t *testing.T) {
	tests := []struct {
		name       string
		loc1, loc2 Location
		want       bool
	}{
		{
			name: "equals",
			loc1: func() Location {
				loc, _ := NewLocation(1, 1)
				return loc
			}(),
			loc2: func() Location {
				loc, _ := NewLocation(1, 1)
				return loc
			}(),
			want: true,
		},
		{
			name: "not equals",
			loc1: func() Location {
				loc, _ := NewLocation(1, 1)
				return loc
			}(),
			loc2: func() Location {
				loc, _ := NewLocation(10, 10)
				return loc
			}(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.loc1.Equals(tt.loc2) != tt.want {
				t.Errorf("Location.Equals: %v = %v, want %v", tt.loc1, tt.loc2, tt.want)
			}
		})
	}
}
