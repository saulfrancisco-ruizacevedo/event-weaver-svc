package repositories

import (
	"context"
	"fmt"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/responses"
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

func (r *TeamRepository) GetAllTeams(ctx context.Context) (*responses.GraphResponseDto, error) {
	return r.BaseRepository.GetAllNodes(ctx, "Team")
}

func (r *TeamRepository) GetTeam(ctx context.Context, teamName string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (t:Team {name: $teamName})
		RETURN t
	`
	params := map[string]interface{}{"teamName": teamName}

	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "t")

	return &responses.GraphResponseDto{
		Nodes:         nodes,
		Relationships: []responses.GraphRelationshipDto{},
	}, nil
}

func (r *TeamRepository) GetTeamWithComponents(ctx context.Context, teamName string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (t:Team {name: $teamName})-[:MANAGES]->(c:Component)
		RETURN t, c
	`
	params := map[string]interface{}{"teamName": teamName}

	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "t", "c")
	rels := r.BaseRepository.BuildRelationships(records, []struct{ SourceKey, TargetKey, Type string }{
		{"t", "c", "MANAGES"},
	}...)

	return &responses.GraphResponseDto{
		Nodes:         nodes,
		Relationships: rels,
	}, nil
}

func (r *TeamRepository) GetTeamWithComponentsAndEvents(ctx context.Context, teamName string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (t:Team {name: $teamName})-[:MANAGES]->(c:Component)-[:PRODUCES]->(e:Event)
		RETURN t, c, e
	`
	params := map[string]interface{}{"teamName": teamName}

	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "t", "c", "e")
	rels := r.BaseRepository.BuildRelationships(records, []struct{ SourceKey, TargetKey, Type string }{
		{"t", "c", "MANAGES"},
		{"c", "e", "PRODUCES"},
	}...)

	return &responses.GraphResponseDto{
		Nodes:         nodes,
		Relationships: rels,
	}, nil
}

func (r *TeamRepository) GetTeamWithComponentsAndEventsAndDomains(ctx context.Context, teamName string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (t:Team {name: $teamName})-[:MANAGES]->(c:Component)-[:PRODUCES]->(e:Event)-[:BELONGS_TO]->(d:Domain)
		RETURN t, c, e, d
	`
	params := map[string]interface{}{"teamName": teamName}

	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "t", "c", "e", "d")
	rels := r.BaseRepository.BuildRelationships(records, []struct{ SourceKey, TargetKey, Type string }{
		{"t", "c", "MANAGES"},
		{"c", "e", "PRODUCES"},
		{"e", "d", "BELONGS_TO"},
	}...)

	return &responses.GraphResponseDto{
		Nodes:         nodes,
		Relationships: rels,
	}, nil
}
