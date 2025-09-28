package queries

const GetTopicsQueryName = "GetTopicsQuery"

type GetTopicsQuery struct{}

func NewGetTopicsQuery() *GetTopicsQuery {
	return &GetTopicsQuery{}
}
