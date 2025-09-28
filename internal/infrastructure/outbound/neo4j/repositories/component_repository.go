package repositories

import (
	"context"
	"fmt"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/responses"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/component"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/outbound/neo4j/mappers"
)

type ComponentRepository struct {
	BaseRepository *BaseRepository
}

func NewComponentRepository(baseRepository *BaseRepository) component.IComponentRepository {
	return &ComponentRepository{
		BaseRepository: baseRepository,
	}
}

var _ component.IComponentRepository = &ComponentRepository{}

func (r *ComponentRepository) GetAllComponentNames(ctx context.Context) ([]*component.Component, error) {
	query := `
		MATCH (c:Component)
		RETURN c.name AS name
	`

	records, err := r.BaseRepository.ExecQuery(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve component names: %w", err)
	}

	return mappers.RecordNamesListToComponentsList(records), nil
}

func (r *ComponentRepository) SaveAll(ctx context.Context, components []*component.Component) error {
	if len(components) == 0 {
		return nil
	}

	query := `
		UNWIND $components AS c
		MERGE (co:Component {name: c.name})
	`

	params := map[string]interface{}{
		"components": make([]map[string]interface{}, len(components)),
	}

	for i, component := range components {
		params["components"].([]map[string]interface{})[i] = map[string]interface{}{
			"name": component.Name,
		}
	}

	_, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to save components: %w", err)
	}

	return nil
}

func (r *ComponentRepository) CreateTeamRelationships(ctx context.Context, components []*component.Component) error {
	if len(components) == 0 {
		return nil
	}

	query := `
		UNWIND $components AS c
		MATCH (t:Team {name: c.team})
		MATCH (co:Component {name: c.name})
		MERGE (t)-[:MANAGES]->(co)
	`

	params := map[string]interface{}{
		"components": make([]map[string]interface{}, len(components)),
	}

	for i, component := range components {
		params["components"].([]map[string]interface{})[i] = map[string]interface{}{
			"name": component.Name,
			"team": component.Team,
		}
	}

	_, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to create team-component relationships: %w", err)
	}

	return nil
}

func (r *ComponentRepository) GetAllComponents(ctx context.Context) (*responses.GraphResponseDto, error) {
	return r.BaseRepository.GetAllNodes(ctx, "Component")
}

func (r *ComponentRepository) GetComponent(ctx context.Context, name string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (c:Component {name: $name})
		RETURN c
	`
	params := map[string]interface{}{"name": name}
	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "c")
	return &responses.GraphResponseDto{Nodes: nodes, Relationships: []responses.GraphRelationshipDto{}}, nil
}

func (r *ComponentRepository) GetComponentWithProducedEvents(ctx context.Context, name string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (c:Component {name: $name})-[:PRODUCES]->(e:Event)
		RETURN c, e
	`
	params := map[string]interface{}{"name": name}
	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "c", "e")
	rels := r.BaseRepository.BuildRelationships(records, []struct{ SourceKey, TargetKey, Type string }{
		{"c", "e", "PRODUCES"},
	}...)
	return &responses.GraphResponseDto{Nodes: nodes, Relationships: rels}, nil
}

func (r *ComponentRepository) GetComponentWithConsumedEvents(ctx context.Context, name string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (c:Component {name: $name})-[:CONSUMES]->(e:Event)
		RETURN c, e
	`
	params := map[string]interface{}{"name": name}
	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "c", "e")
	rels := r.BaseRepository.BuildRelationships(records, []struct{ SourceKey, TargetKey, Type string }{
		{"c", "e", "CONSUMES"},
	}...)
	return &responses.GraphResponseDto{Nodes: nodes, Relationships: rels}, nil
}

func (r *ComponentRepository) GetComponentWithSubscribedTopics(ctx context.Context, name string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (c:Component {name: $name})-[:SUBSCRIBES_TO]->(t:Topic)
		RETURN c, t
	`
	params := map[string]interface{}{"name": name}
	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "c", "t")
	rels := r.BaseRepository.BuildRelationships(records, []struct{ SourceKey, TargetKey, Type string }{
		{"c", "t", "SUBSCRIBES_TO"},
	}...)
	return &responses.GraphResponseDto{Nodes: nodes, Relationships: rels}, nil
}

func (r *ComponentRepository) GetComponentWithProducedTopics(ctx context.Context, name string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (c:Component {name: $name})-[:PRODUCES_TO]->(t:Topic)
		RETURN c, t
	`
	params := map[string]interface{}{"name": name}
	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "c", "t")
	rels := r.BaseRepository.BuildRelationships(records, []struct{ SourceKey, TargetKey, Type string }{
		{"c", "t", "PRODUCES_TO"},
	}...)
	return &responses.GraphResponseDto{Nodes: nodes, Relationships: rels}, nil
}

func (r *ComponentRepository) GetComponentWithTeam(ctx context.Context, name string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (t:Team)-[:MANAGES]->(c:Component {name: $name})
		RETURN t, c
	`
	params := map[string]interface{}{"name": name}
	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "c", "t")
	rels := r.BaseRepository.BuildRelationships(records, []struct{ SourceKey, TargetKey, Type string }{
		{"t", "c", "MANAGES"},
	}...)
	return &responses.GraphResponseDto{Nodes: nodes, Relationships: rels}, nil
}

