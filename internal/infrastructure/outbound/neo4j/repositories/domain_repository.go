package repositories

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/responses"
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
		log.Error().Err(err).Msgf("Error")
		return fmt.Errorf("failed to save domains: %w", err)
	}
	return nil
}

func (r *DomainRepository) GetAllDomains(ctx context.Context) (*responses.GraphResponseDto, error) {
	return r.BaseRepository.GetAllNodes(ctx, "Domain")
}

func (r *DomainRepository) GetDomainWithEvents(ctx context.Context, domainName string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (d:Domain {name: $domainName})<-[:BELONGS_TO]-(e:Event)
		RETURN d, e
	`
	params := map[string]interface{}{"domainName": domainName}

	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "d", "e")
	rels := r.BaseRepository.BuildRelationships(records, []struct{ SourceKey, TargetKey, Type string }{
		{"e", "d", "BELONGS_TO"},
	}...)

	return &responses.GraphResponseDto{
		Nodes:         nodes,
		Relationships: rels,
	}, nil
}

func (r *DomainRepository) GetDomainWithEventsAndComponents(ctx context.Context, domainName string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (d:Domain {name: $domainName})<-[:BELONGS_TO]-(e:Event)<-[:PRODUCES]-(c:Component)
		RETURN d, e, c
	`
	params := map[string]interface{}{"domainName": domainName}

	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "d", "e", "c")
	rels := r.BaseRepository.BuildRelationships(records, []struct{ SourceKey, TargetKey, Type string }{
		{"e", "d", "BELONGS_TO"},
		{"c", "e", "PRODUCES"},
	}...)

	return &responses.GraphResponseDto{
		Nodes:         nodes,
		Relationships: rels,
	}, nil
}

func (r *DomainRepository) GetDomain(ctx context.Context, domainName string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (d:Domain {name: $domainName})
		RETURN d
	`
	params := map[string]interface{}{"domainName": domainName}

	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}

	nodes := r.BaseRepository.BuildNodes(records, "d")

	return &responses.GraphResponseDto{
		Nodes:         nodes,
		Relationships: []responses.GraphRelationshipDto{},
	}, nil
}
