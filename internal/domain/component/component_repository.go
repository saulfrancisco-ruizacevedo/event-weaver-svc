package component

import (
	"context"
)

type IComponentRepository interface {
	GetAllComponentNames(ctx context.Context) ([]*Component, error)
	SaveAll(ctx context.Context, components []*Component) error
	CreateTeamRelationships(ctx context.Context, components []*Component) error
}
