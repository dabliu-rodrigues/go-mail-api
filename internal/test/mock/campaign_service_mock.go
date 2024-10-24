package mock

import (
	"emailn/internal/contract"

	"github.com/stretchr/testify/mock"
)

type CampaingServiceMock struct {
	mock.Mock
}

func (r *CampaingServiceMock) Create(newCampaign contract.NewCampaign) (string, error) {
	args := r.Called(newCampaign)
	return args.String(0), args.Error(1)
}

func (r *CampaingServiceMock) GetByID(id string) (*contract.CampaignResponse, error) {
	args := r.Called(id)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*contract.CampaignResponse), args.Error(1)
}
