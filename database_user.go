package main

// Operations for User database: Insert, Delete, Update

func insertUser(myUser user) error {
	_, err := db.Exec("INSERT INTO users (Username, Password, First, Last) VALUES (?,?,?,?)",
		myUser.Username, myUser.Password, myUser.First, myUser.Last)
	if err != nil {
		return err
	}
	return nil
}
