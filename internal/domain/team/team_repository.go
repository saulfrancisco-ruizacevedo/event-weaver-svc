package team

import "context"

type ITeamRepository interface {
	GetAllTeamNames(ctx context.Context) ([]*Team, error)
	SaveAll(ctx context.Context, team []*Team) error
}
