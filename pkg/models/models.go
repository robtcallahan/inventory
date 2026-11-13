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

type OutputItem struct {
	Quantity     int
	Name         string
	Manufacturer string
	Model        string
	Description  string
	SerialNumber string
	Category     string
	PurchaseDate *time.Time
	Value        int
	TotalCost    int
}

type CSVItem struct {
	Quantity     int    `csv:"Quantity"`
	Name         string `csv:"Name"`
	Manufacturer string `csv:"Manufacturer"`
	Model        string `csv:"Model"`
	Description  string `csv:"Description"`
	SerialNumber string `csv:"Serial Number"`
	Category     string `csv:"Category"`
	PurchaseDate string `csv:"Purchase Date"`
	Value        int    `csv:"Value"`
}
