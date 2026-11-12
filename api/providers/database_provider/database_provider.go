package database_provider

import (
	stdcsv "encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

const (
	ProjectDir  = "/Users/rob/ws/go/src/inventory"
	CSVFileName = "inventory.csv"
)

type DatabaseProviderInterface interface {
	Save(item InventoryItem) error
	List() ([]InventoryItem, error)
	ListByCategory(category string) ([]InventoryItem, error)
	TotalsByCategory(category string) (map[string]int, error)
	ExportCSV(fileName string) error
}
type DatabaseProvider struct {
	Repository []InventoryItem
}

func (p *DatabaseProvider) Save(item InventoryItem) error {
	p.Repository = append(p.Repository, item)
	return nil
}
func (p *DatabaseProvider) List() ([]InventoryItem, error) {
	return p.Repository, nil
}
func (p *DatabaseProvider) ListByCategory(category string) ([]InventoryItem, error) {
	var items []InventoryItem
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

type InventoryItem struct {
	Quantity     int
	Name         string
	Manufacturer string
	Model        string
	Description  string
	SerialNumber string
	Category     string
	PurchaseDate string
	Value        int
}
