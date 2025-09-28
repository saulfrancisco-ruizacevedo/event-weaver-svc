package queries

import "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/requests"

const GetEventQueryName = "GetEventQuery"

type GetEventQuery struct {
	Payload requests.GetEventRequestDto
}

func NewGetEventQuery(payload requests.GetEventRequestDto) *GetEventQuery {
	return &GetEventQuery{
		Payload: payload,
	}
}
