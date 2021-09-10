package models

type Order struct {
	ID         int
	Username   string
	Status     string
	Date       string
	Address    string
	PostalCode int
	Lat        float64
	Lng        float64
}