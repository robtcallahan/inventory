package mysql

import (
	"fmt"
	"inventory/pkg/models"
	repo "inventory/pkg/repository"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type mysqlQueryRepo struct {
	Conn *gorm.DB
}

func NewMySQLQueryRepo(conn *gorm.DB) repo.QueryRepo {
	return &mysqlQueryRepo{
		Conn: conn,
	}
}

func (r *mysqlQueryRepo) Save(item models.InventoryItem) error {
	return r.Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(item).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *mysqlQueryRepo) List() ([]models.InventoryItem, error) {
	var items []models.InventoryItem
	if result := r.Conn.Order("date").Find(&items); result.Error != nil {
		return nil, fmt.Errorf("unable to get items: %w", result.Error)
	}
	return items, nil
}

func (r *mysqlQueryRepo) ListByCategory(category string) ([]models.InventoryItem, error) {
	var items []models.InventoryItem
	if result := r.Conn.Where("category = ?", category).Order("name").Find(&items); result.Error != nil {
		return nil, fmt.Errorf("unable to get items for category: %w", result.Error)
	}
	return items, nil
}

func (r *mysqlQueryRepo) TotalByCategory(category string) (int, error) {
	var results map[string]int
	if result := r.Conn.Select("SUM(value*quantity) as total_cost").Where("category = ?",
		category).Find(&results); result.Error != nil {
		return 0, fmt.Errorf("unable to get items for category: %w", result.Error)
	}
	return results["total_cost"], nil
}

func (r *mysqlQueryRepo) SortedListForCSV() ([]models.InventoryItem, error) {
	return nil, nil
}
