package database

import (
	"GoIndustryProject/models"
	"context"
	"database/sql"
	"time"
)

// Operations for restaurant database: Insert(Create), Select(Read), Update, Delete

// Insert a new restaurant entry into database
func InsertRestaurant(myRestaurant models.Restaurant) error {
	statement := "INSERT INTO restaurants (Name, Description, Halal, Vegan, Address, Unit, Lat, Lng) VALUES(?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := DB.Exec(statement,
		myRestaurant.Name,
		myRestaurant.Description,
		myRestaurant.Halal,
		myRestaurant.Vegan,
		myRestaurant.Address,
		myRestaurant.Unit,
		myRestaurant.Lat,
		myRestaurant.Lng)
	if err != nil {
		return err
	}
	return nil
}

// Select/Read a restaurant entry from database with a ID input
func SelectRestaurant(ID int) (models.Restaurant, error) {
	var myRestaurant models.Restaurant
	query := "SELECT * FROM restaurants WHERE ID=?"

	err := DB.QueryRow(query, ID).Scan(
		&myRestaurant.ID,
		&myRestaurant.Name,
		&myRestaurant.Description,
		&myRestaurant.Halal,
		&myRestaurant.Vegan,
		&myRestaurant.Address,
		&myRestaurant.Unit,
		&myRestaurant.Lat,
		&myRestaurant.Lng,
	)
	return myRestaurant, err
}

// Update a restaurant entry in database
func UpdateRestaurant(myRestaurant models.Restaurant) error {
	statement := "UPDATE restaurants SET Name=?, Description =?, Halal=?, Vegan=?, Address=?, Unit=?, Lat =?, Lng=? " +
		"WHERE ID=?"

	_, err := DB.Exec(statement,
		myRestaurant.Name,
		myRestaurant.Description,
		myRestaurant.Halal,
		myRestaurant.Vegan,
		myRestaurant.Address,
		myRestaurant.Unit,
		myRestaurant.Lat,
		myRestaurant.Lng,
		myRestaurant.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

// Delete a restaurant entry in database
func DeleteRestaurant(ID int) error {
	_, err := DB.Exec("DELETE FROM restaurants WHERE ID=?",
		ID)
	if err != nil {
		return err
	}
	return nil
}

// Select/Read all restaurant entries
func SelectAllRestaurants(db *sql.DB) (myRestaurants []models.Restaurant, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, "SELECT * FROM restaurants")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var myRestaurant models.Restaurant
		err = rows.Scan(
			&myRestaurant.ID,
			&myRestaurant.Name,
			&myRestaurant.Description,
			&myRestaurant.Halal,
			&myRestaurant.Vegan,
			&myRestaurant.Address,
			&myRestaurant.Unit,
			&myRestaurant.Lat,
			&myRestaurant.Lng,
		)

		if err != nil {
			return nil, err
		}
		myRestaurants = append(myRestaurants, myRestaurant)
	}

	return
}
