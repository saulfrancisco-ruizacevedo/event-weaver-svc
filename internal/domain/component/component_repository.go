package component

import (
	"context"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/responses"
)

type IComponentRepository interface {
	GetAllComponentNames(ctx context.Context) ([]*Component, error)
	SaveAll(ctx context.Context, components []*Component) error
	CreateTeamRelationships(ctx context.Context, components []*Component) error
	GetAllComponents(ctx context.Context) (*responses.GraphResponseDto, error)
	GetComponent(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetComponentWithProducedEvents(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetComponentWithConsumedEvents(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetComponentWithSubscribedTopics(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetComponentWithProducedTopics(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetComponentWithTeam(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetComponentWithProducedAndConsumedEvents(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetComponentWithEventsAndSubscriptions(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetComponentWithEventsAndTopicRelations(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetComponentWithEventsAndTeam(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetComponentWithEventsAndTopicsAndTeams(ctx context.Context, name string) (*responses.GraphResponseDto, error)
}
