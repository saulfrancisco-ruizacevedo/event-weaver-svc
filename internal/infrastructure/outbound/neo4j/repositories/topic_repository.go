package repositories

import (
	"context"
	"fmt"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/responses"
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

func (r *TopicRepository) GetAllTopics(ctx context.Context) (*responses.GraphResponseDto, error) {
	return r.BaseRepository.GetAllNodes(ctx, "Topic")
}

func (r *TopicRepository) GetTopic(ctx context.Context, name string) (*responses.GraphResponseDto, error) {
	query := `MATCH (t:Topic {name: $name}) RETURN t`
	params := map[string]interface{}{"name": name}

	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "t")
	return &responses.GraphResponseDto{Nodes: nodes, Relationships: []responses.GraphRelationshipDto{}}, nil
}

func (r *TopicRepository) GetTopicWithEvents(ctx context.Context, name string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (t:Topic {name: $name})<-[:ORIGINATES_FROM]-(e:Event)
		RETURN t, e
	`
	params := map[string]interface{}{"name": name}

	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "t", "e")
	rels := r.BaseRepository.BuildRelationships(records, []struct{ SourceKey, TargetKey, Type string }{
		{"e", "t", "ORIGINATES_FROM"},
	}...)

	return &responses.GraphResponseDto{Nodes: nodes, Relationships: rels}, nil
}

func (r *TopicRepository) GetTopicWithSubscribedComponents(ctx context.Context, name string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (t:Topic {name: $name})<-[:SUBSCRIBES_TO]-(c:Component)
		RETURN t, c
	`
	params := map[string]interface{}{"name": name}

	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "t", "c")
	rels := r.BaseRepository.BuildRelationships(records, []struct{ SourceKey, TargetKey, Type string }{
		{"c", "t", "SUBSCRIBES_TO"},
	}...)

	return &responses.GraphResponseDto{Nodes: nodes, Relationships: rels}, nil
}

func (r *TopicRepository) GetTopicWithProducedComponents(ctx context.Context, name string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (t:Topic {name: $name})<-[:PRODUCES_TO]-(c:Component)
		RETURN t, c
	`
	params := map[string]interface{}{"name": name}

	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "t", "c")
	rels := r.BaseRepository.BuildRelationships(records, []struct{ SourceKey, TargetKey, Type string }{
		{"c", "t", "PRODUCES_TO"},
	}...)

	return &responses.GraphResponseDto{Nodes: nodes, Relationships: rels}, nil
}

func (r *TopicRepository) GetTopicWithEventsAndSubscribedComponents(ctx context.Context, name string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (t:Topic {name: $name})
		OPTIONAL MATCH (t)<-[:ORIGINATES_FROM]-(e:Event)
		OPTIONAL MATCH (t)<-[:SUBSCRIBES_TO]-(c:Component)
		RETURN t, e, c
	`
	params := map[string]interface{}{"name": name}

	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "t", "e", "c")
	rels := r.BaseRepository.BuildRelationships(records, []struct{ SourceKey, TargetKey, Type string }{
		{"e", "t", "ORIGINATES_FROM"},
		{"c", "t", "SUBSCRIBES_TO"},
	}...)

	return &responses.GraphResponseDto{Nodes: nodes, Relationships: rels}, nil
}

func (r *TopicRepository) GetTopicWithEventsAndProducedComponents(ctx context.Context, name string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (t:Topic {name: $name})
		OPTIONAL MATCH (t)<-[:ORIGINATES_FROM]-(e:Event)
		OPTIONAL MATCH (t)<-[:PRODUCES_TO]-(c:Component)
		RETURN t, e, c
	`
	params := map[string]interface{}{"name": name}

	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "t", "e", "c")
	rels := r.BaseRepository.BuildRelationships(records, []struct{ SourceKey, TargetKey, Type string }{
		{"e", "t", "ORIGINATES_FROM"},
		{"c", "t", "PRODUCES_TO"},
	}...)

	return &responses.GraphResponseDto{Nodes: nodes, Relationships: rels}, nil
}

func (r *TopicRepository) GetTopicWithSubscribedAndProducedComponents(ctx context.Context, name string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (t:Topic {name: $name})
		OPTIONAL MATCH (t)<-[:SUBSCRIBES_TO]-(c1:Component)
		OPTIONAL MATCH (t)<-[:PRODUCES_TO]-(c2:Component)
		RETURN t, c1, c2
	`
	params := map[string]interface{}{"name": name}

	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "t", "c1", "c2")
	rels := r.BaseRepository.BuildRelationships(records, []struct{ SourceKey, TargetKey, Type string }{
		{"c1", "t", "SUBSCRIBES_TO"},
		{"c2", "t", "PRODUCES_TO"},
	}...)

	return &responses.GraphResponseDto{Nodes: nodes, Relationships: rels}, nil
}

func (r *TopicRepository) GetTopicWithEventsAndComponents(ctx context.Context, name string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (t:Topic {name: $name})
		OPTIONAL MATCH (t)<-[:ORIGINATES_FROM]-(e:Event)
		OPTIONAL MATCH (t)<-[:SUBSCRIBES_TO]-(c1:Component)
		OPTIONAL MATCH (t)<-[:PRODUCES_TO]-(c2:Component)
		RETURN t, e, c1, c2
	`
	params := map[string]interface{}{"name": name}

	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "t", "e", "c1", "c2")
	rels := r.BaseRepository.BuildRelationships(records, []struct{ SourceKey, TargetKey, Type string }{
		{"e", "t", "ORIGINATES_FROM"},
		{"c1", "t", "SUBSCRIBES_TO"},
		{"c2", "t", "PRODUCES_TO"},
	}...)

	return &responses.GraphResponseDto{Nodes: nodes, Relationships: rels}, nil
}
