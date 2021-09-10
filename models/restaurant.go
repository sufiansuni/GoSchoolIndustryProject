package models

type Restaurant struct {
	ID          int
	Name        string //primary key
	Description string
	Halal       bool
	Vegan       bool
	Address     string
	PostalCode  int
	Lat         float64
	Lng         float64
}

