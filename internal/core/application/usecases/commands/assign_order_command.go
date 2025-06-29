package commands

import (
	"github.com/google/uuid"
)

type AssignOrderCommand struct{ valid bool }

func NewAssignOrderCommand(orderID uuid.UUID, street string, volume int) (AssignOrderCommand, error) {
	return AssignOrderCommand{valid: true}, nil
}

func (c AssignOrderCommand) IsValid() bool { return c.valid }
