package event

import (
	"strings"
	"time"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/validations"
)

// ------------------------
// Aggregate Root: EventSpecification
// ------------------------

type EventSpecification struct {
	Name          string
	Domain        string
	Description   string
	Topic         string
	Schema        Schema
	Producers     []Producer
	Consumers     []Consumer
	RelatedEvents []RelatedEvent
	Tags          []string
	Examples      map[string]interface{}
	Lifecycle     Lifecycle
}

type Schema struct {
	Properties map[string]SchemaProperty
}

type SchemaProperty struct {
	Type        string
	Description string
}

type Owner struct {
	Name  string
	Email string
}

type Producer struct {
	Name string
}

type Consumer struct {
	Name string
}

type Origin struct {
	Type string
	Name string
}

type RelatedEvent struct {
	Name string
}

type Lifecycle struct {
	Version       string
	Deprecated    bool
	EffectiveFrom *time.Time
	EffectiveTo   *time.Time
}

func NewEventSpecification(name, domain string) *EventSpecification {
	return &EventSpecification{
		Name:          name,
		Domain:        domain,
		Producers:     []Producer{},
		Consumers:     []Consumer{},
		RelatedEvents: []RelatedEvent{},
		Tags:          []string{},
		Examples:      make(map[string]interface{}),
	}
}

func (e *EventSpecification) GetName() string {
	return e.Name
}

func (e *EventSpecification) AddProducer(name string) {
	e.Producers = append(e.Producers, Producer{Name: name})
}

func (e *EventSpecification) AddConsumer(name string) {
	e.Consumers = append(e.Consumers, Consumer{Name: name})
}

func (e *EventSpecification) AddRelatedEvent(name string) {
	e.RelatedEvents = append(e.RelatedEvents, RelatedEvent{Name: name})
}

func (e *EventSpecification) AddTag(tag string) {
	e.Tags = append(e.Tags, tag)
}

func (e *EventSpecification) ValidateName() *validations.ValidationError {
	if strings.TrimSpace(e.Name) == "" {
		return &validations.ValidationError{
			EntityName:   e.Name,
			EntityType:   "event",
			ErrorCode:    validations.EmptyField,
			ErrorMessage: "event name cannot be empty",
		}
	}

	return nil
}

func (e *EventSpecification) Validate() validations.ValidationResult {
	result := validations.ValidationResult{Passed: true}

	if err := e.ValidateName(); err != nil {
		result.Errors = append(result.Errors, *err)
	}

	return result
}

func (e *EventSpecification) GetRelatedEventsNames() []*string {
	names := make([]*string, 0, len(e.RelatedEvents))

	for _, rel := range e.RelatedEvents {
		if rel.Name != "" {
			name := rel.Name
			names = append(names, &name)
		}
	}

	return names
}

func (l Lifecycle) EffectiveFromISO() interface{} {
	if l.EffectiveFrom != nil {
		return l.EffectiveFrom.Format(time.RFC3339)
	}

	return nil
}

func (l Lifecycle) EffectiveToISO() interface{} {
	if l.EffectiveTo != nil {
		return l.EffectiveTo.Format(time.RFC3339)
	}

	return nil
}
