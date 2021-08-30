package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var DATABASE_USER string = "root"
var DATABASE_PASSWORD string = "password"
var DATABASE_HOST string = "localhost"
var DATABASE_PORT string = "3306"
var DATABASE_NAME string = "my_db"

//connects to the database according to set values
func connectDatabase() {
	user := DATABASE_USER
	password := DATABASE_PASSWORD
	host := DATABASE_HOST
	port:= DATABASE_PORT
	dbname := DATABASE_NAME
	connectionstring := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)

	//additional parameters
	connectionstring += "?parseTime=True&loc=Local"
	var err error
	db, err = sql.Open("mysql", connectionstring)

	//if there is an error opening the connection, handle it
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
		"Last VARCHAR(255), " +
		"Gender VARCHAR(6), " +
		"Birthday DATE, " +
		"Height SMALLINT, " +
		"Weight SMALLINT, " +
		"CaloriesPerDay FLOAT, " +
		"Halal BOOL, " +
		"Vegan BOOL, " +
		"Address VARCHAR(255), " +
		"PostalCode MEDIUMINT, " +
		"Lat FLOAT, " +
		"Lng FLOAT" +
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