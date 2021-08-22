package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

//connects to the database according to set values
func connectDatabase() {
	user := goDotEnvVariable("DATABASE_USER")
	password := goDotEnvVariable("DATABASE_PASSWORD")
	host := goDotEnvVariable("DATABASE_HOST")
	port:= goDotEnvVariable("DATABASE_PORT")
	dbname := goDotEnvVariable("DATABASE_NAME")
	connectionstring := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)

	connectionstring += "?parseTime=True&loc=Local" //additional parameters
	var err error
	db, err = sql.Open("mysql", connectionstring)

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Database Opened")
	}
}

//ping the database to test connection
func pingDatabase() {
	err := db.Ping()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Database Ping Successful")
	}
}

//creates "users" table
func createUserTable() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " +
		"users" +
		" (" +
		"Username VARCHAR(255) PRIMARY KEY, " +
		"Password BLOB, " +
		"First VARCHAR(255), " +
		"Last VARCHAR(255)" +
		")")

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Table Checked/Created: users")
	}
}

//creates "sessions" table
func createSessionTable() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " +
		"sessions" +
		" (" +
		"UUID VARCHAR(255) PRIMARY KEY, " +
		"Username VARCHAR(255)" +
		")")

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Table Checked/Created: sessions")
	}
}