package kernel

import (
	"crypto/rand"
	"errors"
	"math/big"

	"delivery/internal/pkg/errs"
)

const (
	minCoord = 1
	maxCoord = 10
)

type Location struct {
	x, y  int
	valid bool
}

func NewLocation(x, y int) (Location, error) {
	if x < minCoord || x > maxCoord {
		return Location{}, errs.NewValueIsOutOfRangeError("Location.x", x, minCoord, maxCoord)
	}

	if y < minCoord || y > maxCoord {
		return Location{}, errs.NewValueIsOutOfRangeError("Location.y", y, minCoord, maxCoord)
	}

	return Location{x: x, y: y, valid: true}, nil
}

func NewRandomLocation() Location {
	rnd := func() int {
		n, _ := rand.Int(rand.Reader, big.NewInt(maxCoord-minCoord))
		return int(n.Int64()) + minCoord
	}

	res, err := NewLocation(rnd(), rnd())
	if err != nil {
		panic(err) // should never happen
	}

	return res
}

func (l Location) DistanceTo(target Location) (int, error) {
	if !l.valid {
		cause := errors.New("source location not initialized")
		return 0, errs.NewValueIsInvalidErrorWithCause("Location", cause)
	}

	if !target.valid {
		cause := errors.New("target location not initialized")
		return 0, errs.NewValueIsInvalidErrorWithCause("Location", cause)
	}

	x1, x2 := max(l.x, target.x), min(l.x, target.x)
	y1, y2 := max(l.y, target.y), min(l.y, target.y)

	return (x1 - x2) + (y1 - y2), nil
}

func (l Location) Equals(other Location) bool { return l == other }

func (l Location) IsValid() bool { return l.valid }

func (l Location) X() int { return l.x }

func (l Location) Y() int { return l.y }
