package repositories

import (
	"context"
	"fmt"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/event"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/outbound/neo4j/mappers"
	localModels "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/outbound/neo4j/models"
	"github.com/saulfrancisco-ruizacevedo/go-neopersist"
	"github.com/saulfrancisco-ruizacevedo/gocypher"
)

type EventRepository struct {
	manager   *neopersist.PersistenceManager
	eventRepo *neopersist.Repository[localModels.Event]
}

func NewEventRepository(manager *neopersist.PersistenceManager) (event.IEventRepository, error) {
	eventRepo, err := neopersist.RepositoryFor[localModels.Event](manager)
	if err != nil {
		return nil, fmt.Errorf("failed to create event repository: %w", err)
	}
	return &EventRepository{
		manager:   manager,
		eventRepo: eventRepo,
	}, nil
}

var _ event.IEventRepository = &EventRepository{}

func (r *EventRepository) FindAllNamesByNamesIn(names []*string, ctx context.Context) ([]*event.EventSpecification, error) {
	if len(names) == 0 {
		return []*event.EventSpecification{}, nil
	}

	existing := make(map[string]struct{})
	cleanNames := make([]string, 0, len(names))
	for _, n := range names {
		if n == nil || *n == "" {
			continue
		}
		if _, ok := existing[*n]; !ok {
			existing[*n] = struct{}{}
			cleanNames = append(cleanNames, *n)
		}
	}
	params := map[string]interface{}{
		"names": cleanNames,
	}

	if len(cleanNames) == 0 {
		return []*event.EventSpecification{}, nil
	}

	qb := gocypher.NewQueryBuilder().
		Match(gocypher.N("e", "Event")).
		Where("e.name IN $names").
		Return("e.name AS name").WithParams(params)

	modelEvents, err := r.eventRepo.Find(ctx, qb)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve event names: %w", err)
	}

	return mappers.ModelEventListToSpecificationList(modelEvents), nil
}

func (r *EventRepository) SaveAll(ctx context.Context, events []*event.EventSpecification) error {
	if len(events) == 0 {
		return nil
	}

	modelEvents, err := mappers.DomainEventSpecListToModelList(events)
	if err != nil {
		return fmt.Errorf("failed to map event specifications to model: %w", err)
	}
	return r.eventRepo.SaveAll(ctx, modelEvents)
}

// func (r *EventRepository) CreateProducerRelationships(ctx context.Context, events []*event.EventSpecification) error {
// 	if len(events) == 0 {
// 		return nil
// 	}

// 	query := `
// 		UNWIND $pairs AS rel
// 		MATCH (co:Component {name: rel.component})
// 		MATCH (ev:Event {name: rel.event})
// 		MERGE (co)-[:PRODUCES]->(ev)
// 	`

// 	var pairs []map[string]interface{}
// 	for _, ev := range events {
// 		for _, producer := range ev.Producers {
// 			pairs = append(pairs, map[string]interface{}{
// 				"component": producer.Name,
// 				"event":     ev.Name,
// 			})
// 		}
// 	}

// 	params := map[string]interface{}{"pairs": pairs}

// 	_, err := r.BaseRepository.ExecQuery(ctx, query, params)
// 	if err != nil {
// 		return fmt.Errorf("failed to create producer relationships: %w", err)
// 	}

// 	return nil
// }

func (r *EventRepository) CreateProducerRelationships(ctx context.Context, events []*event.EventSpecification) error {
	if len(events) == 0 {
		return nil
	}

	for _, ev := range events {
		for _, producer := range ev.Producers {
			componentModel := mappers.ProducerToComponentModel(&producer)
			eventModel := mappers.EventSpecToEventModel(ev)

			err := r.manager.CreateRelation(ctx, componentModel, eventModel, "PRODUCES", nil)
			if err != nil {
				return fmt.Errorf("failed to create PRODUCES relationship between component '%s' and event '%s': %w", producer.Name, ev.Name, err)
			}
		}
	}

	return nil
}

// func (r *EventRepository) CreateConsumerRelationships(ctx context.Context, events []*event.EventSpecification) error {
// 	if len(events) == 0 {
// 		return nil
// 	}

// 	query := `
// 		UNWIND $pairs AS rel
// 		MATCH (co:Component {name: rel.component})
// 		MATCH (ev:Event {name: rel.event})
// 		MERGE (co)-[:CONSUMES]->(ev)
// 	`

// 	var pairs []map[string]interface{}
// 	for _, ev := range events {
// 		for _, consumer := range ev.Consumers {
// 			pairs = append(pairs, map[string]interface{}{
// 				"component": consumer.Name,
// 				"event":     ev.Name,
// 			})
// 		}
// 	}

// 	params := map[string]interface{}{"pairs": pairs}

// 	_, err := r.BaseRepository.ExecQuery(ctx, query, params)
// 	if err != nil {
// 		return fmt.Errorf("failed to create consumer relationships: %w", err)
// 	}

// 	return nil
// }

func (r *EventRepository) CreateConsumerRelationships(ctx context.Context, events []*event.EventSpecification) error {
	if len(events) == 0 {
		return nil
	}

	for _, ev := range events {
		for _, consumer := range ev.Consumers {
			componentModel := mappers.ConsumerToComponentModel(&consumer)
			eventModel := mappers.EventSpecToEventModel(ev)

			err := r.manager.CreateRelation(ctx, componentModel, eventModel, "CONSUMES", nil)
			if err != nil {
				return fmt.Errorf("failed to create CONSUMES relationship between component '%s' and event '%s': %w", consumer.Name, ev.Name, err)
			}
		}
	}

	return nil
}

