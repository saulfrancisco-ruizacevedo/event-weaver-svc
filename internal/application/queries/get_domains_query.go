package queries

const GetDomainsQueryName = "GetDomainsQuery"

type GetDomainsQuery struct{}

func NewGetDomainsQuery() *GetDomainsQuery {
	return &GetDomainsQuery{}
}
