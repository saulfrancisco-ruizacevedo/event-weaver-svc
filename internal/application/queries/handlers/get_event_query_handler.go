package handlers

import (
	"context"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/responses"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/queries"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/event"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/mediator"
)

type GetEventQueryHandler struct {
	EventRepository event.IEventRepository
}

var _ mediator.MediatorHandler[*queries.GetEventQuery, *responses.GraphResponseDto] = &GetEventQueryHandler{}

func NewGetEventQueryHandler(EventRepository event.IEventRepository) *GetEventQueryHandler {
	return &GetEventQueryHandler{
		EventRepository: EventRepository,
	}
}

func (h *GetEventQueryHandler) Handle(ctx context.Context, query *queries.GetEventQuery) (*responses.GraphResponseDto, error) {
	dto := query.Payload

	switch {
	case dto.RelatedToEvent && dto.ProducedByComponent && dto.ConsumedByComponent && dto.OriginatesFromTopic && dto.BelongsToDomain:
		return h.EventRepository.GetEventWithEventsAndComponentsAndTopicsAndDomain(ctx, dto.Name)

	case dto.RelatedToEvent && dto.ProducedByComponent && dto.ConsumedByComponent && dto.OriginatesFromTopic:
		return h.EventRepository.GetEventWithEventsProducedConsumedAndTopic(ctx, dto.Name)

	case dto.RelatedToEvent && dto.ProducedByComponent && dto.ConsumedByComponent:
		return h.EventRepository.GetEventWithEventsProducedAndConsumed(ctx, dto.Name)

	case dto.RelatedToEvent && dto.ProducedByComponent:
		return h.EventRepository.GetEventWithEventsAndProduced(ctx, dto.Name)

	case dto.RelatedToEvent:
		return h.EventRepository.GetEventWithEvents(ctx, dto.Name)

	case dto.ProducedByComponent:
		return h.EventRepository.GetEventProducedByComponents(ctx, dto.Name)

	case dto.ConsumedByComponent:
		return h.EventRepository.GetEventConsumedByComponents(ctx, dto.Name)

	case dto.OriginatesFromTopic:
		return h.EventRepository.GetEventOriginatedFromTopic(ctx, dto.Name)

	case dto.BelongsToDomain:
		return h.EventRepository.GetEventBelongsToDomain(ctx, dto.Name)

	default:
		return h.EventRepository.GetEvent(ctx, dto.Name)
	}
}
