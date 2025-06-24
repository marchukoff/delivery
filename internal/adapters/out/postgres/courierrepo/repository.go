package courierrepo

import (
	"context"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
	"github.com/google/uuid"
)

var _ ports.CourierRepository = &Repository{}

type Repository struct {
	tracker Tracker
}

// Add implements ports.CourierRepository.
func (r *Repository) Add(ctx context.Context, aggregate *courier.Courier) error {
	panic("unimplemented")
}

// Get implements ports.CourierRepository.
func (r *Repository) Get(ctx context.Context, ID uuid.UUID) (*courier.Courier, error) {
	panic("unimplemented")
}

// GetAllFree implements ports.CourierRepository.
func (r *Repository) GetAllFree(ctx context.Context) ([]*courier.Courier, error) {
	panic("unimplemented")
}

// Update implements ports.CourierRepository.
func (r *Repository) Update(ctx context.Context, aggregate *courier.Courier) error {
	panic("unimplemented")
}

func NewRepository(tracker Tracker) (*Repository, error) {
	if tracker == nil {
		return nil, errs.NewValueIsRequiredError("tracker")
	}

	return &Repository{
		tracker: tracker,
	}, nil
}
