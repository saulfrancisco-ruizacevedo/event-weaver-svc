package repositories

import (
	"context"

	domainentity "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/domainentity"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/outbound/neo4j/mappers"
	localModels "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/outbound/neo4j/models"
	"github.com/saulfrancisco-ruizacevedo/go-neopersist"
	"github.com/saulfrancisco-ruizacevedo/gocypher"
)

type DomainRepository struct {
	manager    *neopersist.PersistenceManager
	domainRepo *neopersist.Repository[localModels.Domain]
}

func NewDomainRepository(manager *neopersist.PersistenceManager) (domainentity.IDomainRepository, error) {
	domainRepo, err := neopersist.RepositoryFor[localModels.Domain](manager)
	if err != nil {
		panic(err)
	}

	return &DomainRepository{
		manager:    manager,
		domainRepo: domainRepo,
	}, nil
}

var _ domainentity.IDomainRepository = &DomainRepository{}

func (r *DomainRepository) GetAllDomainNames(ctx context.Context) ([]*domainentity.Domain, error) {
	qb := gocypher.NewQueryBuilder().
		Match(gocypher.N("d", "Domain")).
		Return("d.name AS name")

	modelDomains, err := r.domainRepo.Find(ctx, qb)

	if err != nil {
		return nil, err
	}

	return mappers.ModelListToDomainList(modelDomains), nil
}

func (r *DomainRepository) SaveAll(ctx context.Context, domains []*domainentity.Domain) error {
	if len(domains) == 0 {
		return nil
	}

	modelDomains := mappers.DomainListToModel(domains)

	return r.domainRepo.SaveAll(ctx, modelDomains)
}
