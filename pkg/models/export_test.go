package models

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/gocarina/gocsv"
)

type CSVOutputItem struct {
	Quantity     int    `csv:"Quantity"`
	Name         string `csv:"Name"`
	Value        int    `csv:"Replacement Value"`
	ResaleValue  int    `csv:"Resale Value"`
	TotalCost    int    `csv:"Total Cost"`
	Description  string `csv:"Short Description"`
	Manufacturer string `csv:"Manufacturer"`
	Model        string `csv:"Model"`
	Category     string `csv:"Category"`
	PurchaseDate string `csv:"Purchase Date"`
	SerialNumber string `csv:"Serial Number"`
	Details      string `csv:"Details"`
}

func Test_ItemsExportToCSVFile(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "inventory.csv")

	pianoDate, _ := time.Parse("01/02/2006", "07/07/2025")
	macDate, _ := time.Parse("01/02/2006", "10/12/2026")
	denonDate, _ := time.Parse("01/02/2006", "05/22/2021")

	items := []OutputItem{
		{Quantity: 1, Name: "Piano", Value: 2632, ResaleValue: 1800, TotalValue: 2632, Description: "RD-2000 EX Premium Digital Stage", PurchaseDate: &pianoDate},
		{Quantity: 1, Name: "Computer", Value: 2070, TotalValue: 2070, Description: "MacBook Pro 16-inch, M3 Pro, 36GB/512GB", PurchaseDate: &macDate},
		{Quantity: 1, Name: "AV Receiver", Value: 1907, TotalValue: 1907, Description: "AVR-X4700H", PurchaseDate: &denonDate},
		{Quantity: 2, Name: "Studio Monitor", Value: 179, TotalValue: 358, Description: "ROKIT 5 G4"},
	}
	err := ItemsExportToCSVFile(outputFile, items)
	if err != nil {
		t.Fatal(err)
	}
	// read the csv file
	outputFilePtr, err := os.Open(outputFile)
	if err != nil {
		t.Fatalf("could not open file %s: %v", outputFile, err)
	}

	reader := csv.NewReader(outputFilePtr)
	headers, err := reader.Read()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"Qty", "Name", "Replacement Value", "Resale Value", "Total Cost", "Short Description", "Manufacturer", "Model", "Category", "Purchase Date", "Serial Num", "Details"}
	if !slices.Equal(headers, want) {
		t.Errorf("got header %v, want %v", headers, want)
	}
	_ = outputFilePtr.Close()

	outputFilePtr, err = os.Open(outputFile)
	if err != nil {
		t.Fatalf("could not open file %s: %v", outputFile, err)
	}
	var readItems []CSVOutputItem
	if err := gocsv.UnmarshalFile(outputFilePtr, &readItems); err != nil {
		t.Fatalf("could not unmarshal CSV file %s: %v", outputFile, err)
	}
	if len(readItems) != len(items) {
		t.Fatalf("expected %d items, got %d", len(items), len(readItems))
	}
	if readItems[0].Name != items[0].Name || readItems[3].Name != items[3].Name {
		t.Errorf("expected items[0].Name: %v, items[3].Name: %v, got %v and %v", items[0].Name, items[3].Name,
			readItems[0].Name, readItems[3].Name)
	}
	if readItems[3].Value != items[3].Value || items[3].TotalValue != readItems[3].TotalCost {
		t.Errorf("expected items[3].Value: %v, items[3].TotalValue: %v, got %v & %v",
			items[3].Value, items[3].TotalValue, readItems[3].Value, readItems[3].TotalCost)
	}
	if readItems[0].ResaleValue != items[0].ResaleValue {
		t.Errorf("resale value round trip: got %d, want %d", readItems[0].ResaleValue, items[0].ResaleValue)
	}
	if readItems[3].PurchaseDate != "" {
		t.Errorf("expected items[3].PurchaseDate to be empty, got %v", readItems[3].PurchaseDate)
	}
	if readItems[0].PurchaseDate != "07/07/2025" || readItems[1].PurchaseDate != "10/12/2026" {
		t.Errorf("expected items[0].PurchaseDate: \"07/07/2025\", items[1].PurchaseDate: \"10/12/2026\", got %v & %v",
			readItems[0].PurchaseDate, readItems[1].PurchaseDate)
	}
	_ = outputFilePtr.Close()
}
