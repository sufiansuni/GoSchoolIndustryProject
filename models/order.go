package models

type Order struct {
	ID             int
	Username       string
	RestaurantID   string
	RestaurantName string

	Status     string // Started > Awaiting Collection > Completed
	Collection string // Delivery or Self-Collect
	Date       string

	UserAddress string //User Location Data
	UserUnit    string
	UserLat     float64
	UserLng     float64

	RestaurantAddress string //Restaurant Location Data
	RestaurantUnit    string
	RestaurantLat     float64
	RestaurantLng     float64

	TotalItems int
	TotalPrice float64

	TotalCalories int
	BurnCalories  int
}
