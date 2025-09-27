package handlers

import (
	"context"
	"fmt"
	"sync"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/commands"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/common"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/converters"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/responses"
	domainCommon "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/common"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/component"
	domainentity "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/domainentity"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/event"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/spec"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/team"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/topic"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/mediator"
)

type SpecValidationCommandHandler struct {
	teamRepository      team.ITeamRepository
	componentRepository component.IComponentRepository
	domainRepository    domainentity.IDomainRepository
	topicRepository     topic.ITopicRepository
	eventRepository     event.IEventRepository
}

var _ mediator.MediatorHandler[*commands.SpecValidationCommand, *responses.EventSpecValidationResponseDto] = &SpecValidationCommandHandler{}

func NewSpecValidationCommandHandler(teamRepository team.ITeamRepository,
	componentRepository component.IComponentRepository,
	domainRepository domainentity.IDomainRepository,
	topicRepository topic.ITopicRepository,
	eventRepository event.IEventRepository) *SpecValidationCommandHandler {
	return &SpecValidationCommandHandler{
		teamRepository:      teamRepository,
		componentRepository: componentRepository,
		domainRepository:    domainRepository,
		topicRepository:     topicRepository,
		eventRepository:     eventRepository,
	}
}

func (h *SpecValidationCommandHandler) Handle(ctx context.Context, command *commands.SpecValidationCommand) (*responses.EventSpecValidationResponseDto, error) {
	domains := converters.DomainDtoListToDomainList(command.Payload.Domains)
	teams := converters.TeamDtoListToTeamList(command.Payload.Teams)
	components := converters.ComponentDtoListToComponentList(command.Payload.Components)
	topics := converters.TopicDtoListToTopicList(command.Payload.Topics)
	events, err := converters.EventDtoListToEventSpecificationList(command.Payload.Events)
	if err != nil {
		return nil, fmt.Errorf("error mapping events: %w", err)
	}

	dbDomainsCh := make(chan []*domainentity.Domain, 1)
	dbTeamsCh := make(chan []*team.Team, 1)
	dbComponentsCh := make(chan []*component.Component, 1)
	dbTopicsCh := make(chan []*topic.Topic, 1)
	dbEventsCh := make(chan []*event.EventSpecification, 1)
	errCh := make(chan error, 5)

	var wg sync.WaitGroup
	wg.Add(5)

	go common.FetchWithChannel(ctx, h.domainRepository.GetAllDomainNames, dbDomainsCh, errCh, &wg)
	go common.FetchWithChannel(ctx, h.teamRepository.GetAllTeamNames, dbTeamsCh, errCh, &wg)
	go common.FetchWithChannel(ctx, h.componentRepository.GetAllComponentNames, dbComponentsCh, errCh, &wg)
	go common.FetchWithChannel(ctx, h.topicRepository.GetAllTopicNames, dbTopicsCh, errCh, &wg)
	go common.FetchWithChannel(ctx, func(ctx context.Context) ([]*event.EventSpecification, error) {
		var allNames []*string
		for _, ev := range events {
			allNames = append(allNames, ev.GetRelatedEventsNames()...)
		}
		return h.eventRepository.FindAllNamesByNamesIn(allNames, ctx)
	}, dbEventsCh, errCh, &wg)

	wg.Wait()
	close(errCh)

	if e := <-errCh; e != nil {
		return nil, e
	}

	dbDomains := <-dbDomainsCh
	dbTeams := <-dbTeamsCh
	dbComponents := <-dbComponentsCh
	dbTopics := <-dbTopicsCh
	dbEvents := <-dbEventsCh

	validator := spec.NewSpecValidator()

	domainResult := validator.ValidateDomains(domains)
	domains = domainCommon.MergeDistinctByName(domains, dbDomains)

	teamResult := validator.ValidateTeams(teams)
	teams = domainCommon.MergeDistinctByName(teams, dbTeams)

	componentResult := validator.ValidateComponents(components, teams)
	components = domainCommon.MergeDistinctByName(components, dbComponents)

	topicResult := validator.ValidateTopics(topics)
	topics = domainCommon.MergeDistinctByName(topics, dbTopics)

	relatedEvents := domainCommon.MergeDistinctByName(events, dbEvents)
	eventResult := validator.ValidateEvents(events, components, domains, topics, relatedEvents)

	return &responses.EventSpecValidationResponseDto{
		DomainResult:    domainResult,
		TeamResult:      teamResult,
		ComponentResult: componentResult,
		TopicResult:     topicResult,
		EventResult:     eventResult,
	}, nil
}