func (r *ComponentRepository) GetComponentWithEventsAndTopicsAndTeams(ctx context.Context, name string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (c:Component {name: $name})
		OPTIONAL MATCH (c)-[:PRODUCES]->(e:Event)
		OPTIONAL MATCH (c)-[:CONSUMES]->(e2:Event)
		OPTIONAL MATCH (c)-[:SUBSCRIBES_TO]->(t:Topic)
		OPTIONAL MATCH (c)-[:PRODUCES_TO]->(t2:Topic)
		OPTIONAL MATCH (t3:Team)-[:MANAGES]->(c)
		RETURN c, e, e2, t, t2, t3
	`
	params := map[string]interface{}{"name": name}
	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "c", "e", "e2", "t", "t2", "t3")
	rels := r.BaseRepository.BuildRelationships(records, []struct{ SourceKey, TargetKey, Type string }{
		{"c", "e", "PRODUCES"},
		{"c", "e2", "CONSUMES"},
		{"c", "t", "SUBSCRIBES_TO"},
		{"c", "t2", "PRODUCES_TO"},
		{"t3", "c", "MANAGES"},
	}...)
	return &responses.GraphResponseDto{Nodes: nodes, Relationships: rels}, nil
}

func (r *ComponentRepository) GetComponentWithEventsAndTopicRelations(ctx context.Context, componentName string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (c:Component {name: $componentName})
		OPTIONAL MATCH (c)-[:PRODUCES]->(e:Event)
		OPTIONAL MATCH (c)-[:CONSUMES]->(e2:Event)
		OPTIONAL MATCH (c)-[:SUBSCRIBES_TO]->(t:Topic)
		OPTIONAL MATCH (c)-[:PRODUCES_TO]->(t2:Topic)
		RETURN c, e, e2, t, t2
	`
	params := map[string]interface{}{"componentName": componentName}

	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "c", "e", "e2", "t", "t2")
	rels := r.BaseRepository.BuildRelationships(records,
		struct{ SourceKey, TargetKey, Type string }{"c", "e", "PRODUCES"},
		struct{ SourceKey, TargetKey, Type string }{"c", "e2", "CONSUMES"},
		struct{ SourceKey, TargetKey, Type string }{"c", "t", "SUBSCRIBES_TO"},
		struct{ SourceKey, TargetKey, Type string }{"c", "t2", "PRODUCES_TO"},
	)

	return &responses.GraphResponseDto{
		Nodes:         nodes,
		Relationships: rels,
	}, nil
}

func (r *ComponentRepository) GetComponentWithEventsAndSubscriptions(ctx context.Context, componentName string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (c:Component {name: $componentName})
		OPTIONAL MATCH (c)-[:PRODUCES]->(e:Event)
		OPTIONAL MATCH (c)-[:CONSUMES]->(e2:Event)
		OPTIONAL MATCH (c)-[:SUBSCRIBES_TO]->(t:Topic)
		RETURN c, e, e2, t
	`
	params := map[string]interface{}{"componentName": componentName}

	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "c", "e", "e2", "t")
	rels := r.BaseRepository.BuildRelationships(records,
		struct{ SourceKey, TargetKey, Type string }{"c", "e", "PRODUCES"},
		struct{ SourceKey, TargetKey, Type string }{"c", "e2", "CONSUMES"},
		struct{ SourceKey, TargetKey, Type string }{"c", "t", "SUBSCRIBES_TO"},
	)

	return &responses.GraphResponseDto{
		Nodes:         nodes,
		Relationships: rels,
	}, nil
}

func (r *ComponentRepository) GetComponentWithProducedAndConsumedEvents(ctx context.Context, componentName string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (c:Component {name: $componentName})
		OPTIONAL MATCH (c)-[:PRODUCES]->(e:Event)
		OPTIONAL MATCH (c)-[:CONSUMES]->(e2:Event)
		RETURN c, e, e2
	`
	params := map[string]interface{}{"componentName": componentName}

	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "c", "e", "e2")
	rels := r.BaseRepository.BuildRelationships(records,
		struct{ SourceKey, TargetKey, Type string }{"c", "e", "PRODUCES"},
		struct{ SourceKey, TargetKey, Type string }{"c", "e2", "CONSUMES"},
	)

	return &responses.GraphResponseDto{
		Nodes:         nodes,
		Relationships: rels,
	}, nil
}

func (r *ComponentRepository) GetComponentWithEventsAndTeam(ctx context.Context, componentName string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (c:Component {name: $componentName})
		OPTIONAL MATCH (c)-[:PRODUCES]->(e:Event)
		OPTIONAL MATCH (c)-[:CONSUMES]->(e2:Event)
		OPTIONAL MATCH (t:Team)-[:MANAGES]->(c)
		RETURN c, e, e2, t
	`
	params := map[string]interface{}{"componentName": componentName}

	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "c", "e", "e2", "t")
	rels := r.BaseRepository.BuildRelationships(records,
		struct{ SourceKey, TargetKey, Type string }{"c", "e", "PRODUCES"},
		struct{ SourceKey, TargetKey, Type string }{"c", "e2", "CONSUMES"},
		struct{ SourceKey, TargetKey, Type string }{"t", "c", "MANAGES"},
	)

	return &responses.GraphResponseDto{
		Nodes:         nodes,
		Relationships: rels,
	}, nil
}
