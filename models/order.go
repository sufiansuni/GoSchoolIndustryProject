package models

type Order struct {
	ID           int
	Username     string
	RestaurantID string
	Status       string // Cart > Awaiting collection > Completed
	Collection   string // Delivery or Self-Collect
	Date         string

	UserAddress string //User Location Data
	UserUnit    string
	UserLat     float64
	UserLng     float64

	RestaurantAddress string //Restaurant Location Data
	RestaurantUnit    string
	RestaurantLat     float64
	RestaurantLng     float64

	TotalCalories int
	BurnCalories  int
}
