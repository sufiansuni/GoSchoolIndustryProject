package main

import (
	"GoIndustryProject/controllers"
	"GoIndustryProject/database"
)

func main() {
	//Begin connection to database, defer close
	database.Connect()
	defer database.DB.Close()

	//Ping to test connection
	database.Ping()

	//Initialise tables and admin account
	database.Init()

	//Start HTTP Server
	controllers.StartHTTPServer()
}
