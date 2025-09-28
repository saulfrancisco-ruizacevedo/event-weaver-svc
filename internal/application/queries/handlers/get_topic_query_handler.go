package handlers

import (
	"context"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/responses"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/queries"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/topic"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/mediator"
)

type GetTopicQueryHandler struct {
	TopicRepository topic.ITopicRepository
}

var _ mediator.MediatorHandler[*queries.GetTopicQuery, *responses.GraphResponseDto] = &GetTopicQueryHandler{}

func NewGetTopicQueryHandler(TopicRepository topic.ITopicRepository) *GetTopicQueryHandler {
	return &GetTopicQueryHandler{
		TopicRepository: TopicRepository,
	}
}

func (h *GetTopicQueryHandler) Handle(ctx context.Context, query *queries.GetTopicQuery) (*responses.GraphResponseDto, error) {
	dto := query.Payload

	switch {
	case dto.EventsOriginated && dto.ComponentsSubscribed && dto.ComponentsProduced:
		return h.TopicRepository.GetTopicWithEventsAndComponents(ctx, dto.Name)
	case dto.EventsOriginated && dto.ComponentsSubscribed:
		return h.TopicRepository.GetTopicWithEventsAndSubscribedComponents(ctx, dto.Name)
	case dto.EventsOriginated && dto.ComponentsProduced:
		return h.TopicRepository.GetTopicWithEventsAndProducedComponents(ctx, dto.Name)
	case dto.ComponentsSubscribed && dto.ComponentsProduced:
		return h.TopicRepository.GetTopicWithSubscribedAndProducedComponents(ctx, dto.Name)
	case dto.EventsOriginated:
		return h.TopicRepository.GetTopicWithEvents(ctx, dto.Name)
	case dto.ComponentsSubscribed:
		return h.TopicRepository.GetTopicWithSubscribedComponents(ctx, dto.Name)
	case dto.ComponentsProduced:
		return h.TopicRepository.GetTopicWithProducedComponents(ctx, dto.Name)
	default:
		return h.TopicRepository.GetTopic(ctx, dto.Name)
	}
}
