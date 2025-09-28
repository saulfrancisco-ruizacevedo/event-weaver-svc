package domainentity

import (
	"context"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/responses"
)

type IDomainRepository interface {
	GetAllDomainNames(ctx context.Context) ([]*Domain, error)
	SaveAll(ctx context.Context, domains []*Domain) error
	GetAllDomains(ctx context.Context) (*responses.GraphResponseDto, error)
	GetDomain(ctx context.Context, domainName string) (*responses.GraphResponseDto, error)
	GetDomainWithEvents(ctx context.Context, domainName string) (*responses.GraphResponseDto, error)
	GetDomainWithEventsAndComponents(ctx context.Context, domainName string) (*responses.GraphResponseDto, error)
}
