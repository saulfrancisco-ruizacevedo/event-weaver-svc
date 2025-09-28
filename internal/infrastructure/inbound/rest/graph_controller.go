package rest

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/requests"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/responses"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/queries"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/mediator"
)

type GraphController struct{}

func NewGraphController() *GraphController {
	return &GraphController{}
}

func (controller *GraphController) RegisterRoutes(router *gin.RouterGroup) {
	api := router.Group("/graph")

	api.GET("/domains", controller.GetDomains)
	api.GET("/domains/:name", controller.GetDomain)

	api.GET("/teams", controller.GetTeams)
	api.GET("/teams/:name", controller.GetTeam)

	api.GET("/components", controller.GetComponents)
	api.GET("/components/:name", controller.GetComponent)

	api.GET("/events", controller.GetEvents)
	api.GET("/events/:name", controller.GetEvent)

	api.GET("/topics", controller.GetTopics)
	api.GET("/topics/:name", controller.GetTopic) // subscribedFromComponent, producedByComponent, originatedEvents
}

// Domains
func (controller *GraphController) GetDomains(ctx *gin.Context) {
	if !authorize(ctx) {
		return
	}

	response, err := mediator.Send[*queries.GetDomainsQuery, *responses.GraphResponseDto](
		queries.GetDomainsQueryName,
		queries.NewGetDomainsQuery(),
		ctx,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (controller *GraphController) GetDomain(ctx *gin.Context) {
	if !authorize(ctx) {
		return
	}

	opts := requests.GetDomainRequestDto{
		Name:                   ctx.Param("name"),
		RelatedToEvent:         ctx.Query("relatedToEvent") == "true",
		ComponentProducesEvent: ctx.Query("componentProducesEvent") == "true",
	}

	response, err := mediator.Send[*queries.GetDomainQuery, *responses.GraphResponseDto](
		queries.GetDomainQueryName,
		queries.NewGetDomainQuery(opts),
		ctx,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (controller *GraphController) GetTeams(ctx *gin.Context) {
	if !authorize(ctx) {
		return
	}

	response, err := mediator.Send[*queries.GetTeamsQuery, *responses.GraphResponseDto](
		queries.GetTeamsQueryName,
		queries.NewGetTeamsQuery(),
		ctx,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (controller *GraphController) GetTeam(ctx *gin.Context) {
	if !authorize(ctx) {
		return
	}

	opts := requests.GetTeamRequestDto{
		Name:             ctx.Param("name"),
		ManagesComponent: ctx.Query("managesComponent") == "true",
		ProducesEvent:    ctx.Query("producesEvent") == "true",
		BelongsToDomain:  ctx.Query("belongsToDomain") == "true",
	}

	response, err := mediator.Send[*queries.GetTeamQuery, *responses.GraphResponseDto](
		queries.GetTeamQueryName,
		queries.NewGetTeamQuery(opts),
		ctx,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (controller *GraphController) GetComponents(ctx *gin.Context) {
	if !authorize(ctx) {
		return
	}

	response, err := mediator.Send[*queries.GetComponentsQuery, *responses.GraphResponseDto](
		queries.GetComponentsQueryName,
		queries.NewGetComponentsQuery(),
		ctx,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (controller *GraphController) GetComponent(ctx *gin.Context) {
	if !authorize(ctx) {
		return
	}

	opts := requests.GetComponentRequestDto{
		Name:              ctx.Param("name"),
		ProducesEvent:     ctx.Query("producesEvent") == "true",
		ConsumesEvent:     ctx.Query("consumesEvent") == "true",
		SubscribesToTopic: ctx.Query("subscribesToTopic") == "true",
		ProducesToTopic:   ctx.Query("producesToTopic") == "true",
		ManagedByTeam:     ctx.Query("managedByTeam") == "true",
	}

	response, err := mediator.Send[*queries.GetComponentQuery, *responses.GraphResponseDto](
		queries.GetComponentQueryName,
		queries.NewGetComponentQuery(opts),
		ctx,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (controller *GraphController) GetEvents(ctx *gin.Context) {
	if !authorize(ctx) {
		return
	}

	response, err := mediator.Send[*queries.GetEventsQuery, *responses.GraphResponseDto](
		queries.GetEventsQueryName,
		queries.NewGetEventsQuery(),
		ctx,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (controller *GraphController) GetEvent(ctx *gin.Context) {
	if !authorize(ctx) {
		return
	}

	opts := requests.GetEventRequestDto{
		Name:                ctx.Param("name"),
		RelatedToEvent:      ctx.Query("relatedToEvent") == "true",
		BelongsToDomain:     ctx.Query("belongsToDomain") == "true",
		OriginatesFromTopic: ctx.Query("originatesFromTopic") == "true",
		ProducedByComponent: ctx.Query("producedByComponent") == "true",
		ConsumedByComponent: ctx.Query("consumedByComponent") == "true",
	}

	response, err := mediator.Send[*queries.GetEventQuery, *responses.GraphResponseDto](
		queries.GetEventQueryName,
		queries.NewGetEventQuery(opts),
		ctx,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// Topics
func (controller *GraphController) GetTopics(ctx *gin.Context) {
	if !authorize(ctx) {
		return
	}

	response, err := mediator.Send[*queries.GetTopicsQuery, *responses.GraphResponseDto](
		queries.GetTopicsQueryName,
		queries.NewGetTopicsQuery(),
		ctx,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (controller *GraphController) GetTopic(ctx *gin.Context) {
	if !authorize(ctx) {
		return
	}

	opts := requests.GetTopicRequestDto{
		Name:                 ctx.Param("name"),
		EventsOriginated:     ctx.Query("eventsOriginated") == "true",
		ComponentsSubscribed: ctx.Query("componentsSubscibed") == "true",
		ComponentsProduced:   ctx.Query("componentsProduced") == "true",
	}

	response, err := mediator.Send[*queries.GetTopicQuery, *responses.GraphResponseDto](
		queries.GetTopicQueryName,
		queries.NewGetTopicQuery(opts),
		ctx,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// ====================== HELPERS ======================

func authorize(ctx *gin.Context) bool {
	token := ctx.GetHeader("Authorization")
	if token != "Bearer "+os.Getenv("GITHUB_ACTION_SECRET") {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return false
	}
	return true
}
