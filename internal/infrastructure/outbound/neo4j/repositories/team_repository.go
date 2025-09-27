package repositories

import (
	"context"
	"fmt"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/team"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/outbound/neo4j/mappers"
)

type TeamRepository struct {
	BaseRepository *BaseRepository
}

func NewTeamRepository(baseRepository *BaseRepository) team.ITeamRepository {
	return &TeamRepository{
		BaseRepository: baseRepository,
	}
}

var _ team.ITeamRepository = &TeamRepository{}

func (r *TeamRepository) GetAllTeamNames(ctx context.Context) ([]*team.Team, error) {
	query := `
		MATCH (t:Team)
		RETURN t.name AS name
	`

	records, err := r.BaseRepository.ExecQuery(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve team names: %w", err)
	}

	return mappers.RecordNamesListToTeamsList(records), nil
}

func (r *TeamRepository) SaveAll(ctx context.Context, teams []*team.Team) error {
	if len(teams) == 0 {
		return nil
	}

	query := `
		UNWIND $teams AS t
		MERGE (te:Team {name: t.name})
		SET te.lead = t.lead,
			te.email = t.email
	`

	params := map[string]interface{}{
		"teams": make([]map[string]interface{}, len(teams)),
	}

	for i, team := range teams {
		params["teams"].([]map[string]interface{})[i] = map[string]interface{}{
			"name":  team.Name,
			"lead":  team.Lead,
			"email": team.Email,
		}
	}

	_, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to save teams: %w", err)
	}

	return nil
}
