package controller

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/sitetester/sequence-api/api"
	"github.com/sitetester/sequence-api/api/service"
	"gorm.io/gorm"
	"net/http"
)

type SequenceStepsController struct {
	SequenceService      service.SequenceService
	SequenceStepsService service.SequenceStepsService
}

func NewSequenceStepsController(db *gorm.DB) *SequenceStepsController {
	return &SequenceStepsController{
		SequenceService:      service.SequenceService{Db: db},
		SequenceStepsService: service.SequenceStepsService{Db: db},
	}
}

func (ssc *SequenceStepsController) Create(ctx *gin.Context) {
	var sequenceStep api.SequenceStep
	if err := ctx.BindJSON(&sequenceStep); err != nil {
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse{Error: err.Error()})
		return
	}

	_, err := govalidator.ValidateStruct(&sequenceStep)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse{Error: err.Error()})
		return
	}

	var foundSequence *api.Sequence
	foundSequence = ssc.SequenceService.GetByID(sequenceStep.SequenceID)
	if foundSequence.ID == 0 {
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse{Error: "Sequence not found."})
		return
	}

	// Assumption: steps have unique subject per sequence
	subjectAvailablePerSequence := ssc.SequenceStepsService.SubjectAvailablePerSequence(sequenceStep.Subject, sequenceStep.SequenceID)
	if !subjectAvailablePerSequence {
		ctx.JSON(http.StatusConflict, api.ErrorResponse{Error: "Subject already taken."})
		return
	}

	ssc.SequenceStepsService.Create(&sequenceStep)
	ctx.JSON(http.StatusCreated, &sequenceStep)
}

func (ssc *SequenceStepsController) Update(ctx *gin.Context) {
	stepIDStr := ctx.Param("id")
	stepID, err := api.StrToUint(stepIDStr)
	if err != nil {
		msg := fmt.Sprintf("Step ID must be integer: %s", stepIDStr)
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse{Error: msg})
		return
	}

	var foundSequenceStep *api.SequenceStep
	foundSequenceStep = ssc.SequenceStepsService.GetByID(uint(stepID))
	if foundSequenceStep.ID == 0 {
		ctx.JSON(http.StatusNotFound, api.ErrorResponse{Error: "Step not found."})
		return
	}

	var sequenceStep api.SequenceStep
	if err := ctx.BindJSON(&sequenceStep); err != nil {
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse{Error: err.Error()})
		return
	}
	_, err = govalidator.ValidateStruct(&sequenceStep)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse{Error: err.Error()})
		return
	}

	ssc.SequenceStepsService.Update(foundSequenceStep, sequenceStep)
}

func (ssc *SequenceStepsController) Delete(ctx *gin.Context) {
	stepIDStr := ctx.Param("id")
	stepID, err := api.StrToUint(stepIDStr)
	if err != nil {
		msg := fmt.Sprintf("Step ID must be integer: %s", stepIDStr)
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse{Error: msg})
		return
	}

	var foundSequenceStep *api.SequenceStep
	foundSequenceStep = ssc.SequenceStepsService.GetByID(uint(stepID))
	if foundSequenceStep.ID == 0 {
		ctx.JSON(http.StatusNotFound, api.ErrorResponse{Error: "Step not found."})
		return
	}

	ssc.SequenceStepsService.Delete(foundSequenceStep)
}

func (ssc *SequenceStepsController) View(ctx *gin.Context) {
	stepIDStr := ctx.Param("id")
	stepID, err := api.StrToUint(stepIDStr)
	if err != nil {
		msg := fmt.Sprintf("Step ID must be integer: %s", stepIDStr)
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse{Error: msg})
		return
	}

	var foundSequenceStep *api.SequenceStep
	foundSequenceStep = ssc.SequenceStepsService.GetByID(uint(stepID))
	if foundSequenceStep.ID == 0 {
		ctx.JSON(http.StatusNotFound, api.ErrorResponse{Error: "Step not found."})
		return
	}

	ctx.JSON(http.StatusOK, &foundSequenceStep)
}
