package queries

import "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/requests"

const GetDomainQueryName = "GetDomainQuery"

type GetDomainQuery struct {
	Payload requests.GetDomainRequestDto
}

func NewGetDomainQuery(payload requests.GetDomainRequestDto) *GetDomainQuery {
	return &GetDomainQuery{
		Payload: payload,
	}
}
