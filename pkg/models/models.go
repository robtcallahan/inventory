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
	Quantity     int    `csv:"quantity"`
	Name         string `csv:"Name"`
	Manufacturer string `csv:"Manufacturer"`
	Model        string `csv:"Model"`
	Description  string `csv:"Description"`
	SerialNumber string `csv:"SerialNumber"`
	Category     string `csv:"Category"`
	PurchaseDate string `csv:"PurchaseDate"`
	Value        int    `csv:"Value"`
}
