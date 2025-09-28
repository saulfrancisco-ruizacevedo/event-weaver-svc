package queries

const GetComponentsQueryName = "GetComponentsQuery"

type GetComponentsQuery struct{}

func NewGetComponentsQuery() *GetComponentsQuery {
	return &GetComponentsQuery{}
}
