package mappers

import (
	"github.com/neo4j/neo4j-go-driver/v6/neo4j"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/team"
)

func TeamToParams(team *team.Team) map[string]any {
	return map[string]any{
		"name":  team.Name,
		"lead":  team.Lead,
		"email": team.Email,
	}
}

func RecordToTeam(record *neo4j.Record) *team.Team {
	return &team.Team{
		Name:  record.Values[0].(string),
		Lead:  record.Values[1].(string),
		Email: record.Values[2].(string),
	}
}

func RecordNamesListToTeamsList(records []*neo4j.Record) []*team.Team {
	result := make([]*team.Team, 0, len(records))

	for _, record := range records {
		name, _ := record.Get("name")
		result = append(result, &team.Team{
			Name: name.(string),
		})
	}

	return result
}
