//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/commands"
	commandHandlers "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/commands/handlers"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/component"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/domainentity"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/event"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/team"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/domain/topic"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/config"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/inbound/rest"
	neo4jInfra "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/outbound/neo4j"
	neo4jRepositories "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/outbound/neo4j/repositories"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/mediator"
	"github.com/saulfrancisco-ruizacevedo/go-neopersist"
)

type App struct {
	SpecController *rest.SpecController
	Mediator       *MediatorInitializer
}

type MediatorInitializer struct{}

var ProviderSet = wire.NewSet(
	config.NewConfig,
	rest.NewSpecController,

	neo4jInfra.NewNeo4jDriver,
	neo4jInfra.NewNeo4jExecutor,
	neopersist.NewPersistenceManager,

	commandHandlers.NewSpecValidationCommandHandler,
	commandHandlers.NewSpecPersistenceCommandHandler,

	// Repos
	NewTeamRepository,
	NewComponentRepository,
	NewDomainRepository,
	NewTopicRepository,
	NewEventRepository,

	RegisterHandlers,
)

func InitializeMediator() (*App, error) {
	wire.Build(
		ProviderSet,

		wire.Struct(new(App), "*"),
	)
	return &App{}, nil
}

func RegisterHandlers(
	specHandler *commandHandlers.SpecValidationCommandHandler,
	specPersistenceHandler *commandHandlers.SpecPersistenceCommandHandler,
) *MediatorInitializer {
	mediator.Register(commands.SpecValidationCommandName, specHandler)
	mediator.Register(commands.SpecPersictenceCommandName, specPersistenceHandler)
	return &MediatorInitializer{}
}

func NewTeamRepository(manager *neopersist.PersistenceManager) (team.ITeamRepository, error) {
	return neo4jRepositories.NewTeamRepository(manager)
}

func NewComponentRepository(manager *neopersist.PersistenceManager) (component.IComponentRepository, error) {
	return neo4jRepositories.NewComponentRepository(manager)
}

func NewDomainRepository(manager *neopersist.PersistenceManager) (domainentity.IDomainRepository, error) {
	return neo4jRepositories.NewDomainRepository(manager)
}

func NewTopicRepository(manager *neopersist.PersistenceManager) (topic.ITopicRepository, error) {
	return neo4jRepositories.NewTopicRepository(manager)
}

func NewEventRepository(manager *neopersist.PersistenceManager) (event.IEventRepository, error) {
	return neo4jRepositories.NewEventRepository(manager)
}
