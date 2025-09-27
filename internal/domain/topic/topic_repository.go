package topic

import "context"

type ITopicRepository interface {
	GetAllTopicNames(ctx context.Context) ([]*Topic, error)
	SaveAll(ctx context.Context, topics []*Topic) error
}
