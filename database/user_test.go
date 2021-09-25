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

	query := "INSERT INTO users \\(Username, Password, First, Last, Gender, Birthday, Height, Weight, ActivityLevel, CaloriesPerDay, Halal, Vegan, Address, Unit, Lat, Lng\\)" +
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
		myUser.Unit,
		myUser.Lat,
		myUser.Lng,
	).WillReturnResult(sqlmock.NewResult(0, 1))

	err := InsertUser(db, myUser)
	assert.NoError(t, err)
}

func TestInserUserError(t *testing.T) {
	db, mock := NewMock()

	var myUser models.User
	myUser.FillDefaults()

	query := "INSERT INTO users \\(Username, Password, First, Last, Gender, Birthday, Height, Weight, ActivityLevel, CaloriesPerDay, Halal, Vegan, Address, Unit, Lat, Lng\\)" +
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
		myUser.Unit,
		myUser.Lat,
		myUser.Lng,
	).WillReturnResult(sqlmock.NewResult(0, 0))

	err := InsertUser(db, myUser)
	assert.Error(t, err)
}

func TestSelectUserbyUsername(t *testing.T) {
	db, mock := NewMock()

	var myUser models.User
	myUser.FillDefaults()
	myUser.Username = "test-user"

	query := "SELECT \\* FROM users WHERE Username=\\?"

	rows := sqlmock.NewRows([]string{
		"Username",
		"Password",
		"First",
		"Last",
		"Gender",
		"Birthday",
		"Height",
		"Weight",
		"ActivityLevel",
		"CaloriesPerDay",
		"Halal",
		"Vegan",
		"Address",
		"Unit",
		"Lat",
		"Lng",
	}).
		AddRow(
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
			myUser.Unit,
			myUser.Lat,
			myUser.Lng,
		)

	mock.ExpectQuery(query).WithArgs(myUser.Username).WillReturnRows(rows)

	resultUser, err := SelectUserByUsername(db, myUser.Username)
	assert.NotNil(t, resultUser)
	assert.NoError(t, err)
}

func TestSelectUserbyUsernameError(t *testing.T) {
	db, mock := NewMock()

	var myUser models.User
	myUser.FillDefaults()
	myUser.Username = "test-user"

	query := "SELECT \\* FROM users WHERE Username=\\?"

	rows := sqlmock.NewRows([]string{
		"Username",
		"Password",
		"First",
		"Last",
		"Gender",
		"Birthday",
		"Height",
		"Weight",
		"ActivityLevel",
		"CaloriesPerDay",
		"Halal",
		"Vegan",
		"Address",
		"PostalCode",
		"Lat",
		"Lng",
	})

	mock.ExpectQuery(query).WithArgs(myUser.Username).WillReturnRows(rows)

	resultUser, err := SelectUserByUsername(db, myUser.Username)
	assert.Zero(t, resultUser)
	assert.Error(t, err)
}

func TestUpdateUserProfile(t *testing.T) {
	db, mock := NewMock()

	var myUser models.User
	myUser.FillDefaults()

	query := "UPDATE users SET First=\\?, Last=\\?, Gender=\\?, Birthday=\\?, " +
		"Height=\\?, Weight=\\?, ActivityLevel=\\?, CaloriesPerday=\\?, Halal=\\?, Vegan=\\?, " +
		"Address=\\?, Unit=\\?,  Lat=\\?, Lng=\\? " +
		"WHERE Username=\\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(
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
		myUser.Unit,
		myUser.Lat,
		myUser.Lng,
		myUser.Username,
	).WillReturnResult(sqlmock.NewResult(0, 1))

	err := UpdateUserProfile(db, myUser)
	assert.NoError(t, err)
}

func TestUpdateUserProfileError(t *testing.T) {
	db, mock := NewMock()

	var myUser models.User
	myUser.FillDefaults()

	query := "UPDATE users SET First=\\?, Last=\\?, Gender=\\?, Birthday=\\?, " +
		"Height=\\?, Weight=\\?, ActivityLevel=\\?, CaloriesPerday=\\?, Halal=\\?, Vegan=\\?, " +
		"Address=\\?, PostalCode=\\?, Lat=\\?, Lng=\\? " +
		"WHERE Username=\\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(
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
		myUser.Unit,
		myUser.Lat,
		myUser.Lng,
		myUser.Username,
	).WillReturnResult(sqlmock.NewResult(0, 0))

	err := UpdateUserProfile(db, myUser)
	assert.Error(t, err)
}

func TestUpdateUserPassword(t *testing.T) {
	db, mock := NewMock()

	var myUser models.User
	myUser.FillDefaults()
	myUser.Username = "test-user"
	myUser.Password = []byte{98, 99}

	query := "UPDATE users SET Password=\\? WHERE Username=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(
		myUser.Password,
		myUser.Username,
	).WillReturnResult(sqlmock.NewResult(0, 1))

	err := UpdateUserPassword(db, myUser.Username, myUser.Password)
	assert.NoError(t, err)
}

func TestUpdateUserPasswordError(t *testing.T) {
	db, mock := NewMock()

	var myUser models.User
	myUser.FillDefaults()
	myUser.Username = "test-user"
	myUser.Password = []byte{98, 99}

	query2 := "UPDATE users SET Password=\\? WHERE Username=\\?"
	prep := mock.ExpectPrepare(query2)
	prep.ExpectExec().WithArgs(
		myUser.Password,
		myUser.Username,
	).WillReturnResult(sqlmock.NewResult(0, 0))

	err := UpdateUserPassword(db, myUser.Username, myUser.Password)
	assert.Error(t, err)
}

func TestDeleteUser(t *testing.T) {
	db, mock := NewMock()

	query := "DELETE FROM users WHERE Username = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs("name").WillReturnResult(sqlmock.NewResult(0, 1))

	err := DeleteUser(db, "name")
	assert.NoError(t, err)
}

func TestDeleteUserError(t *testing.T) {
	db, mock := NewMock()

	query := "DELETE FROM users WHERE Username = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs("name").WillReturnResult(sqlmock.NewResult(0, 0))

	err := DeleteUser(db, "name")
	assert.Error(t, err)
}
