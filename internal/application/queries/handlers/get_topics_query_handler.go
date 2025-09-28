package handlers

import (
	"context"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/responses"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/queries"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/topic"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/mediator"
)

type GetTopicsQueryHandler struct {
	topicsRepository topic.ITopicRepository
}

var _ mediator.MediatorHandler[*queries.GetTopicsQuery, *responses.GraphResponseDto] = &GetTopicsQueryHandler{}

func NewGetTopicsQueryHandler(topicsRepository topic.ITopicRepository) *GetTopicsQueryHandler {
	return &GetTopicsQueryHandler{
		topicsRepository: topicsRepository,
	}
}

func (h *GetTopicsQueryHandler) Handle(ctx context.Context, query *queries.GetTopicsQuery) (*responses.GraphResponseDto, error) {
	if graph, err := h.topicsRepository.GetAllTopics(ctx); err != nil {
		return nil, err

	} else {
		return graph, nil
	}
}
