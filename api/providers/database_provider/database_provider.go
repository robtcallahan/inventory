package database_provider

import (
	"inventory/pkg/models"
	"inventory/pkg/repository"
)

const (
	ProjectDir  = "/Users/rob/ws/go/src/inventory"
	CSVFileName = "inventory.csv"
)

type DatabaseProviderInterface interface {
	Save(item models.InventoryItem) error
	List() ([]models.InventoryItem, error)
	ListByCategory(category string) ([]models.InventoryItem, error)
	TotalByCategory(category string) (int, error)
	SortedListForCSV() ([]models.InventoryItem, error)
}
type DatabaseProvider struct {
	Repository repository.QueryRepo
}

func (p *DatabaseProvider) Save(item models.InventoryItem) error {
	return p.Repository.Save(item)
}
func (p *DatabaseProvider) List() ([]models.InventoryItem, error) {
	items, err := p.Repository.List()
	if err != nil {
		return nil, err
	}
	return items, nil
}
func (p *DatabaseProvider) ListByCategory(category string) ([]models.InventoryItem, error) {
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
func (p *DatabaseProvider) SortedListForCSV() ([]models.InventoryItem, error) {
	return nil, nil
}
