package models

import (
	"strconv"
	"strings"
	"time"
)

// CurrencyInt custom type wrapping int
type CurrencyInt int

// UnmarshalCSV cleans the currency string before parsing it as an int
func (c *CurrencyInt) UnmarshalCSV(csv string) error {
	// Remove spaces, dollar signs, and thousands-separator commas
	clean := strings.TrimSpace(csv)
	clean = strings.TrimPrefix(clean, "$")
	clean = strings.ReplaceAll(clean, ",", "")
	clean = strings.TrimSpace(clean)

	// Handle empty cells or accounting dashes ("-") as 0
	if clean == "" || clean == "-" {
		*c = 0
		return nil
	}

	val, err := strconv.Atoi(clean)
	if err != nil {
		return err
	}

	*c = CurrencyInt(val)
	return nil
}

type Item struct {
	Quantity     int
	Name         string
	Value        int
	Description  string
	Manufacturer string
	Model        string
	Category     string
	PurchaseDate *time.Time
	SerialNumber string
	Details      string `csv:"Details"`
}

type OutputItem struct {
	Quantity     int        `csv:"Quantity"`
	Name         string     `csv:"Name"`
	Value        int        `csv:"Value"`
	TotalCost    int        `csv:"Total Cost"`
	Description  string     `csv:"Description"`
	Manufacturer string     `csv:"Manufacturer"`
	Model        string     `csv:"Model"`
	Category     string     `csv:"Category"`
	PurchaseDate *time.Time `csv:"Purchase Date"`
	SerialNumber string     `csv:"Serial Number"`
	Details      string     `csv:"Details"`
}

type CSVItem struct {
	Quantity     int         `csv:"Quantity"`
	Name         string      `csv:"Name"`
	Value        CurrencyInt `csv:"Value"`
	Description  string      `csv:"Description"`
	Manufacturer string      `csv:"Manufacturer"`
	Model        string      `csv:"Model"`
	Category     string      `csv:"Category"`
	PurchaseDate string      `csv:"Purchase Date"`
	SerialNumber string      `csv:"Serial Number"`
	Details      string      `csv:"Details"`
}
