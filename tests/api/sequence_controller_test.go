package api

import (
	"encoding/json"
	"fmt"
	"github.com/sitetester/sequence-api/api"
	"github.com/sitetester/sequence-api/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func checkByID(t *testing.T, url string, inputSequence api.Sequence) {
	recorder := performRequest(t, http.MethodGet, url, nil)
	checkStatusCode(t, http.StatusOK, recorder.Code)

	var sequenceWithSteps *api.SequenceWithSteps
	json.NewDecoder(recorder.Body).Decode(&sequenceWithSteps)
	assert.Equal(t, inputSequence.Name, sequenceWithSteps.Sequence.Name)
	assert.Equal(t, inputSequence.OpenTrackingEnabled, sequenceWithSteps.Sequence.OpenTrackingEnabled)
	assert.Equal(t, inputSequence.ClickTrackingEnabled, sequenceWithSteps.Sequence.ClickTrackingEnabled)
}

func checkBindJsonAndValidation(t *testing.T, method string, url string) {
	t.Run("FailsForJSONBinding", func(t *testing.T) {
		inputSequence := map[string]interface{}{
			"Name":                 "Sequence123",
			"OpenTrackingEnabled":  "abc", // a string value
			"ClickTrackingEnabled": true,
		}
		checkFailsWithError(t, method, url, inputSequence, http.StatusBadRequest, "json: cannot unmarshal")
	})

	t.Run("FailsForNameMinLengthValidation", func(t *testing.T) {
		inputSequence := api.Sequence{
			Name:                 "a",
			OpenTrackingEnabled:  false,
			ClickTrackingEnabled: true,
		}
		checkFailsWithError(t, method, url, inputSequence, http.StatusBadRequest, "minstringlength(3)")
	})
}

// Will run sequentially
func TestSequence(t *testing.T) {
	setupTestEnv()

	assertions := assert.New(t)

	sequencesUrl := config.ApiVersion + "/sequences"
	baseSequence := api.Sequence{
		Name:                 "Sequence1",
		OpenTrackingEnabled:  false,
		ClickTrackingEnabled: true,
	}

	var newSequenceID uint
	t.Run("Create", func(t *testing.T) {
		checkBindJsonAndValidation(t, http.MethodPost, sequencesUrl)

		t.Run("Success", func(t *testing.T) {
			deleteSequenceByName(baseSequence.Name)

			recorder := performRequest(t, http.MethodPost, sequencesUrl, baseSequence)
			checkStatusCode(t, http.StatusCreated, recorder.Code)

			var postResult *api.Sequence
			json.NewDecoder(recorder.Body).Decode(&postResult)
			newSequenceID = postResult.ID

			// now check the "by ID" endpoint
			checkByID(t, buildUrl(sequencesUrl, newSequenceID), baseSequence)
		})

		t.Run("FailsForDuplicateName", func(t *testing.T) {
			checkFailsWithError(t, http.MethodPost, sequencesUrl, baseSequence, http.StatusConflict, "Name already assigned")
		})
	})

	// CAUTION! This has dependency on ```postResult.ID``` (from `Create` step)
	t.Run("Update", func(t *testing.T) {
		updateUrl := buildUrl(sequencesUrl, newSequenceID)

		t.Run("FailsForNonExistingSequenceID", func(t *testing.T) {
			checkFailsWih404(t, http.MethodPut, buildUrl(sequencesUrl, 0))
		})

		checkBindJsonAndValidation(t, http.MethodPut, updateUrl)

		t.Run("FailsForDuplicateName(ForAnyOtherSequence)", func(t *testing.T) {
			// let's first create 2 sequences
			inputSequence1 := baseSequence
			inputSequence1.Name = "Sequence11"
			deleteSequenceByName(inputSequence1.Name)

			recorder := performRequest(t, http.MethodPost, sequencesUrl, inputSequence1)
			checkStatusCode(t, http.StatusCreated, recorder.Code)
			var postResult1 *api.Sequence
			json.NewDecoder(recorder.Body).Decode(&postResult1)

			inputSequence2 := baseSequence
			inputSequence2.Name = "Sequence22"
			deleteSequenceByName(inputSequence2.Name)
			recorder = performRequest(t, http.MethodPost, sequencesUrl, inputSequence2)
			checkStatusCode(t, http.StatusCreated, recorder.Code)
			var postResult2 *api.Sequence
			json.NewDecoder(recorder.Body).Decode(&postResult2) // capture it's ID

			// Now try to change inputSequence2 name to inputSequence1
			inputSequence2.Name = inputSequence1.Name
			url := buildUrl(sequencesUrl, postResult2.ID)
			recorder = performRequest(t, http.MethodPut, url, inputSequence2)
			checkStatusCode(t, http.StatusConflict, recorder.Code)
			resp := parseErrorResponse(recorder)
			msg := fmt.Sprintf("Name already assigned to sequence: %d", postResult1.ID)
			assertions.Contains(resp.Error, msg)
		})

		t.Run("Success", func(t *testing.T) {
			updateInput := api.Sequence{
				Name:                 "Sequence123",
				OpenTrackingEnabled:  true,
				ClickTrackingEnabled: false,
			}

			// delete the existing record (if any, when this test is run 2nd time)
			Db.Where("name = ? AND id != ? ", updateInput.Name, newSequenceID).Delete(&api.Sequence{})
			recorder := performRequest(t, http.MethodPut, updateUrl, updateInput)
			checkStatusCode(t, http.StatusOK, recorder.Code)

			// now check the "by ID" endpoint
			checkByID(t, updateUrl, updateInput)
		})
	})

	t.Run("ViewWithSteps", func(t *testing.T) {
		t.Run("FailsForNonExistingSequenceID", func(t *testing.T) {
			checkFailsWih404(t, http.MethodGet, buildUrl(sequencesUrl, 0))
		})

		// `Success` case was already covered in `Create` & `Update` tests above
	})
}
