package queries

type GetIncompleteOrdersQuery struct{ valid bool }

func NewGetIncompleteOrdersQuery() (GetIncompleteOrdersQuery, error) {
	return GetIncompleteOrdersQuery{valid: true}, nil
}

func (q GetIncompleteOrdersQuery) IsValid() bool { return q.valid }
