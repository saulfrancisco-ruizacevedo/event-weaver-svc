package spec

import (
	"fmt"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/validations"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/component"
	domainentity "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/domainentity"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/event"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/team"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/topic"
)

type SpecValidator struct{}

func NewSpecValidator() *SpecValidator {
	return &SpecValidator{}
}

func (s *SpecValidator) ValidateDomains(domains []*domainentity.Domain) *validations.ValidationResult {
	result := validations.ValidationResult{Passed: true}

	for _, dom := range domains {
		if res := dom.Validate(); !res.Passed {
			result.Passed = false
			result.Errors = append(result.Errors, res.Errors...)
		}
	}

	return &result
}

func (s *SpecValidator) ValidateTeams(teams []*team.Team) *validations.ValidationResult {
	result := validations.ValidationResult{Passed: true}

	for _, t := range teams {
		if res := t.Validate(); !res.Passed {
			result.Passed = false
			result.Errors = append(result.Errors, res.Errors...)
		}
	}

	return &result
}

func (s *SpecValidator) ValidateComponents(components []*component.Component, teams []*team.Team) *validations.ValidationResult {
	result := validations.ValidationResult{Passed: true}
	teamMap := make(map[string]struct{})

	for _, t := range teams {
		teamMap[t.Name] = struct{}{}
	}

	for _, c := range components {

		if res := c.Validate(); !res.Passed {
			result.Passed = false
			result.Errors = append(result.Errors, res.Errors...)
		}

		if _, exists := teamMap[c.Team]; !exists {
			result.Passed = false
			result.Errors = append(result.Errors, validations.ValidationError{
				EntityName:   c.Name,
				EntityType:   "component",
				ErrorCode:    validations.NonExistentReference,
				ErrorMessage: fmt.Sprintf("references a non-existent team: '%s'", c.Team),
			})
		}
	}

	return &result
}

func (s *SpecValidator) ValidateEvents(events []*event.EventSpecification,
	components []*component.Component,
	domains []*domainentity.Domain,
	topics []*topic.Topic,
	relatedEvents []*event.EventSpecification) *validations.ValidationResult {

	result := validations.ValidationResult{Passed: true}

	componentMap := make(map[string]struct{})
	for _, c := range components {
		componentMap[c.Name] = struct{}{}
	}

	domainMap := make(map[string]struct{})
	for _, d := range domains {
		domainMap[d.Name] = struct{}{}
	}

	topicsMap := make(map[string]struct{})
	for _, t := range topics {
		topicsMap[t.Name] = struct{}{}
	}

	relatedEvetsMap := make(map[string]struct{})
	for _, e := range relatedEvents {
		relatedEvetsMap[e.Name] = struct{}{}
	}

	for _, e := range events {
		if res := e.Validate(); !res.Passed {
			result.Passed = false
			result.Errors = append(result.Errors, res.Errors...)
		}

		for _, producer := range e.Producers {
			if _, exists := componentMap[producer.Name]; !exists {
				result.Passed = false
				result.Errors = append(result.Errors, validations.ValidationError{
					EntityName:   e.Name,
					EntityType:   "event",
					ErrorCode:    validations.NonExistentReference,
					ErrorMessage: fmt.Sprintf("has a producer(Component) '%s' that doesn't exist", producer.Name),
				})
			}
		}

		for _, consumer := range e.Consumers {
			if _, exists := componentMap[consumer.Name]; !exists {
				result.Passed = false
				result.Errors = append(result.Errors, validations.ValidationError{
					EntityName:   e.Name,
					EntityType:   "event",
					ErrorCode:    validations.NonExistentReference,
					ErrorMessage: fmt.Sprintf("has a consumer(Component) '%s' that doesn't exist", consumer.Name),
				})
			}
		}

		if _, exists := domainMap[e.Domain]; !exists {
			result.Passed = false
			result.Errors = append(result.Errors, validations.ValidationError{
				EntityName:   e.Name,
				EntityType:   "event",
				ErrorCode:    validations.NonExistentReference,
				ErrorMessage: fmt.Sprintf("has a domain '%s' that doesn't exist", e.Domain),
			})
		}

		if _, exists := topicsMap[e.Topic]; !exists {
			result.Passed = false
			result.Errors = append(result.Errors, validations.ValidationError{
				EntityName:   e.Name,
				EntityType:   "event",
				ErrorCode:    validations.NonExistentReference,
				ErrorMessage: fmt.Sprintf("has a topic '%s' that doesn't exist", e.Topic),
			})
		}

		for _, relatedEvent := range e.RelatedEvents {
			if _, exists := relatedEvetsMap[relatedEvent.Name]; !exists {
				result.Passed = false
				result.Errors = append(result.Errors, validations.ValidationError{
					EntityName:   e.Name,
					EntityType:   "event",
					ErrorCode:    validations.NonExistentReference,
					ErrorMessage: fmt.Sprintf("has a related event '%s' that doesn't exist", relatedEvent.Name),
				})
			}
		}
	}
	return &result
}

func (s *SpecValidator) ValidateTopics(topics []*topic.Topic) *validations.ValidationResult {
	result := validations.ValidationResult{Passed: true}

	for _, topic := range topics {
		if res := topic.Validate(); !res.Passed {
			result.Passed = false
			result.Errors = append(result.Errors, res.Errors...)
		}
	}

	return &result
}
