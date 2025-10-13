package mappers

import (
	"github.com/neo4j/neo4j-go-driver/v6/neo4j"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/team"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/outbound/neo4j/models"
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

func ModelTeamListToDomainList(modelTeams []*models.Team) []*team.Team {
	domainTeam := make([]*team.Team, len(modelTeams))
	for i, mc := range modelTeams {
		domainTeam[i] = &team.Team{
			Name:  mc.Name,
			Lead:  mc.Lead,
			Email: mc.Email,
		}
	}
	return domainTeam
}

func DomainTeamListToModelList(domainComponents []*team.Team) []*models.Team {
	modelTeam := make([]*models.Team, len(domainComponents))
	for i, dc := range domainComponents {
		modelTeam[i] = &models.Team{
			Name:  dc.Name,
			Lead:  dc.Lead,
			Email: dc.Email,
		}
	}
	return modelTeam
}
