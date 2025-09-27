package requests

type EventSpecificationRequestDto struct {
	SpecVersion string         `yaml:"specVersion"`
	Domains     []DomainDto    `yaml:"domains,omitempty"`
	Teams       []TeamDto      `yaml:"teams,omitempty"`
	Components  []ComponentDto `yaml:"components,omitempty"`
	Topics      []TopicDto     `yaml:"topics,omitempty"`
	Events      []EventDto     `yaml:"events,omitempty"`
	Migrations  MigrationsDto  `yaml:"migrations,omitempty"`
}

type DomainDto struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
}

type TeamDto struct {
	Name  string `yaml:"name"`
	Lead  string `yaml:"lead,omitempty"`
	Email string `yaml:"email,omitempty"`
}

type ComponentDto struct {
	Name string `yaml:"name"`
	Team string `yaml:"team,omitempty"`
}

type TopicDto struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type EventDto struct {
	Name          string                 `yaml:"name"`
	Domain        string                 `yaml:"domain,omitempty"`
	Description   string                 `yaml:"description,omitempty"`
	Topic         string                 `yaml:"topic,omitempty"`
	Schema        Schema                 `yaml:"schema,omitempty"`
	Producers     []Producer             `yaml:"producers,omitempty"`
	Consumers     []Consumer             `yaml:"consumers,omitempty"`
	RelatedEvents []RelatedEvent         `yaml:"relatedEvents,omitempty"`
	Tags          []string               `yaml:"tags,omitempty"`
	Examples      map[string]interface{} `yaml:"examples,omitempty"`
	Lifecycle     Lifecycle              `yaml:"lifecycle,omitempty"`
}

type SchemaProperty struct {
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
}

type Schema struct {
	Properties map[string]SchemaProperty `yaml:"properties,omitempty"`
}

type Producer struct {
	Name string `yaml:"name"`
}

type Consumer struct {
	Name string `yaml:"name"`
}

type RelatedEvent struct {
	Name string `yaml:"name"`
}

type Lifecycle struct {
	Deprecated    bool    `yaml:"deprecated,omitempty"`
	EffectiveFrom *string `yaml:"effectiveFrom,omitempty"`
	EffectiveTo   *string `yaml:"effectiveTo,omitempty"`
}

type MigrationsDto struct {
	Deletions DeletionsDto `yaml:"deletions,omitempty"`
	Updates   []UpdateDto  `yaml:"updates,omitempty"`
}

type DeletionsDto struct {
	Components []DeletionComponentDto `yaml:"components,omitempty"`
	Teams      []DeletionTeamDto      `yaml:"teams,omitempty"`
}

type DeletionComponentDto struct {
	Name string `yaml:"name"`
}

type DeletionTeamDto struct {
	Name string `yaml:"name"`
}

type UpdateDto struct {
	Type    string `yaml:"type"`
	OldName string `yaml:"oldName"`
	NewName string `yaml:"newName"`
}
