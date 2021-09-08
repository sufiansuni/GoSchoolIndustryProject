package main

// Operations for restaurant database: Insert(Create), Select(Read), Update, Delete

type restaurant struct {
	ID          string
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
func insertRestaurants(myRestaurants restaurant) error {
	statement := "INSERT INTO restaurants(ID, Name, Description, Halal, Vegan,Address, PostalCode, Lat, Lng) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(statement,
		myRestaurants.ID,
		myRestaurants.Name,
		myRestaurants.Description,
		myRestaurants.Halal,
		myRestaurants.Vegan,
		myRestaurants.Address,
		myRestaurants.PostalCode,
		myRestaurants.Lat,
		myRestaurants.Lng)
	if err != nil {
		return err
	}
	return nil
}

// Select/Read a restaurant entry from database with a ID input
func selectRestaurant(ID string) (restaurant, error) {
	var myRestaurants restaurant
	query := "SELECT * FROM restaurants WHERE ID=?"

	err := db.QueryRow(query, ID).Scan(
		&myRestaurants.ID,
		&myRestaurants.Name,
		&myRestaurants.Description,
		&myRestaurants.Halal,
		&myRestaurants.Vegan,
		&myRestaurants.Address,
		&myRestaurants.PostalCode,
		&myRestaurants.Lat,
		&myRestaurants.Lng,
	)
	return myRestaurants, err
}

// Update a restaurant entry in database

func updateRestaurant(myRestaurants restaurant) error {
	statement := "UPDATE restaurants SET Name=?, Description =?, Halal=?, Vegan=?, Address=?, PostalCode=?, Lat =?, Lng=? " +
		"WHERE ID=?"

	_, err := db.Exec(statement,
		myRestaurants.ID,
		myRestaurants.Name,
		myRestaurants.Description,
		myRestaurants.Halal,
		myRestaurants.Vegan,
		myRestaurants.Address,
		myRestaurants.PostalCode,
		myRestaurants.Lat,
		myRestaurants.Lng,
	)
	if err != nil {
		return err
	}
	return nil
}

// Delete a restaurant entry in database
func deleteRestaurant(ID string) error {
	_, err := db.Exec("DELETE FROM restaurants where ID=?",
		ID)
	if err != nil {
		return err
	}
	return nil
}
