package main

// Operations for food database: Insert(Create), Select(Read), Update, Delete

type food struct {
	ID           int //primary key
	RestaurantID int //foreign key
	Name         string
	Price        float64
	Calories     float64
	Halal        bool
	Vegan        bool
}

// Insert a new restaurant entry into database
func insertFood(myFood food) error {
	statement := "INSERT INTO foods (RestaurantID, Name, Price, Calories, Halal, Vegan) VALUES(?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(statement,
		myFood.RestaurantID,
		myFood.Name,
		myFood.Price,
		myFood.Calories,
		myFood.Halal,
		myFood.Vegan,
	)
	if err != nil {
		return err
	}
	return nil
}

// Select/Read a food entry from database with a ID input
func selectFood(ID int) (food, error) {
	var myFood food
	query := "SELECT * FROM foods WHERE ID=?"

	err := db.QueryRow(query, ID).Scan(
		&myFood.ID,
		&myFood.Name,
		&myFood.RestaurantID,
		&myFood.Price,
		&myFood.Calories,
		&myFood.Halal,
		&myFood.Vegan,
	)
	return myFood, err
}

// Update a restaurant entry in database
func updateFood(myFood food) error {
	statement := "UPDATE foods SET RestaurantID=?, Name?, Price =?, Calories =?, Halal=?, Vegan=? " +
		"WHERE ID=?"

	_, err := db.Exec(statement,
		myFood.RestaurantID,
		myFood.Price,
		myFood.Calories,
		myFood.Halal,
		myFood.Vegan,
		myFood.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

// Delete a restaurant entry in database
func deleteFood(ID int) error {
	_, err := db.Exec("DELETE FROM foods WHERE ID=?",
		ID)
	if err != nil {
		return err
	}
	return nil
}
