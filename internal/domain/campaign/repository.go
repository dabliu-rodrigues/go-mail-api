package campaign

type Repository interface {
	Save(campaing *Campaign) error
	List() ([]Campaign, error)
	GetByID(id string) (*Campaign, error)
}
