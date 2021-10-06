package models

type Food struct {
	ID           int //primary key
	RestaurantID int //foreign key
	Name         string
	Price        float64
	Calories     float64
}
