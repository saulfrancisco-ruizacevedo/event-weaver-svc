package repositories

import (
	"context"
	"fmt"

	domainentity "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/domainentity"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/outbound/neo4j/mappers"
)

type DomainRepository struct {
	BaseRepository *BaseRepository
}

func NewDomainRepository(baseRepository *BaseRepository) domainentity.IDomainRepository {
	return &DomainRepository{
		BaseRepository: baseRepository,
	}
}

var _ domainentity.IDomainRepository = &DomainRepository{}

func (r *DomainRepository) GetAllDomainNames(ctx context.Context) ([]*domainentity.Domain, error) {
	query := `
		MATCH (d:Domain)
		RETURN d.name AS name
	`

	records, err := r.BaseRepository.ExecQuery(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve domain names: %w", err)
	}

	return mappers.RecordNamesListToDomainList(records), nil
}

func (r *DomainRepository) SaveAll(ctx context.Context, domains []*domainentity.Domain) error {
	if len(domains) == 0 {
		return nil
	}

	query := `
		UNWIND $domains AS d
		MERGE (dom:Domain {name: d.name})
		SET dom.description = d.description
	`

	params := map[string]interface{}{
		"domains": make([]map[string]interface{}, len(domains)),
	}

	for i, domain := range domains {
		params["domains"].([]map[string]interface{})[i] = map[string]interface{}{
			"name":        domain.Name,
			"description": domain.Description,
		}
	}

	_, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to save domains: %w", err)
	}

	return nil
}
