package database

import "GoIndustryProject/models"

// Operations for restaurant database: Insert(Create), Select(Read), Update, Delete

// Insert a new restaurant entry into database
func InsertRestaurant(myRestaurant models.Restaurant) error {
	statement := "INSERT INTO restaurants (Name, Description, Halal, Vegan, Address, PostalCode, Lat, Lng) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := DB.Exec(statement,
		myRestaurant.Name,
		myRestaurant.Description,
		myRestaurant.Halal,
		myRestaurant.Vegan,
		myRestaurant.Address,
		myRestaurant.PostalCode,
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
		&myRestaurant.PostalCode,
		&myRestaurant.Lat,
		&myRestaurant.Lng,
	)
	return myRestaurant, err
}

// Update a restaurant entry in database

func UpdateRestaurant(myRestaurant models.Restaurant) error {
	statement := "UPDATE restaurants SET Name=?, Description =?, Halal=?, Vegan=?, Address=?, PostalCode=?, Lat =?, Lng=? " +
		"WHERE ID=?"

	_, err := DB.Exec(statement,
		myRestaurant.Name,
		myRestaurant.Description,
		myRestaurant.Halal,
		myRestaurant.Vegan,
		myRestaurant.Address,
		myRestaurant.PostalCode,
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
