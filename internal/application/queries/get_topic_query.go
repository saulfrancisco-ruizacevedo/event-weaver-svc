package queries

import "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/requests"

const GetTopicQueryName = "GetTopicQuery"

type GetTopicQuery struct {
	Payload requests.GetTopicRequestDto
}

func NewGetTopicQuery(payload requests.GetTopicRequestDto) *GetTopicQuery {
	return &GetTopicQuery{
		Payload: payload,
	}
}
