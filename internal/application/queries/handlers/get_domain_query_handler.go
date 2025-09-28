package handlers

import (
	"context"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/responses"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/queries"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/domainentity"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/mediator"
)

type GetDomainQueryHandler struct {
	domainRepository domainentity.IDomainRepository
}

var _ mediator.MediatorHandler[*queries.GetDomainQuery, *responses.GraphResponseDto] = &GetDomainQueryHandler{}

func NewGetDomainQueryHandler(domainRepository domainentity.IDomainRepository) *GetDomainQueryHandler {
	return &GetDomainQueryHandler{
		domainRepository: domainRepository,
	}
}

func (h *GetDomainQueryHandler) Handle(ctx context.Context, query *queries.GetDomainQuery) (*responses.GraphResponseDto, error) {
	dto := query.Payload

	if dto.RelatedToEvent && dto.ComponentProducesEvent {
		return h.domainRepository.GetDomainWithEventsAndComponents(ctx, dto.Name)

	} else if dto.RelatedToEvent {
		return h.domainRepository.GetDomainWithEvents(ctx, dto.Name)

	} else {
		return h.domainRepository.GetDomain(ctx, dto.Name)
	}
}
