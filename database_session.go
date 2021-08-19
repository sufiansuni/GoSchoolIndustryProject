package main

// Operations for Session database: Insert, Delete, Update

func insertSession(mySession session) error {
	_, err := db.Exec("INSERT INTO sessions(UUID, Username) VALUES(?, ?)",
		mySession.UUID, mySession.Username)
	if err != nil {
		return err
	}
	return nil
}

func deleteSession(UUID string) error {
	_, err := db.Exec("DELETE FROM sessions where UUID=?",
		UUID)
	if err != nil {
		return err
	}
	return nil
}
