package database

import (
	"GoIndustryProject/models"
	"context"
	"database/sql"
	"errors"
	"time"
)

// Operations for sessions database: Insert(Create), Select(Read), Update, Delete

// Insert a new session entry into database
func InsertSession(db *sql.DB, mySession models.Session) (err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "INSERT INTO sessions (UUID, Username) VALUES (?, ?)"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		mySession.UUID,
		mySession.Username,
	)
	if err != nil {
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return
	}

	if rowsAffected == 0 {
		err = errors.New("no rows updated")
	}
	return
}

// Select/Read a session entry from database with a UUID input
func SelectSession(db *sql.DB, UUID string) (mySession models.Session, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT * FROM sessions WHERE UUID=?"
	err = db.QueryRowContext(ctx, query, UUID).Scan(
		&mySession.UUID,
		&mySession.Username,
	)
	return
}

// Update a session entry in database
func UpdateSession(db *sql.DB, mySession models.Session) (err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE sessions SET Username=? " +
		"WHERE UUID=?"

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()
	
	result, err := stmt.ExecContext(ctx,
		mySession.Username,
		mySession.UUID,
	)
	if err != nil {
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return
	}

	if rowsAffected == 0 {
		err = errors.New("no rows updated")
	}
	return
}

// Delete a session entry in database
func DeleteSession(db *sql.DB, uuid string) (err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "DELETE FROM sessions WHERE UUID=?"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, uuid)
	if err != nil {
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return
	}

	if rowsAffected == 0 {
		err = errors.New("no rows updated")
	}
	return
}
