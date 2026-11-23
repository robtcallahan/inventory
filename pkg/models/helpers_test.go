package models

import (
	"reflect"
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
			ResaleValue:  1700,
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
			ResaleValue:  250,
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

	want := []HTMLItem{
		{
			ID:           1,
			Number:       1,
			Quantity:     1,
			Name:         "Guitar",
			Value:        2000,
			TotalValue:   2000,
			Description:  "Taylor 610CE",
			Manufacturer: "Taylor",
			Model:        "610CE",
			Category:     "Musical Instruments",
			PurchaseDate: "01/01/2006",
			SerialNumber: "12123434",
			Details:      "Taylor 610CE",
		},
		{
			ID:           2,
			Number:       2,
			Quantity:     2,
			Name:         "Studio Monitors",
			Value:        400,
			TotalValue:   800,
			Description:  "ROKIT 5 G4",
			Manufacturer: "ROKIT",
			Model:        "5 G4",
			Category:     "Music Gear",
			PurchaseDate: "01/01/2020",
			SerialNumber: "050505050",
			Details:      "ROKIT 5 G4 studio monitors",
		},
	}

	htmlItems := OutputItemsToHTML(outputItems)

	if len(htmlItems) != len(outputItems) {
		t.Fatalf("expected %d items, got %d", len(outputItems), len(htmlItems))
	}
	if !reflect.DeepEqual(htmlItems, want) {
		t.Errorf("got %#v, want %#v", htmlItems, want)
	}

	outputItems[0].PurchaseDate = nil
	want[0].PurchaseDate = ""
	htmlItems = OutputItemsToHTML(outputItems)
	if !reflect.DeepEqual(htmlItems, want) {
		t.Errorf("got %#v, want %#v", htmlItems, want)
	}
}
