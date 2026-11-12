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

func (r *mysqlQueryRepo) Save(item models.Item) error {
	return r.Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(item).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *mysqlQueryRepo) List() ([]models.Item, error) {
	var items []models.Item
	if result := r.Conn.Order("date").Find(&items); result.Error != nil {
		return nil, fmt.Errorf("unable to get items: %w", result.Error)
	}
	return items, nil
}

func (r *mysqlQueryRepo) ListByCategory(category string) ([]models.Item, error) {
	var items []models.Item
	if result := r.Conn.Where("category = ?", category).Order("name").Find(&items); result.Error != nil {
		return nil, fmt.Errorf("unable to get items for category: %w", result.Error)
	}
	return items, nil
}

func (r *mysqlQueryRepo) TotalByCategory(category string) (int, error) {
	var total int
	r.Conn.Model(&models.Item{}).
		Select("SUM(value * quantity)").
		Where("category = ?", category).
		Scan(&total)
	return total, nil
}

func (r *mysqlQueryRepo) SortedListForCSV() ([]models.Item, error) {
	var items []models.Item
	if result := r.Conn.Order("category").Order("name").Find(&items); result.Error != nil {
		return nil, fmt.Errorf("unable to get items: %w", result.Error)
	}
	return items, nil
}
