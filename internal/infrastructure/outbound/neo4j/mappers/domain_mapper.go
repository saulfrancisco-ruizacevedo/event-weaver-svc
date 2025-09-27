package mappers

import (
	"github.com/neo4j/neo4j-go-driver/v6/neo4j"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/domainentity"
)

func RecordNamesListToDomainList(records []*neo4j.Record) []*domainentity.Domain {
	result := make([]*domainentity.Domain, 0, len(records))

	for _, record := range records {
		name, _ := record.Get("name")
		result = append(result, &domainentity.Domain{
			Name: name.(string),
		})
	}

	return result
}
