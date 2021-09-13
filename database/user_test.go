package database

import (
	"GoIndustryProject/models"
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestInsertUser(t *testing.T) {
	db, mock := NewMock()

	var myUser models.User
	myUser.FillDefaults()

	query := "INSERT INTO users \\(Username, Password, First, Last, Gender, Birthday, Height, Weight, ActivityLevel, CaloriesPerDay, Halal, Vegan, Address, PostalCode, Lat, Lng\\)" +
	" VALUES \\(\\?,\\?,\\?,\\?,\\?,\\?,\\?,\\?,\\?,\\?,\\?,\\?,\\?,\\?,\\?,\\?\\)"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(
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
	).WillReturnResult(sqlmock.NewResult(0, 1))

	err := InsertUser(db, myUser)
	assert.NoError(t, err)
}
