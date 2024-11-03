package endpoints

import (
	"bytes"
	"context"
	"emailn/internal/contract"
	internalmock "emailn/internal/test/internal-mock"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setup(body contract.NewCampaign, createdByExpected string) (*http.Request, *httptest.ResponseRecorder) {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)
	req, _ := http.NewRequest("POST", "/", &buf)
	ctx := context.WithValue(req.Context(), "email", createdByExpected)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	return req, rr
}

func Test_CampaigsPost_should_save_new_campaign(t *testing.T) {
	assert := assert.New(t)
	createdByExpected := "teste@teste.com"
	body := contract.NewCampaign{
		Name:    "teste",
		Content: "hi everyone!",
		Emails:  []string{"teste@teste.com"},
	}
	service := new(internalmock.CampaingServiceMock)
	service.On("Create", mock.MatchedBy(func(request contract.NewCampaign) bool {
		return request.Name == body.Name &&
			request.Content == body.Content &&
			request.CreatedBy == createdByExpected
	})).Return("123", nil)
	handler := Handler{CampaignService: service}
	req, rr := setup(body, createdByExpected)

	_, status, err := handler.CreateCampaign(rr, req)

	assert.Equal(http.StatusCreated, status)
	assert.Nil(err)
}

func Test_CampaigsPost_should_inform_error_when_exists(t *testing.T) {
	assert := assert.New(t)
	body := contract.NewCampaign{
		Name:    "teste",
		Content: "hi everyone!",
		Emails:  []string{"teste@teste.com"},
	}
	service := new(internalmock.CampaingServiceMock)
	service.On("Create", mock.Anything).Return("", fmt.Errorf("error"))
	handler := Handler{CampaignService: service}

	req, rr := setup(body, "teste@teste.com")
	_, _, err := handler.CreateCampaign(rr, req)

	assert.NotNil(err)
}
