package handlers

import (
	"context"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/responses"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/queries"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/component"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/mediator"
)

type GetComponentsQueryHandler struct {
	componentRepository component.IComponentRepository
}

var _ mediator.MediatorHandler[*queries.GetComponentsQuery, *responses.GraphResponseDto] = &GetComponentsQueryHandler{}

func NewGetComponentsQueryHandler(componentRepository component.IComponentRepository) *GetComponentsQueryHandler {
	return &GetComponentsQueryHandler{
		componentRepository: componentRepository,
	}
}

func (h *GetComponentsQueryHandler) Handle(ctx context.Context, query *queries.GetComponentsQuery) (*responses.GraphResponseDto, error) {
	if graph, err := h.componentRepository.GetAllComponents(ctx); err != nil {
		return nil, err

	} else {
		return graph, nil
	}
}
