package handlers

import (
	"context"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/responses"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/queries"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/domainentity"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/mediator"
)

type GetDomainsQueryHandler struct {
	domainRepository domainentity.IDomainRepository
}

var _ mediator.MediatorHandler[*queries.GetDomainsQuery, *responses.GraphResponseDto] = &GetDomainsQueryHandler{}

func NewGetDomainsQueryHandler(domainRepository domainentity.IDomainRepository) *GetDomainsQueryHandler {
	return &GetDomainsQueryHandler{
		domainRepository: domainRepository,
	}
}

func (h *GetDomainsQueryHandler) Handle(ctx context.Context, query *queries.GetDomainsQuery) (*responses.GraphResponseDto, error) {
	if graph, err := h.domainRepository.GetAllDomains(ctx); err != nil {
		return nil, err

	} else {
		return graph, nil
	}
}
