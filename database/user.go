package database

import (
	"GoIndustryProject/models"
	"context"
	"database/sql"
	"errors"
	"time"
)

// Operations for users database: Insert(Create), Select(Read), Update, Delete

// Insert a new user entry into database
func InsertUser(db *sql.DB, myUser models.User) (err error) {

	myUser.FillDefaults()
	myUser.TitleCaseNames()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "INSERT INTO users (Username, Password, First, Last, Gender, Birthday, Height, Weight, ActivityLevel, CaloriesPerDay, Halal, Vegan, Address, PostalCode, Lat, Lng) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
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

// Select/Read a user entry from database with a username input
func SelectUserByUsername(db *sql.DB, username string) (myUser models.User, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT * FROM users WHERE Username=?"
	err = db.QueryRowContext(ctx, query, username).Scan(
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
	return
}

// Update a user entry in database
// Does not update username and password
func UpdateUserProfile(db *sql.DB, myUser models.User) (err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE users SET First=?, Last=?, Gender=?, Birthday=?, " +
		"Height=?, Weight=?, ActivityLevel=?, CaloriesPerday=?, Halal=?, Vegan=?, Address=?, PostalCode=?, Lat=?, Lng=? " +
		"WHERE Username=?"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
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

// Update a users password entry
func UpdateUserPassword(db *sql.DB, username string, newPassword []byte) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE users SET Password=? WHERE Username=?"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		newPassword,
		username,
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

// Delete a user entry in database
func DeleteUser(db *sql.DB, username string) (err error) {
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "DELETE FROM users WHERE Username = ?"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, username)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return
	}

	if rowsAffected == 0 {
		err = errors.New("no rows updated")
	}
	return
}
