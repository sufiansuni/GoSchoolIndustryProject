package models

type Food struct {
	ID           int //primary key
	RestaurantID int //foreign key
	Name         string
	Description  string
	Price        float64
	Calories     int
}
