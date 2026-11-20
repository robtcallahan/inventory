package mysql

import (
	"fmt"
	"inventory/pkg/models"
	"inventory/pkg/repository"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type mysqlQueryRepo struct {
	Conn *gorm.DB
}

func NewMySQLQueryRepo(conn *gorm.DB) repository.QueryRepo {
	return &mysqlQueryRepo{
		Conn: conn,
	}
}

func (r *mysqlQueryRepo) Save(item models.Item) error {
	return r.Conn.Transaction(func(tx *gorm.DB) error {
		if item.ID == 0 {
			result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&item)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return models.ErrAlreadyExists
			}
			return nil
		}

		result := tx.Save(&item)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
}

func (r *mysqlQueryRepo) SaveItems(items []models.Item) ([]models.Item, error) {
	var skipped []models.Item
	err := r.Conn.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&item)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				skipped = append(skipped, item)
			}
		}
		return nil
	})
	return skipped, err
}

func (r *mysqlQueryRepo) List() ([]models.Item, error) {
	var items []models.Item
	if result := r.Conn.Order("name").Order("description").Find(&items); result.Error != nil {
		return nil, fmt.Errorf("unable to get items: %w", result.Error)
	}
	return items, nil
}

func (r *mysqlQueryRepo) ListByCategory(category string) ([]models.Item, error) {
	var items []models.Item
	if result := r.Conn.Where("category = ?", category).Order("name").Order("description").Find(&items); result.Error != nil {
		return nil, fmt.Errorf("unable to get items for category: %w", result.Error)
	}
	return items, nil
}

func (r *mysqlQueryRepo) TotalByCategory(category string) (int, error) {
	var total int
	if result := r.Conn.Model(&models.Item{}).
		Select("COALESCE(SUM(value * quantity), 0) AS total").
		Where("category = ?", category).
		Scan(&total); result.Error != nil {
		return 0, fmt.Errorf("unable to get total of items: %w", result.Error)
	}
	return total, nil
}

func (r *mysqlQueryRepo) SortedListForOutput() ([]models.OutputItem, error) {
	var items []models.OutputItem
	if result := r.Conn.Table("items").
		Select("*, (quantity * value) AS total_value").
		Order("total_value DESC").Order("name").Order("description").
		Find(&items); result.Error != nil {
		return nil, fmt.Errorf("unable to get items: %w", result.Error)
	}
	return items, nil
}
