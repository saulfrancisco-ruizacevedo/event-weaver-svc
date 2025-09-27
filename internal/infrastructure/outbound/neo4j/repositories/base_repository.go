package repositories

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v6/neo4j"
)

type BaseRepository struct {
	Driver neo4j.Driver
	DbName string
}

func (r *BaseRepository) ExecQuery(ctx context.Context, query string, params map[string]interface{}) ([]*neo4j.Record, error) {
	result, err := neo4j.ExecuteQuery(
		ctx,
		r.Driver,
		query,
		params,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(r.DbName),
	)

	if err != nil {
		return nil, err
	}

	return result.Records, nil
}
