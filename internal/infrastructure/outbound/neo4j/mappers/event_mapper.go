package mappers

import (
	"github.com/neo4j/neo4j-go-driver/v6/neo4j"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/event"
)

func EventNamesListToTopicList(records []*neo4j.Record) []*event.EventSpecification {
	result := make([]*event.EventSpecification, 0, len(records))

	for _, record := range records {
		name, _ := record.Get("name")
		result = append(result, &event.EventSpecification{
			Name: name.(string),
		})
	}

	return result
}
