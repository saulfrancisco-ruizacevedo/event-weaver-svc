package mappers

import (
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/domainentity"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/outbound/neo4j/models"
)

func ModelListToDomainList(modelDomains []*models.Domain) []*domainentity.Domain {

	domains := make([]*domainentity.Domain, len(modelDomains))
	for i, modelDomain := range modelDomains {
		domains[i] = &domainentity.Domain{
			Name: modelDomain.Name,
		}
	}

	return domains
}

func DomainListToModel(domains []*domainentity.Domain) []*models.Domain {
	modelDomains := make([]*models.Domain, len(domains))
	for i, d := range domains {
		modelDomains[i] = &models.Domain{
			Name:        d.Name,
			Description: d.Description,
		}
	}

	return modelDomains
}
