package repositories

import (
	"context"
	"fmt"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/topic"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/outbound/neo4j/mappers"
	localModels "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/outbound/neo4j/models"
	"github.com/saulfrancisco-ruizacevedo/go-neopersist"
	"github.com/saulfrancisco-ruizacevedo/gocypher"
)

type TopicRepository struct {
	manager   *neopersist.PersistenceManager
	topicRepo *neopersist.Repository[localModels.Topic]
}

func NewTopicRepository(manager *neopersist.PersistenceManager) (topic.ITopicRepository, error) {
	topicRepo, err := neopersist.RepositoryFor[localModels.Topic](manager)
	if err != nil {
		return nil, fmt.Errorf("failed to create topic repository: %w", err)
	}
	return &TopicRepository{
		manager:   manager,
		topicRepo: topicRepo,
	}, nil
}

var _ topic.ITopicRepository = &TopicRepository{}

func (r *TopicRepository) GetAllTopicNames(ctx context.Context) ([]*topic.Topic, error) {
	qb := gocypher.NewQueryBuilder().
		Match(gocypher.N("t", "Topic")).
		Return("t.name AS name")

	modelTopics, err := r.topicRepo.Find(ctx, qb)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve topic names: %w", err)
	}

	return mappers.ModelTopicListToDomainList(modelTopics), nil
}

func (r *TopicRepository) SaveAll(ctx context.Context, topics []*topic.Topic) error {
	if len(topics) == 0 {
		return nil
	}

	modelTopics := mappers.DomainTopicListToModelList(topics)

	return r.topicRepo.SaveAll(ctx, modelTopics)
}
