package topic

import "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/validations"

type Topic struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

func NewTopic(name, topicType string) *Topic {
	return &Topic{
		Name: name,
		Type: topicType,
	}
}

func (t *Topic) GetName() string {
	return t.Name
}

func (t *Topic) Validate() validations.ValidationResult {
	result := validations.ValidationResult{Passed: true}

	if t.Name == "" {
		result.Passed = false
		result.Errors = append(result.Errors, validations.ValidationError{
			EntityName:   t.Name,
			EntityType:   "team",
			ErrorCode:    validations.EmptyField,
			ErrorMessage: "topic name cannot be empty",
		})
	}

	if t.Type == "" {
		result.Passed = false
		result.Errors = append(result.Errors, validations.ValidationError{
			EntityName:   t.Name,
			EntityType:   "team",
			ErrorCode:    validations.EmptyField,
			ErrorMessage: "topic type cannot be empty",
		})
	}

	return result
}
