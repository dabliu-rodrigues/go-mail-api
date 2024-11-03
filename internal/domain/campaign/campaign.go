package campaign

import (
	internalerrors "emailn/internal/internal-errors"
	"time"

	"github.com/rs/xid"
)

type Contact struct {
	ID         string `gorm:"size:50"`
	Email      string `validate:"email" gorm:"size:100"`
	CampaignId string `gorm:"size:50"`
}

type Status string

const (
	Status_Pending  Status = "AWAITING"
	Status_Started  Status = "STARTED"
	Status_Done     Status = "DONE"
	Status_Canceled Status = "CANCELED"
	Status_Deleted  Status = "DELETED"
)

type Campaign struct {
	ID        string    `validate:"required" gorm:"size:50"`
	Name      string    `validate:"min=5,max=24" gorm:"size:100"`
	CreatedAt time.Time `validate:"required"`
	Content   string    `validate:"min=5,max=1024" gorm:"size:1024"`
	Contacts  []Contact `validate:"min=1,dive"`
	Status    Status    `gorm:"size:20"`
	CreatedBy string    `validate:"email" gorm:"size:50"`
}

func (c *Campaign) Cancel() {
	c.Status = Status_Canceled
}

func (c *Campaign) Delete() {
	c.Status = Status_Deleted
}

func NewCampaign(name, content string, emails []string, createdBy string) (*Campaign, error) {
	contacts := make([]Contact, len(emails))
	for i, email := range emails {
		contacts[i].Email = email
		contacts[i].ID = xid.New().String()
	}

	campaign := &Campaign{
		ID:        xid.New().String(),
		Name:      name,
		Content:   content,
		CreatedAt: time.Now(),
		Contacts:  contacts,
		Status:    Status_Pending,
		CreatedBy: createdBy,
	}

	err := internalerrors.ValidateStruct(campaign)
	if err == nil {
		return campaign, nil
	}

	return nil, err
}
