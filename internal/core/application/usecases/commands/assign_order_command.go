package commands

type AssignOrderCommand struct{ valid bool }

func NewAssignOrderCommand() (AssignOrderCommand, error) {
	return AssignOrderCommand{valid: true}, nil
}

func (c AssignOrderCommand) IsValid() bool { return c.valid }
