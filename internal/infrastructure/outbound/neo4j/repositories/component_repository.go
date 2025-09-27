package repositories

import (
	"context"
	"fmt"

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
