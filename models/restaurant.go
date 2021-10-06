package models

type Restaurant struct {
	ID          int //primary key
	Name        string
	Description string
	Halal       bool
	Vegan       bool
	Address     string
	PostalCode  int
	Lat         float64
	Lng         float64
}
