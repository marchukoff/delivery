package commands

type MoveCouriersCommand struct{ valid bool }

func NewMoveCouriersCommand() (MoveCouriersCommand, error) {
	return MoveCouriersCommand{valid: true}, nil
}

func (c MoveCouriersCommand) IsValid() bool { return c.valid }
