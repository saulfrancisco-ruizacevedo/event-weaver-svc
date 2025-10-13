package mappers

import (
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/component"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/outbound/neo4j/models"
)

func ModelComponentListToDomainList(modelComponents []*models.Component) []*component.Component {
	domainComponents := make([]*component.Component, len(modelComponents))
	for i, mc := range modelComponents {
		domainComponents[i] = &component.Component{
			Name: mc.Name,
		}
	}
	return domainComponents
}

func DomainComponentListToModelList(domainComponents []*component.Component) []*models.Component {
	modelComponents := make([]*models.Component, len(domainComponents))
	for i, dc := range domainComponents {
		modelComponents[i] = &models.Component{
			Name: dc.Name,
		}
	}
	return modelComponents
}

func DomainComponentToModel(dc *component.Component) *models.Component {
	if dc == nil {
		return nil
	}
	return &models.Component{
		Name: dc.Name,
	}
}

func DomainComponentToTeamModel(dc *component.Component) *models.Team {
	if dc == nil || dc.Team == "" {
		return nil
	}
	return &models.Team{
		Name: dc.Team,
	}
}
