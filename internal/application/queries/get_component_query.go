package queries

import "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/requests"

const GetComponentQueryName = "GetComponentQuery"

type GetComponentQuery struct {
	Payload requests.GetComponentRequestDto
}

func NewGetComponentQuery(payload requests.GetComponentRequestDto) *GetComponentQuery {
	return &GetComponentQuery{
		Payload: payload,
	}
}
