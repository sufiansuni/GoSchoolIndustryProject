package main

// Operations for sessions database: Insert(Create), Select(Read), Update, Delete

// Insert a new session entry into database
func insertSession(mySession session) error {
	_, err := db.Exec("INSERT INTO sessions(UUID, Username) VALUES(?, ?)",
		mySession.UUID, mySession.Username)
	if err != nil {
		return err
	}
	return nil
}

// Select/Read a session entry from database with a UUID input
func selectSession(UUID string) (session, error) {
	var mySession session
	query := "SELECT * FROM sessions WHERE UUID=?"

	err := db.QueryRow(query, UUID).Scan(
		&mySession.UUID,
		&mySession.Username,
	)
	return mySession, err
}

// Update a session entry in database

func updateSession(mySession session) error {
	statement := "UPDATE sessions SET Username=? " +
		"WHERE UUID=?"

	_, err := db.Exec(statement,
		mySession.Username,
		mySession.UUID,
	)
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
