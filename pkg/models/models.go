package models

import "time"

type Item struct {
	Quantity     int
	Name         string
	Manufacturer string
	Model        string
	Description  string
	SerialNumber string
	Category     string
	PurchaseDate *time.Time
	Value        int
}

type CSVItem struct {
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
