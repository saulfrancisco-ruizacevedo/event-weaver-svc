package mappers

import (
	"github.com/neo4j/neo4j-go-driver/v6/neo4j"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/component"
)

func RecordNamesListToComponentsList(records []*neo4j.Record) []*component.Component {
	result := make([]*component.Component, 0, len(records))

	for _, record := range records {
		name, _ := record.Get("name")
		result = append(result, &component.Component{
			Name: name.(string),
		})
	}
	return result
}
