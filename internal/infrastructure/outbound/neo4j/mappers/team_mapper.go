package mappers

import (
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/team"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/outbound/neo4j/models"
)

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
