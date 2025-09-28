package requests

type GetDomainRequestDto struct {
	Name                   string
	RelatedToEvent         bool
	ComponentProducesEvent bool
}
