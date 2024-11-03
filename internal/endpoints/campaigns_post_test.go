package endpoints

import (
	"emailn/internal/contract"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	body = contract.NewCampaign{
		Name:    "teste",
		Content: "hi everyone!",
		Emails:  []string{"teste@teste.com"},
	}
)

func Test_CampaigsPost_Status201(t *testing.T) {
	setup()
	createdByExpected := "teste@teste.com"
	service.On("Create", mock.MatchedBy(func(request contract.NewCampaign) bool {
		return request.Name == body.Name &&
			request.Content == body.Content &&
			request.CreatedBy == createdByExpected
	})).Return("123", nil)
	req, rr := newHttpTest(http.MethodPost, "/", body)
	req = addContext(req, "email", createdByExpected)

	_, status, err := handler.CreateCampaign(rr, req)

	assert.Equal(t, http.StatusCreated, status)
	assert.Nil(t, err)
}

func Test_CampaigsPost_Err(t *testing.T) {
	setup()
	service.On("Create", mock.Anything).Return("", fmt.Errorf("error"))

	req, rr := newHttpTest(http.MethodPost, "/", body)
	req = addContext(req, "email", "teste@teste.com")
	_, _, err := handler.CreateCampaign(rr, req)

	assert.NotNil(t, err)
}
