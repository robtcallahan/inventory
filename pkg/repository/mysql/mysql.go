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

func (r *mysqlQueryRepo) Save(item *models.Item) (*models.Item, error) {
	if item.ID == 0 {
		result := r.Conn.Clauses(clause.OnConflict{DoNothing: true}).Create(&item)
		if result.Error != nil {
			return nil, result.Error
		}
		if result.RowsAffected == 0 {
			return nil, models.ErrAlreadyExists
		}
		return item, nil
	}

	result := r.Conn.Save(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return item, nil
}

func (r *mysqlQueryRepo) Get(id int) (models.Item, error) {
	var item models.Item
	result := r.Conn.First(&item, id)
	if result.Error != nil {
		return models.Item{}, result.Error
	}
	return item, nil
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

// BulkUpdate locks all requested rows and updates only one whitelisted column atomically.
func (r *mysqlQueryRepo) BulkUpdate(ids []int, field string, value any) error {
    allowed := map[string]bool{"room": true, "category": true, "condition": true, "purchase_source": true, "website_url": true, "manual_url": true, "info_complete": true}
    if !allowed[field] || len(ids) == 0 {
        return fmt.Errorf("invalid bulk update field or IDs")
    }
    return r.Conn.Transaction(func(tx *gorm.DB) error {
        var found []models.Item
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id IN ?", ids).Find(&found).Error; err != nil {
            return err
        }
        if len(found) != len(ids) {
            return gorm.ErrRecordNotFound
        }
        return tx.Model(&models.Item{}).Where("id IN ?", ids).Update(field, value).Error
    })
}

func (r *mysqlQueryRepo) Delete(item models.Item) error {
	result := r.Conn.Where("id = ?", item.ID).Delete(item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *mysqlQueryRepo) List() ([]models.Item, error) {
	var items []models.Item
	if result := r.Conn.Order("name").Order("short_description").Find(&items); result.Error != nil {
		return nil, fmt.Errorf("unable to get items: %w", result.Error)
	}
	return items, nil
}

func (r *mysqlQueryRepo) ListByCategory(category string) ([]models.Item, error) {
	var items []models.Item
	if result := r.Conn.Where("category = ?", category).Order("name").Order("short_description").Find(&items); result.Error != nil {
		return nil, fmt.Errorf("unable to get items for category: %w", result.Error)
	}
	return items, nil
}

func (r *mysqlQueryRepo) TotalByCategory(category string) (int, error) {
	var total int
	if result := r.Conn.Model(&models.Item{}).
		Select("COALESCE(SUM(replacement_value * quantity), 0) AS total").
		Where("category = ?", category).
		Scan(&total); result.Error != nil {
		return 0, fmt.Errorf("unable to get total of items: %w", result.Error)
	}
	return total, nil
}

func (r *mysqlQueryRepo) SortedListForOutput() ([]models.OutputItem, error) {
	var items []models.OutputItem
	if result := r.Conn.Table("items").
		Select("items.*, (items.quantity * items.replacement_value) AS total_value, COALESCE(photo_counts.photo_count, 0) AS photo_count").
		Joins("LEFT JOIN (SELECT item_id, COUNT(*) AS photo_count FROM item_photos GROUP BY item_id) AS photo_counts ON photo_counts.item_id = items.id").
		Order("total_value DESC").Order("name").Order("short_description").
		Find(&items); result.Error != nil {
		return nil, fmt.Errorf("unable to get items: %w", result.Error)
	}
	return items, nil
}

func (r *mysqlQueryRepo) SavePhoto(itemPhoto *models.ItemPhoto) (*models.ItemPhoto, error) {
	if err := r.Conn.Create(itemPhoto).Error; err != nil {
		return nil, err
	}
	return itemPhoto, nil
}

func (r *mysqlQueryRepo) ListPhotos(itemID int) ([]models.ItemPhoto, error) {
	var photos []models.ItemPhoto

	result := r.Conn.
		Where("item_id = ?", itemID).
		Order("created_at ASC, id ASC").
		Find(&photos)

	if result.Error != nil {
		return nil, result.Error
	}

	return photos, nil
}
func (r *mysqlQueryRepo) DeletePhoto(photoID int) error {
	return r.Conn.
		Where("id = ?", photoID).
		Delete(&models.ItemPhoto{}).Error
}
