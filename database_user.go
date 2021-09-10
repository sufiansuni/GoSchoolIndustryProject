package main

import "golang.org/x/crypto/bcrypt"

// Operations for users database: Insert(Create), Select(Read), Update, Delete

// Insert a new user entry into database
func insertUser(myUser user) error {

	// set default value for birthday if blank
	if myUser.Birthday == "" {
		myUser.Birthday = "1000-01-01"
	}

	// set default value for activity level if blank
	if myUser.ActivityLevel == 0 {
		myUser.ActivityLevel = 1
	}

	statement := "INSERT INTO users (Username, Password, First, Last, Gender, Birthday, Height, Weight, ActivityLevel, CaloriesPerDay, Halal, Vegan, Address, PostalCode, Lat, Lng) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_, err := db.Exec(statement,
		myUser.Username,
		myUser.Password,
		myUser.First,
		myUser.Last,
		myUser.Gender,
		myUser.Birthday,
		myUser.Height,
		myUser.Weight,
		myUser.ActivityLevel,
		myUser.CaloriesPerDay,
		myUser.Halal,
		myUser.Vegan,
		myUser.Address,
		myUser.PostalCode,
		myUser.Lat,
		myUser.Lng,
	)
	if err != nil {
		return err
	}
	return nil
}

// Select/Read a user entry from database with a username input
func selectUser(username string) (user, error) {
	var myUser user
	query := "SELECT * FROM users WHERE Username=?"

	err := db.QueryRow(query, username).Scan(
		&myUser.Username,
		&myUser.Password,
		&myUser.First,
		&myUser.Last,
		&myUser.Gender,
		&myUser.Birthday,
		&myUser.Height,
		&myUser.Weight,
		&myUser.ActivityLevel,
		&myUser.CaloriesPerDay,
		&myUser.Halal,
		&myUser.Vegan,
		&myUser.Address,
		&myUser.PostalCode,
		&myUser.Lat,
		&myUser.Lng,
	)
	return myUser, err
}

// Update a user entry in database
// Does not include username and password
func updateUserProfile(myUser user) error {
	statement := "UPDATE users SET First=?, Last=?, Gender=?, Birthday=?, " +
		"Height=?, Weight=?, ActivityLevel=?, CaloriesPerday=?, Halal=?, Vegan=?, Address=?, PostalCode=?,  Lat=?, Lng=? " +
		"WHERE Username=?"

	_, err := db.Exec(statement,
		myUser.First,
		myUser.Last,
		myUser.Gender,
		myUser.Birthday,
		myUser.Height,
		myUser.Weight,
		myUser.ActivityLevel,
		myUser.CaloriesPerDay,
		myUser.Halal,
		myUser.Vegan,
		myUser.Address,
		myUser.PostalCode,
		myUser.Lat,
		myUser.Lng,
		myUser.Username,
	)
	if err != nil {
		return err
	}
	return nil
}

// Update a users password entry
func updateUserPassword(username string, oldPassword string, newPassword string) error {
	//find stored password
	var dbPassword []byte
	query := "SELECT Password FROM users WHERE Username=?"

	err := db.QueryRow(query, username).Scan(
		&dbPassword,
	)
	if err != nil {
		return err
	}

	//compare input oldPassword with dbPassword
	err = bcrypt.CompareHashAndPassword(dbPassword, []byte(oldPassword))
	if err != nil {
		return err
	} else {
		bNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.MinCost)
		if err != nil {
			return err
		}
		statement := "UPDATE users SET Password=? WHERE Username=?"

		_, err = db.Exec(statement,
			bNewPassword,
			username,
		)
		if err != nil {
			return err
		}
	}
	return err
}

// Delete a user entry in database
func deleteUser(username string) error {
	_, err := db.Exec("DELETE FROM users WHERE Username=?",
		username)
	if err != nil {
		return err
	}
	return nil
}
