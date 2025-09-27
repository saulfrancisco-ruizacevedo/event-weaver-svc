package handlers

import (
	"context"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/commands"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/converters"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/component"
	domainentity "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/domainentity"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/event"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/team"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/topic"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/mediator"
)

type SpecPersistenceCommandHandler struct {
	teamRepository      team.ITeamRepository
	componentRepository component.IComponentRepository
	domainRepository    domainentity.IDomainRepository
	topicRepository     topic.ITopicRepository
	eventRepository     event.IEventRepository
}

var _ mediator.MediatorHandler[*commands.SpecPersistenceCommand, any] = &SpecPersistenceCommandHandler{}

func NewSpecPersistenceCommandHandler(teamRepository team.ITeamRepository,
	componentRepository component.IComponentRepository,
	domainRepository domainentity.IDomainRepository,
	topicRepository topic.ITopicRepository,
	eventRepository event.IEventRepository) *SpecPersistenceCommandHandler {
	return &SpecPersistenceCommandHandler{
		teamRepository:      teamRepository,
		componentRepository: componentRepository,
		domainRepository:    domainRepository,
		topicRepository:     topicRepository,
		eventRepository:     eventRepository,
	}
}

func (h *SpecPersistenceCommandHandler) Handle(ctx context.Context, command *commands.SpecPersistenceCommand) (any, error) {
	domains := converters.DomainDtoListToDomainList(command.Payload.Domains)
	teams := converters.TeamDtoListToTeamList(command.Payload.Teams)
	components := converters.ComponentDtoListToComponentList(command.Payload.Components)
	topics := converters.TopicDtoListToTopicList(command.Payload.Topics)
	events, _ := converters.EventDtoListToEventSpecificationList(command.Payload.Events)

	err := RunSteps(
		func() error { return h.domainRepository.SaveAll(ctx, domains) },
		func() error { return h.teamRepository.SaveAll(ctx, teams) },
		func() error { return h.componentRepository.SaveAll(ctx, components) },
		func() error { return h.componentRepository.CreateTeamRelationships(ctx, components) },
		func() error { return h.topicRepository.SaveAll(ctx, topics) },
		func() error { return h.eventRepository.SaveAll(ctx, events) },
		func() error { return h.eventRepository.CreateProducerRelationships(ctx, events) },
		func() error { return h.eventRepository.CreateConsumerRelationships(ctx, events) },
		func() error { return h.eventRepository.CreateTopicRelationships(ctx, events) },
		func() error { return h.eventRepository.CreateRelatedEventRelationships(ctx, events) },
		func() error { return h.eventRepository.CreateDomainRelationships(ctx, events) },
	)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func RunSteps(steps ...func() error) error {
	for _, step := range steps {
		if err := step(); err != nil {
			return err
		}
	}
	return nil
}
