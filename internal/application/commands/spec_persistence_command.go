package commands

import "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/requests"

const SpecPersictenceCommandName = "SpecPersistenceCommand"

type SpecPersistenceCommand struct {
	Payload requests.EventSpecificationRequestDto
}

func NewSpecPersistenceCommand(payload requests.EventSpecificationRequestDto) *SpecPersistenceCommand {
	return &SpecPersistenceCommand{
		Payload: payload,
	}
}
