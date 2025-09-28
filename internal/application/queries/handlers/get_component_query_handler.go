package handlers

import (
	"context"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/responses"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/queries"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/component"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/mediator"
)

type GetComponentQueryHandler struct {
	ComponentRepository component.IComponentRepository
}

var _ mediator.MediatorHandler[*queries.GetComponentQuery, *responses.GraphResponseDto] = &GetComponentQueryHandler{}

func NewGetComponentQueryHandler(ComponentRepository component.IComponentRepository) *GetComponentQueryHandler {
	return &GetComponentQueryHandler{
		ComponentRepository: ComponentRepository,
	}
}

func (h *GetComponentQueryHandler) Handle(ctx context.Context, query *queries.GetComponentQuery) (*responses.GraphResponseDto, error) {
	dto := query.Payload

	switch {
	case dto.ProducesEvent && dto.ConsumesEvent && dto.SubscribesToTopic && dto.ProducesToTopic && dto.ManagedByTeam:
		return h.ComponentRepository.GetComponentWithEventsAndTopicsAndTeams(ctx, dto.Name)
	case dto.ProducesEvent && dto.ConsumesEvent && dto.SubscribesToTopic && dto.ProducesToTopic:
		return h.ComponentRepository.GetComponentWithEventsAndTopicRelations(ctx, dto.Name)
	case dto.ProducesEvent && dto.ConsumesEvent && dto.SubscribesToTopic:
		return h.ComponentRepository.GetComponentWithEventsAndSubscriptions(ctx, dto.Name)
	case dto.ProducesEvent && dto.ConsumesEvent:
		return h.ComponentRepository.GetComponentWithProducedAndConsumedEvents(ctx, dto.Name)
	case dto.ProducesEvent && dto.ConsumesEvent && dto.ManagedByTeam:
		return h.ComponentRepository.GetComponentWithEventsAndTeam(ctx, dto.Name)
	case dto.ProducesEvent:
		return h.ComponentRepository.GetComponentWithProducedEvents(ctx, dto.Name)
	case dto.ConsumesEvent:
		return h.ComponentRepository.GetComponentWithConsumedEvents(ctx, dto.Name)
	case dto.SubscribesToTopic:
		return h.ComponentRepository.GetComponentWithSubscribedTopics(ctx, dto.Name)
	case dto.ProducesToTopic:
		return h.ComponentRepository.GetComponentWithProducedTopics(ctx, dto.Name)
	case dto.ManagedByTeam:
		return h.ComponentRepository.GetComponentWithTeam(ctx, dto.Name)
	default:
		return h.ComponentRepository.GetComponent(ctx, dto.Name)
	}
}
