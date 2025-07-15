package courierrepo

import (
	"context"

	"delivery/internal/adapters/out/postgres/shared"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ ports.CourierRepository = &Repository{}

type Repository struct {
	tracker shared.Tracker
}

func NewRepository(tracker shared.Tracker) (ports.CourierRepository, error) {
	if tracker == nil {
		return nil, errs.NewValueIsRequiredError("tracker")
	}
	return &Repository{tracker: tracker}, nil
}

func (r *Repository) Add(ctx context.Context, aggregate *courier.Courier) error {
	r.tracker.Track(aggregate)
	dto := DomainToDTO(aggregate)

	isInTransaction := r.tracker.InTx()
	if !isInTransaction {
		r.tracker.Begin(ctx)
	}
	tx := r.tracker.Tx()

	err := tx.WithContext(ctx).
		Session(&gorm.Session{FullSaveAssociations: true}).
		Create(&dto).
		Error
	if err != nil {
		return err
	}

	if !isInTransaction {
		err := r.tracker.Commit(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, ID uuid.UUID) (*courier.Courier, error) {
	var dto CourierDTO

	tx := r.getTxOrDb()
	res := tx.WithContext(ctx).
		Preload(clause.Associations).
		Find(&dto, ID)
	if res.RowsAffected == 0 {
		return nil, errs.NewObjectNotFoundError("courier.id", ID)
	}

	return DtoToDomain(dto), nil
}

func (r *Repository) GetAllFree(ctx context.Context) ([]*courier.Courier, error) {
	var dtos []CourierDTO

	tx := r.getTxOrDb()
	res := tx.WithContext(ctx).
		Preload(clause.Associations).
		Where(`
        NOT EXISTS (
            SELECT 1 FROM storage_places sp
            WHERE sp.courier_id = couriers.id AND sp.order_id IS NOT NULL
        )`).
		Find(&dtos)
	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, errs.NewObjectNotFoundError("Free couriers", nil)
	}

	couriers := make([]*courier.Courier, 0, len(dtos))
	for _, dto := range dtos {
		couriers = append(couriers, DtoToDomain(dto))
	}

	return couriers, nil
}

func (r *Repository) Update(ctx context.Context, aggregate *courier.Courier) error {
	r.tracker.Track(aggregate)
	dto := DomainToDTO(aggregate)

	isInTransaction := r.tracker.InTx()
	if !isInTransaction {
		r.tracker.Begin(ctx)
	}
	tx := r.tracker.Tx()

	err := tx.WithContext(ctx).
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(&dto).
		Error
	if err != nil {
		return err
	}

	if !isInTransaction {
		err := r.tracker.Commit(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) getTxOrDb() *gorm.DB {
	if tx := r.tracker.Tx(); tx != nil {
		return tx
	}
	return r.tracker.Db()
}
