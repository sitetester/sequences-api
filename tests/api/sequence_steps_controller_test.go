package api

import (
	"encoding/json"
	"github.com/sitetester/sequence-api/api"
	"github.com/sitetester/sequence-api/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func checkStepByID(t *testing.T, url string, inputStep api.SequenceStep) {
	assertions := assert.New(t)
	recorder := performRequest(t, http.MethodGet, url, nil)
	checkStatusCode(t, http.StatusOK, recorder.Code)

	var stepResultByID *api.SequenceStep
	json.NewDecoder(recorder.Body).Decode(&stepResultByID)

	assertions.Equal(inputStep.SequenceID, stepResultByID.SequenceID)
	assertions.Equal(inputStep.Subject, stepResultByID.Subject)
	assertions.Equal(inputStep.Content, stepResultByID.Content)
}

func checkBindJonAndValidation(t *testing.T, method string, url string) {
	t.Run("FailsForBindJSON", func(t *testing.T) {
		inputStep := map[string]interface{}{
			"Subject": 123, // numeric input (must be in double quotes to consider as string)
			"Content": "abc",
		}
		checkFailsWithError(t, method, url, inputStep, http.StatusBadRequest, "json: cannot unmarshal")
	})

	t.Run("FailsForEmptySubject", func(t *testing.T) {
		inputStep := api.SequenceStep{
			Subject: "",
			Content: "blah contents",
		}
		checkFailsWithError(t, method, url, inputStep, http.StatusBadRequest, "Subject: non zero value required")
	})

	t.Run("FailsForContentMinLengthValidation", func(t *testing.T) {
		inputStep := api.SequenceStep{
			Subject: "blah",
			Content: "a",
		}
		checkFailsWithError(t, method, url, inputStep, http.StatusBadRequest, "minstringlength(3)")
	})
}

// Will run sequentially
func TestSequenceSteps(t *testing.T) {
	setupTestEnv()

	stepsUrl := config.ApiVersion + "/sequence-steps"
	baseSequence := api.Sequence{
		Name:                 "Sequence1",
		OpenTrackingEnabled:  false,
		ClickTrackingEnabled: true,
	}

	// `SequenceID` is set under `CreateSequence` stage
	baseStep := api.SequenceStep{
		Subject: "Step1",
		Content: "blah contents",
	}

	// let's make sure we have a sequence available in test db (as tests might run in parallel)
	t.Run("CreateSequence", func(t *testing.T) {
		deleteSequenceByName(baseSequence.Name)

		recorder := performRequest(t, http.MethodPost, config.ApiVersion+"/sequences", baseSequence)
		checkStatusCode(t, http.StatusCreated, recorder.Code)

		var sequenceResult *api.Sequence
		json.NewDecoder(recorder.Body).Decode(&sequenceResult)
		baseStep.SequenceID = sequenceResult.ID // assign newly generated ID
	})

	var newStepId uint

	t.Run("Create", func(t *testing.T) {
		checkBindJonAndValidation(t, http.MethodPost, stepsUrl)

		t.Run("FailsForNonExistingSequenceID", func(t *testing.T) {
			inputStep := baseStep
			inputStep.SequenceID = 0
			checkFailsWithError(t, http.MethodPost, stepsUrl, inputStep, http.StatusBadRequest, "Sequence not found.")
		})

		t.Run("Success", func(t *testing.T) {
			recorder := performRequest(t, http.MethodPost, stepsUrl, baseStep)
			checkStatusCode(t, http.StatusCreated, recorder.Code)

			var stepResult *api.SequenceStep
			json.NewDecoder(recorder.Body).Decode(&stepResult) // capture its ID (will be used in `update` tests below)

			// now check the "by ID" endpoint
			newStepId = stepResult.ID
			checkStepByID(t, buildUrl(stepsUrl, newStepId), baseStep)
		})

		t.Run("FailsWithDuplicateSubject", func(t *testing.T) {
			checkFailsWithError(t, http.MethodPost, stepsUrl, baseStep, http.StatusConflict, "Subject already taken.")
		})
	})

	t.Run("Update", func(t *testing.T) {
		updateStepUrl := buildUrl(stepsUrl, newStepId)

		t.Run("FailsForNonExistingStepID", func(t *testing.T) {
			checkFailsWih404(t, http.MethodPut, buildUrl(stepsUrl, 0))
		})

		checkBindJonAndValidation(t, http.MethodPut, updateStepUrl)

		t.Run("Success", func(t *testing.T) {
			inputStep := baseStep
			inputStep.Subject = "Test Subject 123"
			inputStep.Content = "Test Contents 456"

			recorder := performRequest(t, http.MethodPut, updateStepUrl, inputStep)
			checkStatusCode(t, http.StatusOK, recorder.Code)

			checkStepByID(t, updateStepUrl, inputStep)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		deleteStepUrl := buildUrl(stepsUrl, newStepId)
		t.Run("FailsForNonExistingStepID", func(t *testing.T) {
			checkFailsWih404(t, http.MethodDelete, buildUrl(stepsUrl, 0))
		})

		t.Run("Success", func(t *testing.T) {
			recorder := performRequest(t, http.MethodDelete, deleteStepUrl, nil)
			checkStatusCode(t, http.StatusOK, recorder.Code)

			// verify "by ID" returns 404
			checkFailsWih404(t, http.MethodDelete, deleteStepUrl)
		})
	})

}
