package handlers

import (
	"context"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/responses"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/queries"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/team"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/mediator"
)

type GetTeamQueryHandler struct {
	TeamRepository team.ITeamRepository
}

var _ mediator.MediatorHandler[*queries.GetTeamQuery, *responses.GraphResponseDto] = &GetTeamQueryHandler{}

func NewGetTeamQueryHandler(TeamRepository team.ITeamRepository) *GetTeamQueryHandler {
	return &GetTeamQueryHandler{
		TeamRepository: TeamRepository,
	}
}

func (h *GetTeamQueryHandler) Handle(ctx context.Context, query *queries.GetTeamQuery) (*responses.GraphResponseDto, error) {
	dto := query.Payload

	if dto.ManagesComponent && dto.ProducesEvent && dto.BelongsToDomain {
		return h.TeamRepository.GetTeamWithComponentsAndEventsAndDomains(ctx, dto.Name)

	} else if dto.ManagesComponent && dto.ProducesEvent {
		return h.TeamRepository.GetTeamWithComponentsAndEvents(ctx, dto.Name)

	} else if dto.ManagesComponent {
		return h.TeamRepository.GetTeamWithComponents(ctx, dto.Name)

	} else {
		return h.TeamRepository.GetTeam(ctx, dto.Name)
	}
}
