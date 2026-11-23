package database_provider

import (
	"inventory/pkg/config"
	"inventory/pkg/models"
	"inventory/pkg/repository"
	"inventory/pkg/repository/mysql"
)

type DatabaseProviderInterface interface {
	Save(item *models.Item) (*models.Item, error)
	Get(id int) (models.Item, error)
	SaveItems(items []models.Item) ([]models.Item, error)
	Delete(item models.Item) error
	List() ([]models.Item, error)
	ListByCategory(category string) ([]models.Item, error)
	TotalByCategory(category string) (int, error)
	// SortedListForOutput includes photo counts for inventory filtering.
	SortedListForOutput() ([]models.OutputItem, error)
	SavePhoto(itemPhoto *models.ItemPhoto) (*models.ItemPhoto, error)
	ListPhotos(itemID int) ([]models.ItemPhoto, error)
	DeletePhoto(photoID int) error
}
type DatabaseProvider struct {
	Repository repository.QueryRepo
}

func NewMySQLProvider(dbCfg *config.DBConfig) (DatabaseProviderInterface, error) {
	conn, err := repository.Connect(dbCfg)
	if err != nil {
		return nil, err
	}
	return &DatabaseProvider{Repository: mysql.NewMySQLQueryRepo(conn)}, nil
}

func (p *DatabaseProvider) Save(item *models.Item) (*models.Item, error) {
	return p.Repository.Save(item)
}
func (p *DatabaseProvider) Get(id int) (models.Item, error) { return p.Repository.Get(id) }
func (p *DatabaseProvider) SaveItems(items []models.Item) ([]models.Item, error) {
	return p.Repository.SaveItems(items)
}
func (p *DatabaseProvider) Delete(item models.Item) error {
	return p.Repository.Delete(item)
}
func (p *DatabaseProvider) List() ([]models.Item, error) {
	items, err := p.Repository.List()
	if err != nil {
		return nil, err
	}
	return items, nil
}
func (p *DatabaseProvider) ListByCategory(category string) ([]models.Item, error) {
	items, err := p.Repository.ListByCategory(category)
	if err != nil {
		return nil, err
	}
	return items, nil
}
func (p *DatabaseProvider) TotalByCategory(category string) (int, error) {
	totalValue, err := p.Repository.TotalByCategory(category)
	if err != nil {
		return 0, err
	}
	return totalValue, nil
}
func (p *DatabaseProvider) SortedListForOutput() ([]models.OutputItem, error) {
	items, err := p.Repository.SortedListForOutput()
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (p *DatabaseProvider) SavePhoto(itemPhoto *models.ItemPhoto) (*models.ItemPhoto, error) {
	return p.Repository.SavePhoto(itemPhoto)
}
func (p *DatabaseProvider) ListPhotos(itemID int) ([]models.ItemPhoto, error) {
	return p.Repository.ListPhotos(itemID)
}
func (p *DatabaseProvider) DeletePhoto(photoID int) error {
	return p.Repository.DeletePhoto(photoID)
}
