package mappers

import (
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/topic"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/outbound/neo4j/models"
)

func ModelTopicListToDomainList(modelTopics []*models.Topic) []*topic.Topic {
	domainTopics := make([]*topic.Topic, len(modelTopics))
	for i, mc := range modelTopics {
		domainTopics[i] = &topic.Topic{
			Name: mc.Name,
			Type: mc.Type,
		}
	}
	return domainTopics
}

func DomainTopicListToModelList(domainComponents []*topic.Topic) []*models.Topic {
	modelTopics := make([]*models.Topic, len(domainComponents))
	for i, dc := range domainComponents {
		modelTopics[i] = &models.Topic{
			Name: dc.Name,
			Type: dc.Type,
		}
	}
	return modelTopics
}
