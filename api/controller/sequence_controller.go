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

type SequenceController struct {
	service service.SequenceService
}

func NewSequenceController(db *gorm.DB) *SequenceController {
	return &SequenceController{
		service: service.SequenceService{Db: db},
	}
}

func (sc *SequenceController) Create(ctx *gin.Context) {
	var sequence api.Sequence

	if err := ctx.BindJSON(&sequence); err != nil {
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse{Error: err.Error()})
		return
	}
	_, err := govalidator.ValidateStruct(&sequence)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse{Error: err.Error()})
		return
	}

	var foundSequence *api.Sequence
	foundSequence = sc.service.GetByName(sequence.Name)
	if foundSequence.ID > 0 {
		msg := fmt.Sprintf("Name already assigned to sequence: %d", foundSequence.ID)
		ctx.JSON(http.StatusConflict, api.ErrorResponse{Error: msg})
		return
	}

	sc.service.Create(&sequence)

	ctx.JSON(http.StatusCreated, &sequence)
}

func (sc *SequenceController) Update(ctx *gin.Context) {
	sequenceIDStr := ctx.Param("id")
	sequenceID, err := api.StrToUint(sequenceIDStr)
	if err != nil {
		msg := fmt.Sprintf("Sequence ID must be integer: %s", sequenceIDStr)
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse{Error: msg})
		return
	}

	var foundSequence *api.Sequence
	foundSequence = sc.service.GetByID(uint(sequenceID))
	if foundSequence.ID == 0 {
		ctx.JSON(http.StatusNotFound, api.ErrorResponse{Error: "Sequence not found."})
		return
	}

	var sequence api.Sequence
	if err := ctx.BindJSON(&sequence); err != nil {
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse{Error: err.Error()})
		return
	}
	_, err = govalidator.ValidateStruct(&sequence)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse{Error: err.Error()})
		return
	}

	// Check other sequence with same name
	var otherSequence *api.Sequence
	otherSequence = sc.service.GetOtherSequenceWithSameName(sequence.Name, uint(sequenceID))
	if otherSequence.ID > 0 {
		msg := fmt.Sprintf("Name already assigned to sequence: %d", otherSequence.ID)
		ctx.JSON(http.StatusConflict, api.ErrorResponse{Error: msg})
		return
	}

	// finally update
	sc.service.Update(foundSequence, sequence)
	// auto returns 200 status
}

func (sc *SequenceController) ViewWithSteps(ctx *gin.Context) {
	sequenceIDStr := ctx.Param("id")
	sequenceID, err := api.StrToUint(sequenceIDStr)
	if err != nil {
		msg := fmt.Sprintf("Sequence ID must be integer: %s", sequenceIDStr)
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse{Error: msg})
		return
	}

	var foundSequence *api.Sequence
	foundSequence = sc.service.GetWithSteps(sequenceID)
	if foundSequence.ID == 0 {
		ctx.JSON(http.StatusNotFound, api.ErrorResponse{Error: "Sequence not found."})
		return
	}

	sequenceWithSteps := api.SequenceWithSteps{
		Sequence: foundSequence,
		Steps:    &foundSequence.SequenceSteps,
	}
	ctx.JSON(http.StatusOK, sequenceWithSteps)
}

// TODO: Implement "delete" endpoint (when needed)
