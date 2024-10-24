package endpoints

import (
	"emailn/internal/contract"
	internalmock "emailn/internal/test/mock"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CampaigsGetById_should_return_campaign(t *testing.T) {
	assert := assert.New(t)
	campaign := contract.CampaignResponse{
		ID:      "123",
		Name:    "Test",
		Content: "Hi",
		Status:  "Pending",
	}
	service := new(internalmock.CampaingServiceMock)
	service.On("GetByID", mock.Anything).Return(&campaign, nil)
	handler := Handler{CampaignService: service}
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	response, status, _ := handler.GetById(rr, req)

	assert.Equal(http.StatusOK, status)
	assert.Equal(campaign.ID, response.(*contract.CampaignResponse).ID)
	assert.Equal(campaign.Name, response.(*contract.CampaignResponse).Name)
}

func Test_CampaigsGetById_should_return_error_when_something_wrong_occurs(t *testing.T) {
	assert := assert.New(t)
	service := new(internalmock.CampaingServiceMock)
	errExpected := errors.New("something wrong")
	service.On("GetByID", mock.Anything).Return(nil, errExpected)
	handler := Handler{CampaignService: service}
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	_, _, errReturned := handler.GetById(rr, req)
	assert.Equal(errExpected.Error(), errReturned.Error())
}