package main

// Operations for Session database: Insert, Delete, Update

//insert a new session entry
func insertSession(mySession session) error {
	_, err := db.Exec("INSERT INTO sessions(UUID, Username) VALUES(?, ?)",
		mySession.UUID, mySession.Username)
	if err != nil {
		return err
	}
	return nil
}

//delete a session entry
func deleteSession(UUID string) error {
	_, err := db.Exec("DELETE FROM sessions where UUID=?",
		UUID)
	if err != nil {
		return err
	}
	return nil
}
