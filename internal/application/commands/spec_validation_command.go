package commands

import (
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/requests"
)

const SpecValidationCommandName = "SpecValidationCommand"

type SpecValidationCommand struct {
	Payload requests.EventSpecificationRequestDto
}

func NewSpecValidationCommand(payload requests.EventSpecificationRequestDto) *SpecValidationCommand {
	return &SpecValidationCommand{
		Payload: payload,
	}
}
