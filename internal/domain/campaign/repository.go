package campaign

type Repository interface {
	Create(campaing *Campaign) error
	List() ([]Campaign, error)
	GetByID(id string) (*Campaign, error)
	Delete(campaing *Campaign) error
	Update(campaign *Campaign) error
}
