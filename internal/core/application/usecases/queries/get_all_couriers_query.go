package queries

type GetAllCouriersQuery struct{ valid bool }

func NewGetAllCouriersQuery() (GetAllCouriersQuery, error) {
	return GetAllCouriersQuery{valid: true}, nil
}

func (q GetAllCouriersQuery) IsValid() bool { return q.valid }
