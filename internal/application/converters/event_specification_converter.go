package converters

import (
	"fmt"
	"time"

	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/requests"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/component"
	domainentity "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/domainentity"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/event"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/team"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/topic"
)

func EventDtoListToEventSpecificationList(dtos []requests.EventDto) ([]*event.EventSpecification, error) {
	events := make([]*event.EventSpecification, 0, len(dtos))

	for _, dto := range dtos {
		eventSpec := event.NewEventSpecification(dto.Name, dto.Domain)

		eventSpec.Description = dto.Description
		eventSpec.Topic = dto.Topic

		eventSpec.Schema = event.Schema{
			Properties: make(map[string]event.SchemaProperty),
		}
		for propName, propDto := range dto.Schema.Properties {
			eventSpec.Schema.Properties[propName] = event.SchemaProperty{
				Type:        propDto.Type,
				Description: propDto.Description,
			}
		}

		for _, p := range dto.Producers {
			eventSpec.AddProducer(p.Name)
		}
		for _, c := range dto.Consumers {
			eventSpec.AddConsumer(c.Name)
		}
		for _, r := range dto.RelatedEvents {
			eventSpec.AddRelatedEvent(r.Name)
		}
		for _, tag := range dto.Tags {
			eventSpec.AddTag(tag)
		}

		var effectiveFrom *time.Time
		if dto.Lifecycle.EffectiveFrom != nil {
			parsedTime, err := time.Parse("2006-01-02", *dto.Lifecycle.EffectiveFrom)
			if err != nil {
				return nil, fmt.Errorf("failed to parse effectiveFrom date for event '%s': %w", dto.Name, err)
			}
			effectiveFrom = &parsedTime
		}

		var effectiveTo *time.Time
		if dto.Lifecycle.EffectiveTo != nil {
			parsedTime, err := time.Parse("2006-01-02", *dto.Lifecycle.EffectiveTo)
			if err != nil {
				return nil, fmt.Errorf("failed to parse effectiveTo date for event '%s': %w", dto.Name, err)
			}
			effectiveTo = &parsedTime
		}

		eventSpec.Lifecycle = event.Lifecycle{
			Deprecated:    dto.Lifecycle.Deprecated,
			EffectiveFrom: effectiveFrom,
			EffectiveTo:   effectiveTo,
		}

		eventSpec.Examples = dto.Examples

		events = append(events, eventSpec)
	}

	return events, nil
}

func TeamDtoListToTeamList(dtos []requests.TeamDto) []*team.Team {
	teams := make([]*team.Team, 0, len(dtos))

	for _, dto := range dtos {
		team := team.NewTeam(dto.Name, dto.Lead, dto.Email)
		teams = append(teams, team)
	}

	return teams
}

func ComponentDtoListToComponentList(dtos []requests.ComponentDto) []*component.Component {
	components := make([]*component.Component, 0, len(dtos))

	for _, dto := range dtos {
		component := component.NewComponent(dto.Name, dto.Team)
		components = append(components, component)
	}

	return components
}

func DomainDtoListToDomainList(dtos []requests.DomainDto) []*domainentity.Domain {
	domains := make([]*domainentity.Domain, 0, len(dtos))

	for _, dto := range dtos {
		dom := domainentity.NewDomain(dto.Name, dto.Description)
		domains = append(domains, dom)
	}

	return domains
}

func TopicDtoListToTopicList(dtos []requests.TopicDto) []*topic.Topic {
	topics := make([]*topic.Topic, 0, len(dtos))

	for _, dto := range dtos {
		dom := topic.NewTopic(dto.Name, dto.Type)
		topics = append(topics, dom)
	}

	return topics
}
