package endpoints

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CampaignStart_Status200(t *testing.T) {
	setup()
	campaignId := "xpto"
	service.On("Start", mock.MatchedBy(func(id string) bool {
		return id == campaignId
	})).Return(nil)
	req, rr := newHttpTest(http.MethodPatch, "/", nil)
	req = addParameter(req, "id", campaignId)

	_, status, err := handler.CampaignStart(rr, req)

	assert.Equal(t, http.StatusOK, status)
	assert.Nil(t, err)
}

func Test_CampaignStart_Err(t *testing.T) {
	setup()
	errExpected := errors.New("something wrong")
	service.On("Start", mock.Anything).Return(errExpected)
	req, rr := newHttpTest(http.MethodPatch, "/", nil)

	_, _, err := handler.CampaignStart(rr, req)

	assert.Equal(t, errExpected, err)
}
