package models

import (
	"testing"
	"time"
)

func Test_OutputItemsToHTML(t *testing.T) {
	guitarDate, _ := time.Parse("01/02/2006", "01/01/2006")
	monitorsDate, _ := time.Parse("01/02/2006", "01/01/2020")
	outputItems := []OutputItem{
		{
			ID:           1,
			Quantity:     1,
			Name:         "Guitar",
			Value:        2000,
			TotalValue:   2000,
			Description:  "Taylor 610CE",
			Manufacturer: "Taylor",
			Model:        "610CE",
			Category:     "Musical Instruments",
			PurchaseDate: &guitarDate,
			SerialNumber: "12123434",
			Details:      "Taylor 610CE",
		},
		{
			ID:           2,
			Quantity:     2,
			Name:         "Studio Monitors",
			Value:        400,
			TotalValue:   800,
			Description:  "ROKIT 5 G4",
			Manufacturer: "ROKIT",
			Model:        "5 G4",
			Category:     "Music Gear",
			PurchaseDate: &monitorsDate,
			SerialNumber: "050505050",
			Details:      "ROKIT 5 G4 studio monitors",
		},
	}
	htmlItems := OutputItemsToHTML(outputItems)

	if htmlItems[0].ID != 1 || htmlItems[1].ID != 2 {
		t.Errorf("OutputItemsToHTML returned wrong item ID %d, %d != %d, %d", htmlItems[0].ID, htmlItems[1].ID,
			outputItems[0].ID, outputItems[1].ID)
	}
	if htmlItems[0].Number != 1 || htmlItems[1].Number != 2 {
		t.Errorf("OutputItemsToHTML returned wrong item number %d, %d != %d, %d", htmlItems[0].Number, htmlItems[1].Number, 1, 2)
	}
	if htmlItems[0].PurchaseDate != "01/01/2006" || htmlItems[1].PurchaseDate != "01/01/2020" {
		t.Errorf("OutputItemsToHTML returned wrong item purchase date %s, %s != %s, %s", htmlItems[0].PurchaseDate,
			htmlItems[1].PurchaseDate, "01/01/2006", "01/01/2020")
	}
}
