package event

import (
	"context"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/responses"
)

type IEventRepository interface {
	FindAllNamesByNamesIn(names []*string, ctx context.Context) ([]*EventSpecification, error)
	SaveAll(ctx context.Context, events []*EventSpecification) error
	CreateProducerRelationships(ctx context.Context, events []*EventSpecification) error
	CreateConsumerRelationships(ctx context.Context, events []*EventSpecification) error
	CreateTopicRelationships(ctx context.Context, events []*EventSpecification) error
	CreateRelatedEventRelationships(ctx context.Context, events []*EventSpecification) error
	CreateDomainRelationships(ctx context.Context, events []*EventSpecification) error
	GetAllEvents(ctx context.Context) (*responses.GraphResponseDto, error)

	GetEvent(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetEventWithEvents(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetEventWithEventsAndProduced(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetEventWithEventsProducedAndConsumed(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetEventWithEventsProducedConsumedAndTopic(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetEventWithEventsAndComponentsAndTopicsAndDomain(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetEventProducedByComponents(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetEventConsumedByComponents(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetEventOriginatedFromTopic(ctx context.Context, name string) (*responses.GraphResponseDto, error)
	GetEventBelongsToDomain(ctx context.Context, name string) (*responses.GraphResponseDto, error)
}