// func (r *EventRepository) CreateTopicRelationships(ctx context.Context, events []*event.EventSpecification) error {
// 	if len(events) == 0 {
// 		return nil
// 	}

// 	query := `
// 		UNWIND $pairs AS rel
// 		MATCH (ev:Event {name: rel.event})
// 		MERGE (to:Topic {name: rel.topic})
// 		MERGE (ev)-[:ORIGINATES_FROM]->(to)
// 	`

// 	var pairs []map[string]interface{}
// 	for _, ev := range events {
// 		if ev.Topic == "" {
// 			continue
// 		}
// 		pairs = append(pairs, map[string]interface{}{
// 			"event": ev.Name,
// 			"topic": ev.Topic,
// 		})
// 	}

// 	if len(pairs) == 0 {
// 		return nil
// 	}

// 	params := map[string]interface{}{"pairs": pairs}

// 	_, err := r.BaseRepository.ExecQuery(ctx, query, params)
// 	if err != nil {
// 		return fmt.Errorf("failed to create topic relationships: %w", err)
// 	}

// 	return nil
// }

// func (r *EventRepository) CreateRelatedEventRelationships(ctx context.Context, events []*event.EventSpecification) error {
// 	if len(events) == 0 {
// 		return nil
// 	}

// 	query := `
// 		UNWIND $pairs AS rel
// 		MATCH (src:Event {name: rel.source})
// 		MATCH (dst:Event {name: rel.target})
// 		MERGE (src)-[:RELATED_TO]->(dst)
// 	`

// 	var pairs []map[string]interface{}
// 	for _, ev := range events {
// 		for _, related := range ev.RelatedEvents {
// 			pairs = append(pairs, map[string]interface{}{
// 				"source": ev.Name,
// 				"target": related.Name,
// 			})
// 		}
// 	}

// 	if len(pairs) == 0 {
// 		return nil
// 	}

// 	params := map[string]interface{}{"pairs": pairs}

// 	_, err := r.BaseRepository.ExecQuery(ctx, query, params)
// 	if err != nil {
// 		return fmt.Errorf("failed to create related event relationships: %w", err)
// 	}

// 	return nil
// }

// func (r *EventRepository) CreateDomainRelationships(ctx context.Context, events []*event.EventSpecification) error {
// 	if len(events) == 0 {
// 		return nil
// 	}

// 	query := `
// 		UNWIND $pairs AS rel
// 		MATCH (ev:Event {name: rel.event})
// 		MERGE (dom:Domain {name: rel.domain})
// 		MERGE (ev)-[:BELONGS_TO]->(dom)
// 	`

// 	var pairs []map[string]interface{}
// 	for _, ev := range events {
// 		if ev.Domain == "" {
// 			continue
// 		}
// 		pairs = append(pairs, map[string]interface{}{
// 			"event":  ev.Name,
// 			"domain": ev.Domain,
// 		})
// 	}

// 	if len(pairs) == 0 {
// 		return nil
// 	}

// 	params := map[string]interface{}{"pairs": pairs}

// 	_, err := r.BaseRepository.ExecQuery(ctx, query, params)
// 	if err != nil {
// 		return fmt.Errorf("failed to create domain relationships: %w", err)
// 	}

// 	return nil
// }

func (r *EventRepository) CreateTopicRelationships(ctx context.Context, events []*event.EventSpecification) error {
	if len(events) == 0 {
		return nil
	}

	for _, ev := range events {
		if ev.Topic == "" {
			continue
		}

		eventModel := mappers.EventSpecToEventModel(ev)
		topicModel := mappers.TopicNameToTopicModel(ev.Topic)

		err := r.manager.CreateRelation(ctx, eventModel, topicModel, "ORIGINATES_FROM", nil)
		if err != nil {
			return fmt.Errorf("failed to create ORIGINATES_FROM relationship for event '%s': %w", ev.Name, err)
		}
	}

	return nil
}

func (r *EventRepository) CreateRelatedEventRelationships(ctx context.Context, events []*event.EventSpecification) error {
	if len(events) == 0 {
		return nil
	}

	for _, ev := range events {
		for _, related := range ev.RelatedEvents {
			sourceEventModel := mappers.EventSpecToEventModel(ev)
			targetEventModel := mappers.RelatedEventToEventModel(&related)

			err := r.manager.CreateRelation(ctx, sourceEventModel, targetEventModel, "RELATED_TO", nil)
			if err != nil {
				return fmt.Errorf("failed to create RELATED_TO relationship between event '%s' and '%s': %w", ev.Name, related.Name, err)
			}
		}
	}

	return nil
}

func (r *EventRepository) CreateDomainRelationships(ctx context.Context, events []*event.EventSpecification) error {
	if len(events) == 0 {
		return nil
	}

	for _, ev := range events {
		if ev.Domain == "" {
			continue
		}

		eventModel := mappers.EventSpecToEventModel(ev)
		domainModel := mappers.DomainNameToDomainModel(ev.Domain)

		err := r.manager.CreateRelation(ctx, eventModel, domainModel, "BELONGS_TO", nil)
		if err != nil {
			return fmt.Errorf("failed to create BELONGS_TO relationship for event '%s': %w", ev.Name, err)
		}
	}

	return nil
}
