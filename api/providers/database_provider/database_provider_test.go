package database_provider

import (
	stdcsv "encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gocarina/gocsv"
)

func Test_SaveAndList(t *testing.T) {
	provider := &DatabaseProvider{}
	item := InventoryItem{
		Quantity:     1,
		Name:         "Item 1",
		Manufacturer: "Manufacturer",
		Model:        "Model",
		Description:  "Description",
		SerialNumber: "Serial Number",
		Category:     "Category",
		PurchaseDate: "",
		Value:        100,
	}

	err := provider.Save(item)
	if err != nil {
		fmt.Errorf("got error %v, want nil", err)
	}
	items, err := provider.List()
	if err != nil {
		fmt.Errorf("got error %v, want nil", err)
	}
	if len(items) != 1 {
		t.Fatalf("got len(items) = %d, want 1", len(items))
	}
	if items[0] != item {
		t.Errorf("got item = %v, want %v", item, item)
	}
}

func Test_ListByCategory(t *testing.T) {
	provider := &DatabaseProvider{
		Repository: []InventoryItem{
			{Name: "Taylor 610CE", Category: "Guitars"},
			{Name: "Paradym Loud", Category: "Entertainment"},
		},
	}
	got, err := provider.ListByCategory("Guitars")
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	if len(got) != 1 {
		t.Fatalf("got len(items) = %d, want 1", len(got))
	}
	if got[0] != provider.Repository[0] {
		t.Errorf("got item = %v, want %v", got, provider.Repository[0])
	}
}

func Test_TotalByCategory(t *testing.T) {
	provider := &DatabaseProvider{
		Repository: []InventoryItem{
			{Name: "KRK Studio Monitors", Value: 179, Quantity: 2, Category: "Music Gear"},
			{Name: "Macbook Pro", Value: 2070, Quantity: 1, Category: "Computers"},
		},
	}
	want := 358
	got, err := provider.TotalByCategory("Music Gear")
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if got != want {
		t.Errorf("got = %d, want %d", got, want)
	}
}

func Test_ExportCSV(t *testing.T) {
	dirName := "tmp"
	csvPath := filepath.Join(ProjectDir, dirName, CSVFileName)

	provider := &DatabaseProvider{
		Repository: []InventoryItem{
			{Name: "KRK Studio Monitors", Value: 179, Quantity: 2, Category: "Music Gear"},
			{Name: "Macbook Pro", Value: 2070, Quantity: 1, Category: "Computers"},
		},
	}
	err := provider.ExportCSV("tmp")
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if err := checkCSVHeaders(csvPath, InventoryItem{}); err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	csvFilePtr, err := os.Open(csvPath)
	if err != nil {
		t.Fatalf("could not open file %s: %v", csvPath, err)
	}
	defer csvFilePtr.Close()

	var inventory []*InventoryItem
	if err := gocsv.UnmarshalFile(csvFilePtr, &inventory); err != nil {
		t.Fatalf("could not unmarshal CSV file %s: %v", csvPath, err)
	}
	if len(inventory) != 2 {
		t.Fatalf("got len(inventory) = %d, want 2", len(inventory))
	}
	wantValue := 358
	if inventory[0].Name != provider.Repository[0].Name {
		t.Errorf("got inventory[0].Name = %s, want %s", inventory[0].Name, provider.Repository[0].Name)
	}
	if inventory[0].Value != wantValue {
		t.Errorf("got inventory[0].Value = %d, want %d", inventory[0].Value, wantValue)
	}
	if inventory[1].Name != provider.Repository[1].Name {
		t.Errorf("got inventory[1].Name = %s, want %s", inventory[1].Name, provider.Repository[1].Name)
	}
	wantValue = 2070
	if inventory[1].Value != wantValue {
		t.Errorf("got inventory[1].Value = %d, want %d", inventory[1].Value, wantValue)
	}
}

// normalizeHeader canonicalizes a CSV header label for matching: surrounding
// whitespace is ignored and matching is case-insensitive.
func normalizeHeader(h string) string {
	return strings.ToLower(strings.TrimSpace(h))
}

// csvHeaders returns the `csv:` struct tags of a gocsv-mapped struct: the
// headers the bank's file must contain.
func csvHeaders(sample any) []string {
	t := reflect.TypeOf(sample)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	var headers []string
	for field := range t.Fields() {
		if tag, ok := field.Tag.Lookup("csv"); ok && tag != "" && tag != "-" {
			headers = append(headers, tag)
		}
	}
	return headers
}

// checkCSVHeaders reads a CSV file's header row and verifies it contains every
// header the given gocsv struct expects. gocsv maps columns by header name
// but stays silent when one is missing (the field just keeps its zero value),
// so a renamed bank column would silently zero a field; this makes it loud.
func checkCSVHeaders(csvFile string, sample any) error {
	f, err := os.Open(csvFile)
	if err != nil {
		return fmt.Errorf("could not open file %s: %w", csvFile, err)
	}
	defer f.Close()
	r := stdcsv.NewReader(f)
	header, err := r.Read()
	if err != nil {
		return fmt.Errorf("could not read header row of %s: %w", csvFile, err)
	}
	have := make(map[string]bool)
	for _, h := range header {
		have[normalizeHeader(h)] = true
	}
	var missing []string
	for _, want := range csvHeaders(sample) {
		if !have[normalizeHeader(want)] {
			missing = append(missing, want)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("CSV file %s is missing expected headers: %s", csvFile, strings.Join(missing, ", "))
	}
	return nil
}
