package database_provider

import (
	stdcsv "encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"inventory/pkg/models"
)

const (
	ProjectDir  = "/Users/rob/ws/go/src/inventory"
	CSVFileName = "inventory.csv"
)

type DatabaseProviderInterface interface {
	Save(item models.InventoryItem) error
	List() ([]models.InventoryItem, error)
	ListByCategory(category string) ([]models.InventoryItem, error)
	TotalsByCategory(category string) (map[string]int, error)
	ExportCSV(fileName string) error
}
type DatabaseProvider struct {
	Repository []models.InventoryItem
}

func (p *DatabaseProvider) Save(item models.InventoryItem) error {
	p.Repository = append(p.Repository, item)
	return nil
}
func (p *DatabaseProvider) List() ([]models.InventoryItem, error) {
	return p.Repository, nil
}
func (p *DatabaseProvider) ListByCategory(category string) ([]models.InventoryItem, error) {
	var items []models.InventoryItem
	for _, item := range p.Repository {
		if item.Category == category {
			items = append(items, item)
		}
	}
	return items, nil
}
func (p *DatabaseProvider) TotalByCategory(category string) (int, error) {
	var totals = 0
	for _, item := range p.Repository {
		if item.Category == category {
			totals += item.Quantity * item.Value
		}
	}
	return totals, nil
}
func (p *DatabaseProvider) ExportCSV(dirName string) error {
	csvPath := filepath.Join(ProjectDir, dirName, CSVFileName)
	file, err := os.Create(csvPath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", csvPath, err)
	}
	defer file.Close()

	w := stdcsv.NewWriter(file)
	// Quantity	Name	Manufacturer	Model	Value	Category	Purchase Date	Serial Number
	record := []string{"Qty", "Name", "Manufacturer", "Model", "Value", "Category", "Purchase Date", "Serial Num"}
	if err := w.Write(record); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	for _, t := range p.Repository {
		record = []string{
			strconv.Itoa(t.Quantity),
			t.Name,
			t.Manufacturer,
			t.Model,
			strconv.Itoa(t.Value * t.Quantity),
			t.Category,
			t.PurchaseDate,
			t.SerialNumber,
		}
		if err := w.Write(record); err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}
