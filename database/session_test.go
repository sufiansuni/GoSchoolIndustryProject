package database

import (
	"GoIndustryProject/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestInsertSession(t *testing.T) {
	db, mock := NewMock()

	var mySession models.Session

	query := "INSERT INTO sessions \\(UUID, Username\\) VALUES \\(\\?, \\?\\)"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(
		mySession.UUID,
		mySession.Username,

	).WillReturnResult(sqlmock.NewResult(0, 1))

	err := InsertSession(db, mySession)
	assert.NoError(t, err)
}

func TestInsertSessionError(t *testing.T) {
	db, mock := NewMock()

	var mySession models.Session

	query := "INSERT INTO sessions \\(UUID, Username\\) VALUES \\(\\?, \\?\\)"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(
		mySession.UUID,
		mySession.Username,

	).WillReturnResult(sqlmock.NewResult(0, 0))

	err := InsertSession(db, mySession)
	assert.Error(t, err)
}

func TestSelectSession(t *testing.T) {
	db, mock := NewMock()

	var mySession models.Session
	mySession.UUID = "test-uuid"
	mySession.Username = "test-user"

	query := "SELECT \\* FROM sessions WHERE UUID=\\?"

	rows := sqlmock.NewRows([]string{
		"UUID",
		"Username",
	}).
		AddRow(
			mySession.UUID,
			mySession.Username,
		)

	mock.ExpectQuery(query).WithArgs(mySession.UUID).WillReturnRows(rows)

	resultSession, err := SelectSession(db, mySession.UUID)
	assert.NotNil(t, resultSession)
	assert.NoError(t, err)
}

func TestSelectSessionError(t *testing.T) {
	db, mock := NewMock()

	var mySession models.Session
	mySession.UUID = "test-uuid"
	mySession.Username = "test-user"

	query := "SELECT \\* FROM sessions WHERE UUID=\\?"

	rows := sqlmock.NewRows([]string{
		"UUID",
		"Username",
	})

	mock.ExpectQuery(query).WithArgs(mySession.UUID).WillReturnRows(rows)

	resultSession, err := SelectSession(db, mySession.UUID)
	assert.Zero(t, resultSession)
	assert.Error(t, err)
}

func TestUpdateSession(t *testing.T) {
	db, mock := NewMock()

	var mySession models.Session
	mySession.UUID = "test-uuid"
	mySession.Username = "test-user"

	query := "UPDATE sessions SET Username=\\? WHERE UUID=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(
		mySession.Username,
		mySession.UUID,
		
	).WillReturnResult(sqlmock.NewResult(0, 1))

	err := UpdateSession(db, mySession)
	assert.NoError(t, err)
}

func TestUpdateSessionError(t *testing.T) {
	db, mock := NewMock()

	var mySession models.Session
	mySession.UUID = "test-uuid"
	mySession.Username = "test-user"

	query := "UPDATE sessions SET Username=\\? WHERE UUID=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(
		mySession.Username,
		mySession.UUID,
		
	).WillReturnResult(sqlmock.NewResult(0, 0))

	err := UpdateSession(db, mySession)
	assert.Error(t, err)
}

func TestDeleteSession(t *testing.T) {
	db, mock := NewMock()

	query := "DELETE FROM sessions WHERE UUID=\\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs("uuid").WillReturnResult(sqlmock.NewResult(0, 1))

	err := DeleteSession(db, "uuid")
	assert.NoError(t, err)
}

func TestDeleteSessionError(t *testing.T) {
	db, mock := NewMock()

	query := "DELETE FROM sessions WHERE UUID=\\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs("uuid").WillReturnResult(sqlmock.NewResult(0, 0))

	err := DeleteSession(db, "uuid")
	assert.Error(t, err)
}
