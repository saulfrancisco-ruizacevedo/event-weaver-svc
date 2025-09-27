package domainentity

import "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/validations"

type Domain struct {
	Name        string
	Description string
}

func NewDomain(name, description string) *Domain {
	return &Domain{
		Name:        name,
		Description: description,
	}
}

func (d *Domain) GetName() string {
	return d.Name
}

func (d *Domain) Validate() validations.ValidationResult {
	result := validations.ValidationResult{Passed: true}

	if d.Name == "" {
		result.Passed = false
		result.Errors = append(result.Errors, validations.ValidationError{
			EntityName:   d.Name,
			EntityType:   "domain",
			ErrorCode:    validations.EmptyField,
			ErrorMessage: "domain name cannot be empty",
		})
	}
	if d.Description == "" {
		result.Passed = false
		result.Errors = append(result.Errors, validations.ValidationError{
			EntityName:   d.Name,
			EntityType:   "domain",
			ErrorCode:    validations.EmptyField,
			ErrorMessage: "domain description cannot be empty",
		})
	}

	return result
}
