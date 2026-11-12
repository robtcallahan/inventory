package database_provider

type DatabaseProviderInterface interface {
	Save(item InventoryItem) error
	List() ([]InventoryItem, error)
}
type DatabaseProvider struct {
	Repository []InventoryItem
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
