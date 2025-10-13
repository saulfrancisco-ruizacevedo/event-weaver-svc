package repositories

import (
	"context"
	"fmt"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/component"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/outbound/neo4j/mappers"
	localModels "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/outbound/neo4j/models"
	"github.com/saulfrancisco-ruizacevedo/go-neopersist"
	"github.com/saulfrancisco-ruizacevedo/gocypher"
)

type ComponentRepository struct {
	manager       *neopersist.PersistenceManager
	componentRepo *neopersist.Repository[localModels.Component]
}

func NewComponentRepository(manager *neopersist.PersistenceManager) (component.IComponentRepository, error) {
	componentRepo, err := neopersist.RepositoryFor[localModels.Component](manager)
	if err != nil {
		return nil, fmt.Errorf("failed to create component repository: %w", err)
	}
	return &ComponentRepository{
		manager:       manager,
		componentRepo: componentRepo,
	}, nil
}

var _ component.IComponentRepository = &ComponentRepository{}

func (r *ComponentRepository) GetAllComponentNames(ctx context.Context) ([]*component.Component, error) {
	qb := gocypher.NewQueryBuilder().
		Match(gocypher.N("c", "Component")).
		Return("c.name AS name")

	modelComponents, err := r.componentRepo.Find(ctx, qb)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve component names: %w", err)
	}

	return mappers.ModelComponentListToDomainList(modelComponents), nil
}

func (r *ComponentRepository) SaveAll(ctx context.Context, components []*component.Component) error {
	if len(components) == 0 {
		return nil
	}

	modelComponents := mappers.DomainComponentListToModelList(components)

	return r.componentRepo.SaveAll(ctx, modelComponents)
}

func (r *ComponentRepository) CreateTeamRelationships(ctx context.Context, components []*component.Component) error {
	if len(components) == 0 {
		return nil
	}

	for _, c := range components {
		if c.Team == "" {
			continue
		}
		teamModel := mappers.DomainComponentToTeamModel(c)
		componentModel := mappers.DomainComponentToModel(c)

		err := r.manager.CreateRelation(ctx, teamModel, componentModel, "MANAGES", nil)
		if err != nil {
			return fmt.Errorf("failed to create MANAGES relationship for component '%s': %w", c.Name, err)
		}
	}

	return nil
}
