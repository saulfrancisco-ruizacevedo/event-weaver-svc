package team

import (
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/validations"
)

// ------------------------
// Aggregate Root: Team
// ------------------------

type Team struct {
	Name  string
	Lead  string
	Email string
}

func NewTeam(name, lead, email string) *Team {
	return &Team{
		Name:  name,
		Lead:  lead,
		Email: email,
	}
}

func (t *Team) GetName() string {
	return t.Name
}

func (t *Team) Validate() validations.ValidationResult {
	result := validations.ValidationResult{Passed: true}

	if t.Name == "" {
		result.Passed = false
		result.Errors = append(result.Errors, validations.ValidationError{
			EntityName:   t.Name,
			EntityType:   "team",
			ErrorCode:    validations.EmptyField,
			ErrorMessage: "team name cannot be empty",
		})
	}
	if t.Lead == "" {
		result.Passed = false
		result.Errors = append(result.Errors, validations.ValidationError{
			EntityName:   t.Name,
			EntityType:   "team",
			ErrorCode:    validations.EmptyField,
			ErrorMessage: "team lead cannot be empty",
		})
	}
	if t.Email == "" {
		result.Passed = false
		result.Errors = append(result.Errors, validations.ValidationError{
			EntityName:   t.Name,
			EntityType:   "team",
			ErrorCode:    validations.EmptyField,
			ErrorMessage: "team email cannot be empty",
		})
	}

	return result
}
