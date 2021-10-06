package models

type OrderItem struct {
	ID               int //primary key
	OrderID          int
	FoodID           int
	FoodName         string
	Quantity         int
	SubtotalPrice    float64
	SubtotalCalories int
}
