package neo4j

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v6/neo4j"

	"github.com/rs/zerolog/log"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/config"
	"github.com/saulfrancisco-ruizacevedo/go-neopersist"
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

func NewNeo4jExecutor(cfg *config.Config) (neopersist.DBRunner, error) {
	executor, err := neopersist.NewNeo4jExecutor(cfg.Neo4jURI, cfg.Neo4jUser, cfg.Neo4jPassword, cfg.DBName)
	if err != nil {
		return nil, fmt.Errorf("could not create neo4j executor: %w", err)
	}

	if err := executor.Verify(context.Background()); err != nil {
		return nil, fmt.Errorf("could not verify neo4j connection: %w", err)
	}

	log.Info().Msgf("Successfully connected to Neo4j database '%s' via go-neopersist", cfg.DBName)
	return executor, nil
}
