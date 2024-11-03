package campaign_test

import (
	"emailn/internal/contract"
	"emailn/internal/domain/campaign"
	internalerrors "emailn/internal/internal-errors"
	internalmock "emailn/internal/test/internal-mock"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var (
	newCampaign = contract.NewCampaign{
		Name:      "Test Y",
		Content:   "Body Hi!",
		Emails:    []string{"teste1@test.com"},
		CreatedBy: "teste@teste.com.br",
	}
	pendingCampaign *campaign.Campaign
	startedCampaign *campaign.Campaign
	repositoryMock  *internalmock.CampaignRepositoryMock
	service         = campaign.ServiceImp{}
)

func setUp() {
	repositoryMock = new(internalmock.CampaignRepositoryMock)
	service.Repository = repositoryMock
	pendingCampaign, _ = campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreatedBy)
	startedCampaign = &campaign.Campaign{ID: "1", Status: campaign.Status_Started}
}

func setUpGetByIdRepository(campaign *campaign.Campaign) {
	repositoryMock.On("GetByID", mock.Anything).Return(campaign, nil)
}

func setUpUpdateRepository() {
	repositoryMock.On("Update", mock.Anything).Return(nil)
}

func setUpSendMailSuccefully() {
	sendMail := func(campaign *campaign.Campaign) error {
		return nil
	}
	service.SendMail = sendMail
}

//Method_Context_ReturnOrAction

func Test_Create_Campaign(t *testing.T) {
	setUp()
	repositoryMock.On("Create", mock.Anything).Return(nil)

	id, err := service.Create(newCampaign)

	assert.NotNil(t, id)
	assert.Nil(t, err)
}

func Test_Create_ValidateDomainError(t *testing.T) {
	setUp()

	_, err := service.Create(contract.NewCampaign{})

	assert.False(t, errors.Is(err, internalerrors.ErrInternal))
}

func Test_Create_SaveCampaign(t *testing.T) {
	setUp()
	repositoryMock.On("Create", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		if campaign.Name != newCampaign.Name ||
			campaign.Content != newCampaign.Content ||
			len(campaign.Contacts) != len(newCampaign.Emails) {
			return false
		}

		return true
	})).Return(nil)

	service.Create(newCampaign)
	repositoryMock.AssertExpectations(t)
}

func Test_Create_ValidateRepositorySave(t *testing.T) {
	setUp()
	repositoryMock.On("Create", mock.Anything).Return(errors.New("error to save on database"))

	_, err := service.Create(newCampaign)

	assert.True(t, errors.Is(err, internalerrors.ErrInternal))
}

func Test_GetById_Return_Campaign_Response_Contract(t *testing.T) {
	setUp()
	repositoryMock.On("GetByID", mock.MatchedBy(func(id string) bool {
		return id == pendingCampaign.ID
	})).Return(pendingCampaign, nil)

	campaignReturned, _ := service.GetByID(pendingCampaign.ID)

	assert.Equal(t, pendingCampaign.ID, campaignReturned.ID)
	assert.Equal(t, pendingCampaign.Name, campaignReturned.Name)
	assert.Equal(t, pendingCampaign.Content, campaignReturned.Content)
	assert.Equal(t, pendingCampaign.CreatedBy, campaignReturned.CreatedBy)
	assert.Equal(t, string(pendingCampaign.Status), campaignReturned.Status)
}

func Test_GetById_Return_ErrorWhenSomethingWrongExists(t *testing.T) {
	setUp()
	repositoryMock.On("GetByID", mock.Anything).Return(pendingCampaign, errors.New("something wrong!"))

	_, err := service.GetByID("invalid_campaign_id")

	assert.Equal(t, internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Delete_ReturnRecordNotFound_When_Campaign_Does_Not_Exists(t *testing.T) {
	setUp()
	repositoryMock.On("GetByID", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	err := service.Delete("invalid_campaign_id")

	assert.Equal(t, err.Error(), gorm.ErrRecordNotFound.Error())
}

func Test_Delete_ReturnStatusInvalid_when_campaign_has_not_equals_pending(t *testing.T) {
	setUp()
	setUpGetByIdRepository(startedCampaign)

	err := service.Delete(startedCampaign.ID)

	assert.Equal(t, "Campaign status invalid", err.Error())
}

func Test_Delete_ReturnInternalError_when_delete_has_problems(t *testing.T) {
	setUp()
	setUpGetByIdRepository(pendingCampaign)
	repositoryMock.On("Delete", mock.Anything).Return(errors.New("error to delete campaign"))

	err := service.Delete(pendingCampaign.ID)

	assert.Equal(t, internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Delete_ReturnNil_when_delete_success(t *testing.T) {
	setUp()
	setUpGetByIdRepository(pendingCampaign)
	repositoryMock.On("Delete", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		return pendingCampaign == campaign
	})).Return(nil)

	err := service.Delete(pendingCampaign.ID)

	assert.Nil(t, err)
}

func Test_Start_ReturnRecordNotFound_When_Campaign_Does_Not_Exists(t *testing.T) {
	setUp()
	repositoryMock.On("GetByID", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	err := service.Start("invalid_campaign_id")

	assert.Equal(t, err.Error(), gorm.ErrRecordNotFound.Error())
}

func Test_Start_ReturnStatusInvalid_when_campaign_has_not_equals_pending(t *testing.T) {
	setUp()
	setUpGetByIdRepository(startedCampaign)

	err := service.Start(startedCampaign.ID)

	assert.Equal(t, "Campaign status invalid", err.Error())
}

func Test_Start_should_send_mail(t *testing.T) {
	setUp()
	setUpGetByIdRepository(pendingCampaign)
	setUpUpdateRepository()
	emailWasSent := false
	sendMail := func(campaign *campaign.Campaign) error {
		if campaign.ID == pendingCampaign.ID {
			emailWasSent = true
		}
		return nil
	}
	service.SendMail = sendMail

	service.Start(pendingCampaign.ID)

	assert.True(t, emailWasSent)
}

func Test_Start_ReturnError_when_func_SendMail_fails(t *testing.T) {
	setUp()
	setUpGetByIdRepository(pendingCampaign)
	sendMail := func(campaign *campaign.Campaign) error {
		return errors.New("error to send mail")
	}
	service.SendMail = sendMail

	err := service.Start(pendingCampaign.ID)

	assert.Equal(t, internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Start_ReturnNil_when_updated_to_done(t *testing.T) {
	setUp()
	setUpGetByIdRepository(pendingCampaign)
	setUpSendMailSuccefully()
	repositoryMock.On("Update", mock.MatchedBy(func(campaignToUpdate *campaign.Campaign) bool {
		return pendingCampaign.ID == campaignToUpdate.ID && campaignToUpdate.Status == campaign.Status_Done
	})).Return(nil)

	service.Start(pendingCampaign.ID)

	assert.Equal(t, campaign.Status_Done, pendingCampaign.Status)
}
