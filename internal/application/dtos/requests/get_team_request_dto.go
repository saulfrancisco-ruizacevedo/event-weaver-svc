package requests

type GetTeamRequestDto struct {
	Name             string
	ManagesComponent bool
	ProducesEvent    bool
	BelongsToDomain  bool
}
