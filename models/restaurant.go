package models

type Restaurant struct {
	ID          int //primary key
	Name        string
	Description string
	Halal       bool
	Vegan       bool
	Address     string
	Unit        string
	Lat         float64
	Lng         float64
}
