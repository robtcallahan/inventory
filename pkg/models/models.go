package models

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var ErrAlreadyExists = errors.New("item already exists")

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
	ID           int        `json:"id"`
	Quantity     int        `json:"quantity"`
	Name         string     `json:"name"`
	ReplacementValue        int        `json:"replacement_value"`
	ResaleValue  int        `json:"resale_value"`
	ShortDescription  string     `gorm:"type:varchar(255)" json:"short_description"`
	Manufacturer string     `json:"manufacturer"`
	Model        string     `json:"model"`
	Category     string     `json:"category"`
	PurchaseDate *time.Time `json:"purchase_date"`
	SerialNumber string     `json:"serial_number"`
	Room         string     `json:"room"`
	Condition    string     `json:"condition"`
	PurchaseSource string   `json:"purchase_source"`
	PurchasePrice int       `json:"purchase_price"`
	WebsiteURL   string     `json:"website_url"`
	ManualURL    string     `json:"manual_url"`
	Details      string     `gorm:"type:text" json:"details"`
}

type ItemPhoto struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	ItemID    int       `gorm:"not null;index" json:"item_id"`
	Filename  string    `gorm:"size:255;not null" json:"filename"`
	CreatedAt time.Time `json:"created_at"`
}

type OutputItem struct {
	ID           int        `csv:"id" json:"id"`
	Quantity     int        `csv:"Quantity" json:"quantity"`
	Name         string     `csv:"Name" json:"name"`
	ReplacementValue        int        `csv:"Replacement Value" json:"replacement_value"`
	ResaleValue  int        `csv:"Resale Value" json:"resale_value"`
	TotalValue   int        `csv:"Total Cost" json:"total_value"`
	ShortDescription  string     `csv:"Short Description" json:"short_description"`
	Manufacturer string     `csv:"Manufacturer" json:"manufacturer"`
	Model        string     `csv:"Model" json:"model"`
	Category     string     `csv:"Category" json:"category"`
	PurchaseDate *time.Time `csv:"Purchase Date" json:"purchase_date"`
	SerialNumber string     `csv:"Serial Number" json:"serial_number"`
	Room string `csv:"Room" json:"room"`
	Condition string `csv:"Condition" json:"condition"`
	PurchaseSource string `csv:"Purchase Source" json:"purchase_source"`
	PurchasePrice int `csv:"Purchase Price" json:"purchase_price"`
	WebsiteURL string `csv:"Website URL" json:"website_url"`
	ManualURL string `csv:"Manual URL" json:"manual_url"`
	Details      string     `csv:"Details" json:"details"`
	PhotoCount int `json:"photo_count" gorm:"column:photo_count" csv:"-"`
}

type HTMLItem struct {
	ID           int    `json:"id"`
	Number       int    `json:"number"`
	Quantity     int    `json:"quantity"`
	Name         string `json:"name"`
	ReplacementValue        int    `json:"replacement_value"`
	ResaleValue  int    `json:"resale_value"`
	TotalValue   int    `json:"total_value"`
	ShortDescription  string `json:"short_description"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	Category     string `json:"category"`
	PurchaseDate string `json:"purchase_date"`
	SerialNumber string `json:"serial_number"`
	Room string `json:"room"`
	Condition string `json:"condition"`
	PurchaseSource string `json:"purchase_source"`
	PurchasePrice int `json:"purchase_price"`
	WebsiteURL string `json:"website_url"`
	ManualURL string `json:"manual_url"`
	Details      string `json:"details"`
	PhotoCount int `json:"photo_count"`
}

type CSVItem struct {
	Quantity     int         `csv:"Quantity"`
	Name         string      `csv:"Name"`
	ReplacementValue        CurrencyInt `csv:"Replacement Value"`
	ResaleValue  CurrencyInt `csv:"Resale Value"`
	ShortDescription  string      `csv:"Short Description"`
	Manufacturer string      `csv:"Manufacturer"`
	Model        string      `csv:"Model"`
	Category     string      `csv:"Category"`
	PurchaseDate string      `csv:"Purchase Date"`
	SerialNumber string      `csv:"Serial Number"`
	Room string `csv:"Room"`
	Condition string `csv:"Condition"`
	PurchaseSource string `csv:"Purchase Source"`
	PurchasePrice CurrencyInt `csv:"Purchase Price"`
	WebsiteURL string `csv:"Website URL"`
	ManualURL string `csv:"Manual URL"`
	Details      string      `csv:"Details"`
}
