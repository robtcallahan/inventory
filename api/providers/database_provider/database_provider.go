package database_provider

import (
	"inventory/pkg/config"
	"inventory/pkg/models"
	"inventory/pkg/repository"
	"inventory/pkg/repository/mysql"
)

type DatabaseProviderInterface interface {
	Save(item models.Item) error
	SaveItems(items []models.Item) error
	List() ([]models.Item, error)
	ListByCategory(category string) ([]models.Item, error)
	TotalByCategory(category string) (int, error)
	SortedListForOutput() ([]models.OutputItem, error)
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

func (p *DatabaseProvider) Save(item models.Item) error {
	return p.Repository.Save(item)
}
func (p *DatabaseProvider) SaveItems(items []models.Item) error {
	return p.Repository.SaveItems(items)
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
