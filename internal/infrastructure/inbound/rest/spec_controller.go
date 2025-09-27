package rest

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/commands"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/requests"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/application/dtos/responses"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/mediator"
	"gopkg.in/yaml.v3"
)

type SpecController struct {
}

func NewSpecController() *SpecController {
	return &SpecController{}
}

func (controller *SpecController) RegisterRoutes(router *gin.RouterGroup) {
	api := router.Group("/spec")
	api.POST("/validate/yml", controller.ValidateContract)
	api.POST("/save/yml", controller.SaveContract)
}

func (controller *SpecController) ValidateContract(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")

	if token != "Bearer "+os.Getenv("GITHUB_ACTION_SECRET") {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	yamlData, err := ctx.GetRawData()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body: " + err.Error()})
		return
	}

	var eventSpecificationValidationRequestDto requests.EventSpecificationRequestDto
	err = yaml.Unmarshal(yamlData, &eventSpecificationValidationRequestDto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := mediator.Send[*commands.SpecValidationCommand, *responses.EventSpecValidationResponseDto](
		commands.SpecValidationCommandName,
		commands.NewSpecValidationCommand(eventSpecificationValidationRequestDto),
		ctx,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (controller *SpecController) SaveContract(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")

	if token != "Bearer "+os.Getenv("GITHUB_ACTION_SECRET") {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	yamlData, err := ctx.GetRawData()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body: " + err.Error()})
		return
	}

	var eventSpecificationValidationRequestDto requests.EventSpecificationRequestDto
	err = yaml.Unmarshal(yamlData, &eventSpecificationValidationRequestDto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err = mediator.Send[*commands.SpecPersistenceCommand, any](
		commands.SpecPersictenceCommandName,
		commands.NewSpecPersistenceCommand(eventSpecificationValidationRequestDto),
		ctx,
	); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusCreated)
}
