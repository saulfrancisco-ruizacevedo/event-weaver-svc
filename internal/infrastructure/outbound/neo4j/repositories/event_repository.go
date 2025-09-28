package repositories

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/responses"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/event"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/outbound/neo4j/mappers"
)

type EventRepository struct {
	BaseRepository *BaseRepository
}

func NewEventRepository(BaseRepository *BaseRepository) event.IEventRepository {
	return &EventRepository{
		BaseRepository: BaseRepository,
	}
}

var _ event.IEventRepository = &EventRepository{}

func (r *EventRepository) FindAllNamesByNamesIn(names []*string, ctx context.Context) ([]*event.EventSpecification, error) {
	if len(names) == 0 {
		return []*event.EventSpecification{}, nil
	}

	existing := make(map[string]struct{})
	cleanNames := make([]*string, 0, len(names))
	for _, n := range names {
		if n == nil || *n == "" {
			continue
		}
		if _, ok := existing[*n]; ok {
			continue
		}
		existing[*n] = struct{}{}
		cleanNames = append(cleanNames, n)
	}

	if len(cleanNames) == 0 {
		return []*event.EventSpecification{}, nil
	}

	namesParam := make([]interface{}, len(cleanNames))
	for i, n := range cleanNames {
		namesParam[i] = *n
	}

	query := `
        MATCH (e:Event)
        WHERE e.name IN $names
        RETURN e.name AS name
    `
	params := map[string]interface{}{
		"names": namesParam,
	}

	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve event names: %w", err)
	}

	return mappers.EventNamesListToTopicList(records), nil
}

func (r *EventRepository) SaveAll(ctx context.Context, events []*event.EventSpecification) error {
	if len(events) == 0 {
		return nil
	}

	query := `
		UNWIND $events AS e
		MERGE (ev:Event {name: e.name})
		SET ev.description   = e.description,
		    ev.schema        = e.schema,
		    ev.examples      = e.examples,
		    ev.deprecated    = e.deprecated,
		    ev.effectiveFrom = e.effectiveFrom,
		    ev.effectiveTo   = e.effectiveTo
	`

	params := map[string]interface{}{
		"events": make([]map[string]interface{}, len(events)),
	}

	for i, event := range events {

		schemaJSON, err := json.Marshal(event.Schema.Properties)
		if err != nil {
			return fmt.Errorf("failed to marshal schema for event %s: %w", event.Name, err)
		}

		exampleJSON, err := json.Marshal(event.Examples)
		if err != nil {
			return fmt.Errorf("failed to marshal example for event %s: %w", event.Name, err)
		}

		params["events"].([]map[string]interface{})[i] = map[string]interface{}{
			"name":          event.Name,
			"description":   event.Description,
			"deprecated":    event.Lifecycle.Deprecated,
			"effectiveFrom": event.Lifecycle.EffectiveFromISO(),
			"effectiveTo":   event.Lifecycle.EffectiveToISO(),
			"schema":        string(schemaJSON),
			"examples":      string(exampleJSON),
		}
	}

	_, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to save events: %w", err)
	}
	return nil
}

func (r *EventRepository) CreateProducerRelationships(ctx context.Context, events []*event.EventSpecification) error {
	if len(events) == 0 {
		return nil
	}

	query := `
		UNWIND $pairs AS rel
		MATCH (co:Component {name: rel.component})
		MATCH (ev:Event {name: rel.event})
		MERGE (co)-[:PRODUCES]->(ev)
	`

	var pairs []map[string]interface{}
	for _, ev := range events {
		for _, producer := range ev.Producers {
			pairs = append(pairs, map[string]interface{}{
				"component": producer.Name,
				"event":     ev.Name,
			})
		}
	}

	params := map[string]interface{}{"pairs": pairs}

	_, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to create producer relationships: %w", err)
	}

	return nil
}

func (r *EventRepository) CreateConsumerRelationships(ctx context.Context, events []*event.EventSpecification) error {
	if len(events) == 0 {
		return nil
	}

	query := `
		UNWIND $pairs AS rel
		MATCH (co:Component {name: rel.component})
		MATCH (ev:Event {name: rel.event})
		MERGE (co)-[:CONSUMES]->(ev)
	`

	var pairs []map[string]interface{}
	for _, ev := range events {
		for _, consumer := range ev.Consumers {
			pairs = append(pairs, map[string]interface{}{
				"component": consumer.Name,
				"event":     ev.Name,
			})
		}
	}

	params := map[string]interface{}{"pairs": pairs}

	_, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to create consumer relationships: %w", err)
	}

	return nil
}

func (r *EventRepository) CreateTopicRelationships(ctx context.Context, events []*event.EventSpecification) error {
	if len(events) == 0 {
		return nil
	}

	query := `
		UNWIND $pairs AS rel
		MATCH (ev:Event {name: rel.event})
		MERGE (to:Topic {name: rel.topic})
		MERGE (ev)-[:ORIGINATES_FROM]->(to)
	`

	var pairs []map[string]interface{}
	for _, ev := range events {
		if ev.Topic == "" {
			continue
		}
		pairs = append(pairs, map[string]interface{}{
			"event": ev.Name,
			"topic": ev.Topic,
		})
	}

	if len(pairs) == 0 {
		return nil
	}

	params := map[string]interface{}{"pairs": pairs}

	_, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to create topic relationships: %w", err)
	}

	return nil
}

