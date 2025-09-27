package repositories

import (
	"context"
	"fmt"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/topic"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/outbound/neo4j/mappers"
)

type TopicRepository struct {
	BaseRepository *BaseRepository
}

func NewTopicRepository(baseRepository *BaseRepository) topic.ITopicRepository {
	return &TopicRepository{
		BaseRepository: baseRepository,
	}
}

var _ topic.ITopicRepository = &TopicRepository{}

func (r *TopicRepository) GetAllTopicNames(ctx context.Context) ([]*topic.Topic, error) {
	query := `
		MATCH (t:Topic)
		RETURN t.name AS name
	`

	records, err := r.BaseRepository.ExecQuery(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve topic names: %w", err)
	}

	return mappers.TopicNamesListToTopicList(records), nil
}

func (r *TopicRepository) SaveAll(ctx context.Context, topics []*topic.Topic) error {
	if len(topics) == 0 {
		return nil
	}

	query := `
		UNWIND $topics AS t
		MERGE (top:Topic {name: t.name})
		SET top.type = t.type
	`

	params := map[string]interface{}{
		"topics": make([]map[string]interface{}, len(topics)),
	}

	for i, topic := range topics {
		params["topics"].([]map[string]interface{})[i] = map[string]interface{}{
			"name": topic.Name,
			"type": topic.Type,
		}
	}

	_, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to save topics: %w", err)
	}

	return nil
}
