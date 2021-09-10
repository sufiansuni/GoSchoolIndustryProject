package models

type User struct {
	Username       string //primary key
	Password       []byte
	First          string
	Last           string
	Gender         string
	Birthday       string
	Height         int
	Weight         float64
	ActivityLevel  int
	CaloriesPerDay int
	Halal          bool
	Vegan          bool
	Address        string
	PostalCode     int
	Lat            float64
	Lng            float64
}
