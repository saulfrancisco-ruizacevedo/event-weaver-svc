package queries

const GetEventsQueryName = "GetEventsQuery"

type GetEventsQuery struct{}

func NewGetEventsQuery() *GetEventsQuery {
	return &GetEventsQuery{}
}
