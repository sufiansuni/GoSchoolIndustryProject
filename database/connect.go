package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB
var DATABASE_USER string = "root"
var DATABASE_PASSWORD string = "password"
var DATABASE_HOST string = "localhost"
var DATABASE_PORT string = "3306"
var DATABASE_NAME string = "my_db"

// Connects to the database according to set values
func Connect() {
	user := DATABASE_USER
	password := DATABASE_PASSWORD
	host := DATABASE_HOST
	port := DATABASE_PORT
	dbname := DATABASE_NAME
	connectionstring := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)

	// additional parameters
	connectionstring += "?parseTime=True&loc=Local"
	var err error
	DB, err = sql.Open("mysql", connectionstring)

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Database Opened")
	}
}

// Ping the database to test connection
func Ping() {
	err := DB.Ping()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Database Ping Successful")
	}
}