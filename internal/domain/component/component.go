package component

import (
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/validations"
)

// ------------------------
// Aggregate Root: Component
// ------------------------

type Component struct {
	Name string
	Team string
}

func NewComponent(name, team string) *Component {
	return &Component{
		Name: name,
		Team: team,
	}
}

func (c *Component) GetName() string {
	return c.Name
}

func (c *Component) Validate() validations.ValidationResult {
	result := validations.ValidationResult{Passed: true}

	if c.Name == "" {
		result.Passed = false
		result.Errors = append(result.Errors, validations.ValidationError{
			EntityName:   c.Name,
			EntityType:   "component",
			ErrorCode:    validations.EmptyField,
			ErrorMessage: "component name cannot be empty",
		})
	}

	if c.Team == "" {
		result.Passed = false
		result.Errors = append(result.Errors, validations.ValidationError{
			EntityName:   c.Name,
			EntityType:   "component",
			ErrorCode:    validations.EmptyField,
			ErrorMessage: "component team cannot be empty",
		})
	}

	return result
}
