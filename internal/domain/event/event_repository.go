package event

import (
	"context"
)

type IEventRepository interface {
	FindAllNamesByNamesIn(names []*string, ctx context.Context) ([]*EventSpecification, error)
	SaveAll(ctx context.Context, events []*EventSpecification) error
	CreateProducerRelationships(ctx context.Context, events []*EventSpecification) error
	CreateConsumerRelationships(ctx context.Context, events []*EventSpecification) error
	CreateTopicRelationships(ctx context.Context, events []*EventSpecification) error
	CreateRelatedEventRelationships(ctx context.Context, events []*EventSpecification) error
	CreateDomainRelationships(ctx context.Context, events []*EventSpecification) error
}
