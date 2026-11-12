package repository

import (
	"inventory/pkg/models"
)

// QueryRepo represent the repositories
type QueryRepo interface {
	Save(item models.InventoryItem) error
	List() ([]models.InventoryItem, error)
	ListByCategory(category string) ([]models.InventoryItem, error)
	TotalsByCategory(category string) (map[string]int, error)
	ExportCSV(fileName string) error
}
