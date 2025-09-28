package topic

import (
	"context"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/responses"
)

type ITopicRepository interface {
	GetAllTopicNames(ctx context.Context) ([]*Topic, error)
	SaveAll(ctx context.Context, topics []*Topic) error
	GetAllTopics(ctx context.Context) (*responses.GraphResponseDto, error)
	GetTopic(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetTopicWithEvents(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetTopicWithSubscribedComponents(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetTopicWithProducedComponents(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetTopicWithEventsAndSubscribedComponents(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetTopicWithEventsAndProducedComponents(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetTopicWithSubscribedAndProducedComponents(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetTopicWithEventsAndComponents(ctx context.Context, name string) (*responses.GraphResponseDto, error)
}
