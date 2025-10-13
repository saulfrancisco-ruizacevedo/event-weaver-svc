package domainentity

import (
	"context"
)

type IDomainRepository interface {
	GetAllDomainNames(ctx context.Context) ([]*Domain, error)
	SaveAll(ctx context.Context, domains []*Domain) error
}
