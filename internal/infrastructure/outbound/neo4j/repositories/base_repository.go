package repositories

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v6/neo4j"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/responses"
)

type BaseRepository struct {
	Driver neo4j.Driver
	DbName string
}

func (r *BaseRepository) ExecQuery(ctx context.Context, query string, params map[string]interface{}) ([]*neo4j.Record, error) {
	result, err := neo4j.ExecuteQuery(
		ctx,
		r.Driver,
		query,
		params,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(r.DbName),
	)

	if err != nil {
		return nil, err
	}

	return result.Records, nil
}

func (r *BaseRepository) GetAllNodes(ctx context.Context, nodeType string) (*responses.GraphResponseDto, error) {
	query := fmt.Sprintf(`MATCH (n:%s) RETURN n`, nodeType)

	records, err := r.ExecQuery(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve %s: %w", nodeType, err)
	}

	nodes := r.BuildNodes(records, "n")

	return &responses.GraphResponseDto{
		Nodes:         nodes,
		Relationships: []responses.GraphRelationshipDto{},
	}, nil
}

func (r *BaseRepository) BuildNodes(records []*neo4j.Record, keys ...string) []responses.GraphNodeDto {
	nodes := make([]responses.GraphNodeDto, 0)
	seen := make(map[string]bool)

	for _, record := range records {
		for _, key := range keys {
			val, ok := record.Get(key)
			if !ok || val == nil {
				continue
			}
			node, ok := val.(neo4j.Node)
			if !ok {
				continue
			}
			nodeID := fmt.Sprintf("%v", node.Props["name"])
			if !seen[nodeID] {
				props := make(map[string]interface{}, len(node.Props))
				for k, v := range node.Props {
					props[k] = v
				}
				nodes = append(nodes, responses.GraphNodeDto{
					ID:         nodeID,
					Type:       string(node.Labels[0]),
					Properties: props,
				})
				seen[nodeID] = true
			}
		}
	}

	return nodes
}

func (r *BaseRepository) BuildRelationships(records []*neo4j.Record, rels ...struct {
	SourceKey string
	TargetKey string
	Type      string
}) []responses.GraphRelationshipDto {
	relationships := []responses.GraphRelationshipDto{}

	for _, record := range records {
		for _, rel := range rels {
			sourceVal, ok1 := record.Get(rel.SourceKey)
			targetVal, ok2 := record.Get(rel.TargetKey)
			if !ok1 || !ok2 || sourceVal == nil || targetVal == nil {
				continue
			}
			sourceNode, ok1 := sourceVal.(neo4j.Node)
			targetNode, ok2 := targetVal.(neo4j.Node)
			if !ok1 || !ok2 {
				continue
			}

			relationships = append(relationships, responses.GraphRelationshipDto{
				ID:       fmt.Sprintf("%s-%s", sourceNode.Props["name"], targetNode.Props["name"]),
				Type:     rel.Type,
				SourceID: fmt.Sprintf("%v", sourceNode.Props["name"]),
				TargetID: fmt.Sprintf("%v", targetNode.Props["name"]),
			})
		}
	}

	return relationships
}
