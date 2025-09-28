package handlers

import (
	"context"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/responses"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/queries"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/team"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/mediator"
)

type GetTeamsQueryHandler struct {
	teamsRepository team.ITeamRepository
}

var _ mediator.MediatorHandler[*queries.GetTeamsQuery, *responses.GraphResponseDto] = &GetTeamsQueryHandler{}

func NewGetTeamsQueryHandler(teamsRepository team.ITeamRepository) *GetTeamsQueryHandler {
	return &GetTeamsQueryHandler{
		teamsRepository: teamsRepository,
	}
}

func (h *GetTeamsQueryHandler) Handle(ctx context.Context, query *queries.GetTeamsQuery) (*responses.GraphResponseDto, error) {
	if graph, err := h.teamsRepository.GetAllTeams(ctx); err != nil {
		return nil, err

	} else {
		return graph, nil
	}
}
