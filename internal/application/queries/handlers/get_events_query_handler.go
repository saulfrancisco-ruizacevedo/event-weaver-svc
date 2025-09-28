package handlers

import (
	"context"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/responses"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/queries"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/event"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/mediator"
)

type GetEventsQueryHandler struct {
	eventsRepository event.IEventRepository
}

var _ mediator.MediatorHandler[*queries.GetEventsQuery, *responses.GraphResponseDto] = &GetEventsQueryHandler{}

func NewGetEventsQueryHandler(eventsRepository event.IEventRepository) *GetEventsQueryHandler {
	return &GetEventsQueryHandler{
		eventsRepository: eventsRepository,
	}
}

func (h *GetEventsQueryHandler) Handle(ctx context.Context, query *queries.GetEventsQuery) (*responses.GraphResponseDto, error) {
	if graph, err := h.eventsRepository.GetAllEvents(ctx); err != nil {
		return nil, err

	} else {
		return graph, nil
	}
}