func (r *EventRepository) CreateRelatedEventRelationships(ctx context.Context, events []*event.EventSpecification) error {
	if len(events) == 0 {
		return nil
	}

	query := `
		UNWIND $pairs AS rel
		MATCH (src:Event {name: rel.source})
		MATCH (dst:Event {name: rel.target})
		MERGE (src)-[:RELATED_TO]->(dst)
	`

	var pairs []map[string]interface{}
	for _, ev := range events {
		for _, related := range ev.RelatedEvents {
			pairs = append(pairs, map[string]interface{}{
				"source": ev.Name,
				"target": related.Name,
			})
		}
	}

	if len(pairs) == 0 {
		return nil
	}

	params := map[string]interface{}{"pairs": pairs}

	_, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to create related event relationships: %w", err)
	}

	return nil
}

func (r *EventRepository) CreateDomainRelationships(ctx context.Context, events []*event.EventSpecification) error {
	if len(events) == 0 {
		return nil
	}

	query := `
		UNWIND $pairs AS rel
		MATCH (ev:Event {name: rel.event})
		MERGE (dom:Domain {name: rel.domain})
		MERGE (ev)-[:BELONGS_TO]->(dom)
	`

	var pairs []map[string]interface{}
	for _, ev := range events {
		if ev.Domain == "" {
			continue
		}
		pairs = append(pairs, map[string]interface{}{
			"event":  ev.Name,
			"domain": ev.Domain,
		})
	}

	if len(pairs) == 0 {
		return nil
	}

	params := map[string]interface{}{"pairs": pairs}

	_, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to create domain relationships: %w", err)
	}

	return nil
}

func (r *EventRepository) GetAllEvents(ctx context.Context) (*responses.GraphResponseDto, error) {
	return r.BaseRepository.GetAllNodes(ctx, "Event")
}

func (r *EventRepository) GetEvent(ctx context.Context, name string) (*responses.GraphResponseDto, error) {
	query := `MATCH (e:Event {name: $name}) RETURN e`
	params := map[string]interface{}{"name": name}
	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}
	nodes := r.BaseRepository.BuildNodes(records, "e")
	return &responses.GraphResponseDto{Nodes: nodes, Relationships: []responses.GraphRelationshipDto{}}, nil
}

func (r *EventRepository) GetEventWithEvents(ctx context.Context, name string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (e:Event {name: $name})-[:RELATED_TO]->(related:Event)
		RETURN e, related
	`
	params := map[string]interface{}{"name": name}
	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}
	nodes := r.BaseRepository.BuildNodes(records, "e", "related")
	rels := r.BaseRepository.BuildRelationships(records, struct{ SourceKey, TargetKey, Type string }{"e", "related", "RELATED_TO"})
	return &responses.GraphResponseDto{Nodes: nodes, Relationships: rels}, nil
}

func (r *EventRepository) GetEventWithEventsAndProduced(ctx context.Context, name string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (e:Event {name: $name})-[:RELATED_TO]->(related:Event)
		OPTIONAL MATCH (c:Component)-[:PRODUCES]->(e)
		RETURN e, related, c
	`
	params := map[string]interface{}{"name": name}
	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}
	nodes := r.BaseRepository.BuildNodes(records, "e", "related", "c")
	rels := r.BaseRepository.BuildRelationships(records,
		struct{ SourceKey, TargetKey, Type string }{"e", "related", "RELATED_TO"},
		struct{ SourceKey, TargetKey, Type string }{"c", "e", "PRODUCES"},
	)
	return &responses.GraphResponseDto{Nodes: nodes, Relationships: rels}, nil
}

func (r *EventRepository) GetEventWithEventsProducedAndConsumed(ctx context.Context, name string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (e:Event {name: $name})-[:RELATED_TO]->(related:Event)
		OPTIONAL MATCH (c1:Component)-[:PRODUCES]->(e)
		OPTIONAL MATCH (c2:Component)-[:CONSUMES]->(e)
		RETURN e, related, c1, c2
	`
	params := map[string]interface{}{"name": name}
	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}
	nodes := r.BaseRepository.BuildNodes(records, "e", "related", "c1", "c2")
	rels := r.BaseRepository.BuildRelationships(records,
		struct{ SourceKey, TargetKey, Type string }{"e", "related", "RELATED_TO"},
		struct{ SourceKey, TargetKey, Type string }{"c1", "e", "PRODUCES"},
		struct{ SourceKey, TargetKey, Type string }{"c2", "e", "CONSUMES"},
	)
	return &responses.GraphResponseDto{Nodes: nodes, Relationships: rels}, nil
}

func (r *EventRepository) GetEventWithEventsProducedConsumedAndTopic(ctx context.Context, name string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (e:Event {name: $name})-[:RELATED_TO]->(related:Event)
		OPTIONAL MATCH (c1:Component)-[:PRODUCES]->(e)
		OPTIONAL MATCH (c2:Component)-[:CONSUMES]->(e)
		OPTIONAL MATCH (e)-[:ORIGINATES_FROM]->(t:Topic)
		RETURN e, related, c1, c2, t
	`
	params := map[string]interface{}{"name": name}
	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}
	nodes := r.BaseRepository.BuildNodes(records, "e", "related", "c1", "c2", "t")
	rels := r.BaseRepository.BuildRelationships(records,
		struct{ SourceKey, TargetKey, Type string }{"e", "related", "RELATED_TO"},
		struct{ SourceKey, TargetKey, Type string }{"c1", "e", "PRODUCES"},
		struct{ SourceKey, TargetKey, Type string }{"c2", "e", "CONSUMES"},
		struct{ SourceKey, TargetKey, Type string }{"e", "t", "ORIGINATES_FROM"},
	)
	return &responses.GraphResponseDto{Nodes: nodes, Relationships: rels}, nil
}

