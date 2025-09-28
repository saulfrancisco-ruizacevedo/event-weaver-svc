package team

import (
	"context"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/responses"
)

type ITeamRepository interface {
	GetAllTeamNames(ctx context.Context) ([]*Team, error)
	SaveAll(ctx context.Context, team []*Team) error
	GetAllTeams(ctx context.Context) (*responses.GraphResponseDto, error)
	GetTeamWithComponentsAndEventsAndDomains(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetTeamWithComponentsAndEvents(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetTeamWithComponents(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetTeam(ctx context.Context, name string) (*responses.GraphResponseDto, error)
}
