package requests

type GetTopicRequestDto struct {
	Name                 string
	EventsOriginated     bool
	ComponentsSubscribed bool
	ComponentsProduced   bool
}
