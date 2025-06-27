package courierrepo

import (
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/kernel"

	"github.com/google/uuid"
)

func DomainToDTO(courier *courier.Courier) CourierDTO {
	places := make([]*StoragePlaceDTO, 0, len(courier.StoragePlaces()))
	for _, v := range courier.StoragePlaces() {
		oid := uuid.Nil
		if v.OrderID() != nil {
			oid = *v.OrderID()
		}

		places = append(places, &StoragePlaceDTO{
			ID:          v.ID(),
			Name:        v.Name(),
			TotalVolume: v.TotalVolume(),
			OrderID:     oid,
		})
	}

	return CourierDTO{
		ID:    courier.ID(),
		Name:  courier.Name(),
		Speed: courier.Speed(),
		Location: LocationDTO{
			X: courier.Location().X(),
			Y: courier.Location().Y(),
		},
		StoragePlaces: places,
	}
}

func DtoToDomain(dto CourierDTO) *courier.Courier {
	places := make([]*courier.StoragePlace, 0, len(dto.StoragePlaces))
	for _, place := range dto.StoragePlaces {
		places = append(places, courier.RestoreStoragePlace(place.ID, place.Name, place.TotalVolume, place.OrderID))
	}

	loc, _ := kernel.NewLocation(dto.Location.X, dto.Location.Y)
	return courier.RestoreCourier(dto.ID, dto.Name, dto.Speed, loc, places)
}
