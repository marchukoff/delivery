package commands

import (
	"strings"

	"delivery/internal/pkg/errs"
)

type CreateCourierCommand struct {
	name  string
	speed int
	valid bool
}

func NewCreateCourierCommand(name string, speed int) (CreateCourierCommand, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return CreateCourierCommand{}, errs.NewValueIsRequiredError("name")
	}

	if speed <= 0 {
		return CreateCourierCommand{}, errs.NewValueIsRequiredError("speed")
	}

	return CreateCourierCommand{
		name:  name,
		speed: speed,
		valid: true,
	}, nil
}

func (c CreateCourierCommand) Name() string { return c.name }

func (c CreateCourierCommand) Speed() int { return c.speed }

func (c CreateCourierCommand) IsValid() bool { return c.valid }
