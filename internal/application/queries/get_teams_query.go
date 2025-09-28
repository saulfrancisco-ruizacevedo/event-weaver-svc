package queries

const GetTeamsQueryName = "GetTeamsQuery"

type GetTeamsQuery struct{}

func NewGetTeamsQuery() *GetTeamsQuery {
	return &GetTeamsQuery{}
}
