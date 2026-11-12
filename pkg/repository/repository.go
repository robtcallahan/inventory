package repository

import (
	"inventory/pkg/models"
)

// QueryRepo represent the repositories
type QueryRepo interface {
	Save(item models.InventoryItem) error
	List() ([]models.InventoryItem, error)
	ListByCategory(category string) ([]models.InventoryItem, error)
	TotalByCategory(category string) (int, error)
	SortedListForCSV() ([]models.InventoryItem, error)
}
