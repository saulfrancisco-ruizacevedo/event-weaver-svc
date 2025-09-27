package responses

import "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/validations"

type EventSpecValidationResponseDto struct {
	DomainResult    *validations.ValidationResult `json:"domainResult,omitempty"`
	TeamResult      *validations.ValidationResult `json:"teamResult,omitempty"`
	ComponentResult *validations.ValidationResult `json:"componentResult,omitempty"`
	TopicResult     *validations.ValidationResult `json:"topicResult,omitempty"`
	EventResult     *validations.ValidationResult `json:"eventResult,omitempty"`
}
