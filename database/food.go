package database

import (
	"GoIndustryProject/models"
	"context"
	"database/sql"
	"time"
)

// Operations for food database: Insert(Create), Select(Read), Update, Delete

// Insert a new restaurant entry into database
func InsertFood(myFood models.Food) error {
	statement := "INSERT INTO foods (RestaurantID, Name, Price, Calories) VALUES(?, ?, ?, ?)"
	_, err := DB.Exec(statement,
		myFood.RestaurantID,
		myFood.Name,
		myFood.Price,
		myFood.Calories,
	)
	if err != nil {
		return err
	}
	return nil
}

// Select/Read a food entry from database with a ID input
func SelectFood(ID int) (models.Food, error) {
	var myFood models.Food
	query := "SELECT * FROM foods WHERE ID=?"

	err := DB.QueryRow(query, ID).Scan(
		&myFood.ID,
		&myFood.RestaurantID,
		&myFood.Name,
		&myFood.Price,
		&myFood.Calories,
	)
	return myFood, err
}

// Update a restaurant entry in database
func UpdateFood(myFood models.Food) error {
	statement := "UPDATE foods SET RestaurantID=?, Name=?, Price =?, Calories =? " +
		"WHERE ID=?"

	_, err := DB.Exec(statement,
		myFood.RestaurantID,
		myFood.Name,
		myFood.Price,
		myFood.Calories,
		myFood.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

// Delete a restaurant entry in database
func DeleteFood(ID int) error {
	_, err := DB.Exec("DELETE FROM foods WHERE ID=?",
		ID)
	if err != nil {
		return err
	}
	return nil
}

// Select/Read All User Entries
func SelectAllFoodsByRestaurantID(db *sql.DB, restaurantID int) (myFoods []models.Food, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, "SELECT * FROM foods WHERE RestaurantID=?", restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var myFood models.Food
		err = rows.Scan(
			&myFood.ID,
			&myFood.RestaurantID,
			&myFood.Name,
			&myFood.Price,
			&myFood.Calories,
		)
		if err != nil {
			return nil, err
		}
		myFoods = append(myFoods, myFood)
	}

	return
}
