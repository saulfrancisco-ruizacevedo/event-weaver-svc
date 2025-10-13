package repositories

import (
	"context"
	"fmt"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/team"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/outbound/neo4j/mappers"
	localModels "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/outbound/neo4j/models"
	"github.com/saulfrancisco-ruizacevedo/go-neopersist"
	"github.com/saulfrancisco-ruizacevedo/gocypher"
)

type TeamRepository struct {
	manager  *neopersist.PersistenceManager
	teamRepo *neopersist.Repository[localModels.Team]
}

func NewTeamRepository(manager *neopersist.PersistenceManager) (team.ITeamRepository, error) {
	teamRepo, err := neopersist.RepositoryFor[localModels.Team](manager)
	if err != nil {
		return nil, fmt.Errorf("failed to create team repository: %w", err)
	}
	return &TeamRepository{
		manager:  manager,
		teamRepo: teamRepo,
	}, nil
}

var _ team.ITeamRepository = &TeamRepository{}

func (r *TeamRepository) GetAllTeamNames(ctx context.Context) ([]*team.Team, error) {
	qb := gocypher.NewQueryBuilder().
		Match(gocypher.N("t", "Team")).
		Return("t.name AS name")

	modelTeams, err := r.teamRepo.Find(ctx, qb)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve team names: %w", err)
	}

	return mappers.ModelTeamListToDomainList(modelTeams), nil
}

func (r *TeamRepository) SaveAll(ctx context.Context, teams []*team.Team) error {
	if len(teams) == 0 {
		return nil
	}

	modelTeams := mappers.DomainTeamListToModelList(teams)

	return r.teamRepo.SaveAll(ctx, modelTeams)
}