func (r *EventRepository) GetEventWithEventsAndComponentsAndTopicsAndDomain(ctx context.Context, name string) (*responses.GraphResponseDto, error) {
	query := `
		MATCH (e:Event {name: $name})-[:RELATED_TO]->(related:Event)
		OPTIONAL MATCH (c:Component)-[:PRODUCES]->(e)
		OPTIONAL MATCH (c2:Component)-[:CONSUMES]->(e)
		OPTIONAL MATCH (e)-[:ORIGINATES_FROM]->(t:Topic)
		OPTIONAL MATCH (e)-[:BELONGS_TO]->(d:Domain)
		RETURN e, related, c, c2, t, d
	`
	params := map[string]interface{}{"name": name}
	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}
	nodes := r.BaseRepository.BuildNodes(records, "e", "related", "c", "c2", "t", "d")
	rels := r.BaseRepository.BuildRelationships(records,
		struct{ SourceKey, TargetKey, Type string }{"e", "related", "RELATED_TO"},
		struct{ SourceKey, TargetKey, Type string }{"c", "e", "PRODUCES"},
		struct{ SourceKey, TargetKey, Type string }{"c2", "e", "CONSUMES"},
		struct{ SourceKey, TargetKey, Type string }{"e", "t", "ORIGINATES_FROM"},
		struct{ SourceKey, TargetKey, Type string }{"e", "d", "BELONGS_TO"},
	)
	return &responses.GraphResponseDto{Nodes: nodes, Relationships: rels}, nil
}

func (r *EventRepository) GetEventProducedByComponents(ctx context.Context, name string) (*responses.GraphResponseDto, error) {
	query := `MATCH (c:Component)-[:PRODUCES]->(e:Event {name: $name}) RETURN e, c`
	params := map[string]interface{}{"name": name}
	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}
	nodes := r.BaseRepository.BuildNodes(records, "e", "c")
	rels := r.BaseRepository.BuildRelationships(records, struct{ SourceKey, TargetKey, Type string }{"c", "e", "PRODUCES"})
	return &responses.GraphResponseDto{Nodes: nodes, Relationships: rels}, nil
}

func (r *EventRepository) GetEventConsumedByComponents(ctx context.Context, name string) (*responses.GraphResponseDto, error) {
	query := `MATCH (c:Component)-[:CONSUMES]->(e:Event {name: $name}) RETURN e, c`
	params := map[string]interface{}{"name": name}
	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}
	nodes := r.BaseRepository.BuildNodes(records, "e", "c")
	rels := r.BaseRepository.BuildRelationships(records, struct{ SourceKey, TargetKey, Type string }{"c", "e", "CONSUMES"})
	return &responses.GraphResponseDto{Nodes: nodes, Relationships: rels}, nil
}

func (r *EventRepository) GetEventOriginatedFromTopic(ctx context.Context, name string) (*responses.GraphResponseDto, error) {
	query := `MATCH (e:Event {name: $name})-[:ORIGINATES_FROM]->(t:Topic) RETURN e, t`
	params := map[string]interface{}{"name": name}
	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}
	nodes := r.BaseRepository.BuildNodes(records, "e", "t")
	rels := r.BaseRepository.BuildRelationships(records, struct{ SourceKey, TargetKey, Type string }{"e", "t", "ORIGINATES_FROM"})
	return &responses.GraphResponseDto{Nodes: nodes, Relationships: rels}, nil
}

func (r *EventRepository) GetEventBelongsToDomain(ctx context.Context, name string) (*responses.GraphResponseDto, error) {
	query := `MATCH (e:Event {name: $name})-[:BELONGS_TO]->(d:Domain) RETURN e, d`
	params := map[string]interface{}{"name": name}
	records, err := r.BaseRepository.ExecQuery(ctx, query, params)
	if err != nil {
		return nil, err
	}
	nodes := r.BaseRepository.BuildNodes(records, "e", "d")
	rels := r.BaseRepository.BuildRelationships(records, struct{ SourceKey, TargetKey, Type string }{"e", "d", "BELONGS_TO"})
	return &responses.GraphResponseDto{Nodes: nodes, Relationships: rels}, nil
}
