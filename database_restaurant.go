package main

// Operations for restaurant database: Insert(Create), Select(Read), Update, Delete

type restaurant struct {
	ID          int
	Name        string //primary key
	Description string
	Halal       bool
	Vegan       bool
	Address     string
	PostalCode  int
	Lat         float64
	Lng         float64
}

// Insert a new restaurant entry into database
func insertRestaurant(myRestaurant restaurant) error {
	statement := "INSERT INTO restaurants (Name, Description, Halal, Vegan, Address, PostalCode, Lat, Lng) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(statement,
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
func selectRestaurant(ID int) (restaurant, error) {
	var myRestaurant restaurant
	query := "SELECT * FROM restaurants WHERE ID=?"

	err := db.QueryRow(query, ID).Scan(
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

func updateRestaurant(myRestaurant restaurant) error {
	statement := "UPDATE restaurants SET Name=?, Description =?, Halal=?, Vegan=?, Address=?, PostalCode=?, Lat =?, Lng=? " +
		"WHERE ID=?"

	_, err := db.Exec(statement,
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
func deleteRestaurant(ID int) error {
	_, err := db.Exec("DELETE FROM restaurants WHERE ID=?",
		ID)
	if err != nil {
		return err
	}
	return nil
}
