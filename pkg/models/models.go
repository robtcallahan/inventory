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
	Quantity     int        `csv:"Quantity" json:"quantity"`
	Name         string     `csv:"Name" json:"name"`
	Value        int        `csv:"Value" json:"value"`
	TotalCost    int        `csv:"Total Cost" json:"total_cost"`
	Description  string     `csv:"Description" json:"description"`
	Manufacturer string     `csv:"Manufacturer" json:"manufacturer"`
	Model        string     `csv:"Model" json:"model"`
	Category     string     `csv:"Category" json:"category"`
	PurchaseDate *time.Time `csv:"Purchase Date" json:"purchase_date"`
	SerialNumber string     `csv:"Serial Number" json:"serial_number"`
	Details      string     `csv:"Details" json:"details"`
}

type HTMLItem struct {
	Number       string `json:"number"`
	Quantity     string `json:"quantity"`
	Name         string `json:"name"`
	Value        string `json:"value"`
	TotalCost    string `json:"total_cost"`
	Description  string `json:"description"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	Category     string `json:"category"`
	PurchaseDate string `json:"purchase_date"`
	SerialNumber string `json:"serial_number"`
	Details      string `json:"details"`
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
