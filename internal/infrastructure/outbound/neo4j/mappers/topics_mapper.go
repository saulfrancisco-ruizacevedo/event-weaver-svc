package mappers

import (
	"github.com/neo4j/neo4j-go-driver/v6/neo4j"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/topic"
)

func TopicNamesListToTopicList(records []*neo4j.Record) []*topic.Topic {
	result := make([]*topic.Topic, 0, len(records))

	for _, record := range records {
		name, _ := record.Get("name")
		result = append(result, &topic.Topic{
			Name: name.(string),
		})
	}

	return result
}
