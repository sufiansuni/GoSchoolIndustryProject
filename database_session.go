package main

// Operations for Session database: Insert, Delete, Update

// Insert a new session entry into database
func insertSession(mySession session) error {
	_, err := db.Exec("INSERT INTO sessions(UUID, Username) VALUES(?, ?)",
		mySession.UUID, mySession.Username)
	if err != nil {
		return err
	}
	return nil
}

// Delete a session entry in database
func deleteSession(UUID string) error {
	_, err := db.Exec("DELETE FROM sessions where UUID=?",
		UUID)
	if err != nil {
		return err
	}
	return nil
}
