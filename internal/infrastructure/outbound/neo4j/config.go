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

	err = driver.VerifyConnectivity(context.Background())

	if err != nil {
		panic(err)

	} else {
		log.Info().Msg("Successfully connected to Neo4j database")
	}

	return driver, nil
}
