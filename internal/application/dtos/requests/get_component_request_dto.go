package requests

type GetComponentRequestDto struct {
	Name              string
	ProducesEvent     bool
	ConsumesEvent     bool
	SubscribesToTopic bool
	ProducesToTopic   bool
	ManagedByTeam     bool
}
