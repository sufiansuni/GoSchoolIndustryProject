package database

import (
	"GoIndustryProject/models"
)

// Operations for sessions database: Insert(Create), Select(Read), Update, Delete

// Insert a new session entry into database
func InsertSession(mySession models.Session) error {
	_, err := DB.Exec("INSERT INTO sessions(UUID, Username) VALUES(?, ?)",
		mySession.UUID, mySession.Username)
	if err != nil {
		return err
	}
	return nil
}

// Select/Read a session entry from database with a UUID input
func SelectSession(UUID string) (models.Session, error) {
	var mySession models.Session
	query := "SELECT * FROM sessions WHERE UUID=?"

	err := DB.QueryRow(query, UUID).Scan(
		&mySession.UUID,
		&mySession.Username,
	)
	return mySession, err
}

// Update a session entry in database

func UpdateSession(mySession models.Session) error {
	statement := "UPDATE sessions SET Username=? " +
		"WHERE UUID=?"

	_, err := DB.Exec(statement,
		mySession.Username,
		mySession.UUID,
	)
	if err != nil {
		return err
	}
	return nil
}

// Delete a session entry in database
func DeleteSession(UUID string) error {
	_, err := DB.Exec("DELETE FROM sessions WHERE UUID=?",
		UUID)
	if err != nil {
		return err
	}
	return nil
}
