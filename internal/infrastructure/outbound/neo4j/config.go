package neo4j

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v6/neo4j"

	"github.com/rs/zerolog/log"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/config"
)

func NewNeo4jDriver(cfg *config.Config) (neo4j.Driver, error) {
	driver, err := neo4j.NewDriver(cfg.Neo4jURI, neo4j.BasicAuth(cfg.Neo4jUser, cfg.Neo4jPassword, ""))
	if err != nil {
		panic(err)
	}

	session := driver.NewSession(context.Background(), neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeRead,
		DatabaseName: cfg.DBName,
	})
	defer session.Close(context.Background())

	if _, err := session.Run(context.Background(), "RETURN 1", nil); err != nil {
		return nil, err
	}

	log.Info().Msgf("Successfully connected to Neo4j database '%s'", cfg.DBName)
	return driver, nil
}
