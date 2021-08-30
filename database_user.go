package main

// Operations for User database: Insert, Delete, Update

//insert a new user entry
func insertUser(myUser user) error {

	//set default value for birthday if blank
	if myUser.Birthday == "" {
		myUser.Birthday = "1000-01-01"
	}

	statement := "INSERT INTO users (Username, Password, First, Last, Gender, Birthday, Height, Weight, CaloriesPerDay, Halal, Vegan, Address, PostalCode, Lat, Lng) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_, err := db.Exec(statement,
		myUser.Username,
		myUser.Password,
		myUser.First,
		myUser.Last,
		myUser.Gender,
		myUser.Birthday,
		myUser.Height,
		myUser.Weight,
		myUser.CaloriesPerDay,
		myUser.Halal,
		myUser.Vegan,
		myUser.Address,
		myUser.PostalCode,
		myUser.Lat,
		myUser.Lng)
	if err != nil {
		return err
	}
	return nil
}
