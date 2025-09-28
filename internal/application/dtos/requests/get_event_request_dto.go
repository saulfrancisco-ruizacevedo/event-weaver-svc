package requests

type GetEventRequestDto struct {
	Name                string
	RelatedToEvent      bool
	ProducedByComponent bool
	ConsumedByComponent bool
	OriginatesFromTopic bool
	BelongsToDomain     bool
}
