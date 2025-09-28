//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/neo4j/neo4j-go-driver/v6/neo4j"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/commands"
	commandHandlers "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/commands/handlers"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/queries"
	queryHandlers "github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/queries/handlers"
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
)

type App struct {
	SpecController  *rest.SpecController
	GraphController *rest.GraphController
	Mediator        *MediatorInitializer
}

type MediatorInitializer struct{}

var ProviderSet = wire.NewSet(
	config.NewConfig,
	rest.NewSpecController,
	rest.NewGraphController,
	neo4jInfra.NewNeo4jDriver,

	commandHandlers.NewSpecValidationCommandHandler,
	commandHandlers.NewSpecPersistenceCommandHandler,

	queryHandlers.NewGetDomainsQueryHandler,
	queryHandlers.NewGetTeamsQueryHandler,
	queryHandlers.NewGetComponentsQueryHandler,
	queryHandlers.NewGetEventsQueryHandler,
	queryHandlers.NewGetTopicsQueryHandler,
	queryHandlers.NewGetDomainQueryHandler,
	queryHandlers.NewGetTeamQueryHandler,
	queryHandlers.NewGetComponentQueryHandler,
	queryHandlers.NewGetEventQueryHandler,
	queryHandlers.NewGetTopicQueryHandler,

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
	getDomainsQueryHandler *queryHandlers.GetDomainsQueryHandler,
	getTeamsQueryHandler *queryHandlers.GetTeamsQueryHandler,
	getComponentsQueryHandler *queryHandlers.GetComponentsQueryHandler,
	getEventsQueryHandler *queryHandlers.GetEventsQueryHandler,
	getTopicsQueryHandler *queryHandlers.GetTopicsQueryHandler,
	getDomainQueryHandler *queryHandlers.GetDomainQueryHandler,
	getTeamQueryHandler *queryHandlers.GetTeamQueryHandler,
	getComponentQueryHandler *queryHandlers.GetComponentQueryHandler,
	getEventQueryHandler *queryHandlers.GetEventQueryHandler,
	getTopicQueryHandler *queryHandlers.GetTopicQueryHandler,
) *MediatorInitializer {
	mediator.Register(commands.SpecValidationCommandName, specHandler)
	mediator.Register(commands.SpecPersictenceCommandName, specPersistenceHandler)

	mediator.Register(queries.GetDomainsQueryName, getDomainsQueryHandler)
	mediator.Register(queries.GetTeamsQueryName, getTeamsQueryHandler)
	mediator.Register(queries.GetComponentsQueryName, getComponentsQueryHandler)
	mediator.Register(queries.GetEventsQueryName, getEventsQueryHandler)
	mediator.Register(queries.GetTopicsQueryName, getTopicsQueryHandler)
	mediator.Register(queries.GetDomainQueryName, getDomainQueryHandler)
	mediator.Register(queries.GetTeamQueryName, getTeamQueryHandler)
	mediator.Register(queries.GetComponentQueryName, getComponentQueryHandler)
	mediator.Register(queries.GetEventQueryName, getEventQueryHandler)
	mediator.Register(queries.GetTopicQueryName, getTopicQueryHandler)

	return &MediatorInitializer{}
}

func NewTeamRepository(driver neo4j.Driver, cfg *config.Config) team.ITeamRepository {
	return &neo4jRepositories.TeamRepository{
		BaseRepository: &neo4jRepositories.BaseRepository{
			Driver: driver,
			DbName: cfg.DBName,
		},
	}
}

func NewComponentRepository(driver neo4j.Driver, cfg *config.Config) component.IComponentRepository {
	return &neo4jRepositories.ComponentRepository{
		BaseRepository: &neo4jRepositories.BaseRepository{
			Driver: driver,
			DbName: cfg.DBName,
		},
	}
}

func NewDomainRepository(driver neo4j.Driver, cfg *config.Config) domainentity.IDomainRepository {
	return &neo4jRepositories.DomainRepository{
		BaseRepository: &neo4jRepositories.BaseRepository{
			Driver: driver,
			DbName: cfg.DBName,
		},
	}
}

func NewTopicRepository(driver neo4j.Driver, cfg *config.Config) topic.ITopicRepository {
	return &neo4jRepositories.TopicRepository{
		BaseRepository: &neo4jRepositories.BaseRepository{
			Driver: driver,
			DbName: cfg.DBName,
		},
	}
}

func NewEventRepository(driver neo4j.Driver, cfg *config.Config) event.IEventRepository {
	return &neo4jRepositories.EventRepository{
		BaseRepository: &neo4jRepositories.BaseRepository{
			Driver: driver,
			DbName: cfg.DBName,
		},
	}
}
