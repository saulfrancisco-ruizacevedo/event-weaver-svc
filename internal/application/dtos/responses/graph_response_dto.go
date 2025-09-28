package responses

type GraphResponseDto struct {
	Nodes         []GraphNodeDto         `json:"nodes"`
	Relationships []GraphRelationshipDto `json:"relationships"`
}

type GraphNodeDto struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
}

type GraphRelationshipDto struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	SourceID string `json:"sourceId"`
	TargetID string `json:"targetId"`
}
