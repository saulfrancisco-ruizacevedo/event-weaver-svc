package queries

import "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/requests"

const GetTeamQueryName = "GetTeamQuery"

type GetTeamQuery struct {
	Payload requests.GetTeamRequestDto
}

func NewGetTeamQuery(payload requests.GetTeamRequestDto) *GetTeamQuery {
	return &GetTeamQuery{
		Payload: payload,
	}
}
