package database_provider

type DatabaseProviderInterface interface {
	Save(item InventoryItem) error
	List() ([]InventoryItem, error)
	ListByCategory(category string) ([]InventoryItem, error)
	TotalsByCategory(category string) (map[string]int, error)
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
