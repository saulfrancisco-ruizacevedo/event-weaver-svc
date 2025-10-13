package mappers

import (
	"encoding/json"
	"fmt"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/event"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/outbound/neo4j/models"
)

func ModelEventListToSpecificationList(modelEvents []*models.Event) []*event.EventSpecification {
	domainEvent := make([]*event.EventSpecification, len(modelEvents))

	for i, me := range modelEvents {
		domainEvent[i] = &event.EventSpecification{
			Name: me.Name,
		}
	}
	return domainEvent
}

func DomainEventSpecListToModelList(domainEvents []*event.EventSpecification) ([]*models.Event, error) {
	modelEvents := make([]*models.Event, 0, len(domainEvents))

	for _, domainEvent := range domainEvents {
		schemaJSON, err := json.Marshal(domainEvent.Schema.Properties)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal schema for event %s: %w", domainEvent.Name, err)
		}

		exampleJSON, err := json.Marshal(domainEvent.Examples)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal example for event %s: %w", domainEvent.Name, err)
		}

		model := &models.Event{
			Name:          domainEvent.Name,
			Description:   domainEvent.Description,
			Schema:        string(schemaJSON),
			Examples:      string(exampleJSON),
			Deprecated:    domainEvent.Lifecycle.Deprecated,
			EffectiveFrom: domainEvent.Lifecycle.EffectiveFromISO(),
			EffectiveTo:   domainEvent.Lifecycle.EffectiveToISO(),
		}
		modelEvents = append(modelEvents, model)
	}

	return modelEvents, nil
}

func ProducerToComponentModel(producer *event.Producer) *models.Component {
	return &models.Component{
		Name: producer.Name,
	}
}

func EventSpecToEventModel(spec *event.EventSpecification) *models.Event {
	return &models.Event{
		Name: spec.Name,
	}
}

func ConsumerToComponentModel(consumer *event.Consumer) *models.Component {
	return &models.Component{
		Name: consumer.Name,
	}
}

func TopicNameToTopicModel(name string) *models.Topic {
	return &models.Topic{
		Name: name,
	}
}

func RelatedEventToEventModel(relatedEvent *event.RelatedEvent) *models.Event {
	return &models.Event{
		Name: relatedEvent.Name,
	}
}

func DomainNameToDomainModel(name string) *models.Domain {
	return &models.Domain{
		Name: name,
	}
}
