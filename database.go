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

// Connects to the database according to set values
func connectDatabase() {
	user := DATABASE_USER
	password := DATABASE_PASSWORD
	host := DATABASE_HOST
	port := DATABASE_PORT
	dbname := DATABASE_NAME
	connectionstring := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)

	// additional parameters
	connectionstring += "?parseTime=True&loc=Local"
	var err error
	db, err = sql.Open("mysql", connectionstring)

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Database Opened")
	}
}

// Ping the database to test connection
func pingDatabase() {
	err := db.Ping()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Database Ping Successful")
	}
}

// Creates "users" table in database
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
		"Height SMALLINT UNSIGNED, " +
		"Weight SMALLINT UNSIGNED, " +
		"CaloriesPerDay FLOAT UNSIGNED, " +
		"Halal BOOL, " +
		"Vegan BOOL, " +
		"Address VARCHAR(255), " +
		"PostalCode MEDIUMINT UNSIGNED, " +
		"Lat FLOAT, " +
		"Lng FLOAT" +
		")")

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Table Checked/Created: users")
	}
}

// Creates "sessions" table in database
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

// Creates "restaurants" table in database
func createRestaurantTable() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " +
		"restaurants" +
		" (" +
		"ID MEDIUMINT UNSIGNED PRIMARY KEY AUTO_INCREMENT, " +
		"Name VARCHAR(255), " +
		"Description VARCHAR(255), " +
		"Halal BOOL, " +
		"Vegan BOOL, " +
		"Address VARCHAR(255), " +
		"PostalCode MEDIUMINT UNSIGNED, " +
		"Lat FLOAT, " +
		"Lng FLOAT" +
		")")

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Table Checked/Created: restaurants")
	}
}

// Creates "foods" table in database
func createFoodTable() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " +
		"foods" +
		" (" +
		"ID MEDIUMINT UNSIGNED PRIMARY KEY AUTO_INCREMENT, " +
		"RestaurantID VARCHAR(255), " +
		"Name VARCHAR(255), " +
		"Price FLOAT UNSIGNED, " +
		"Calories FLOAT UNSIGNED, " +
		"Halal BOOL, " +
		"Vegan BOOL " +
		")")

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Table Checked/Created: foods")
	}
}
