package main

func main() {
	//Begin connection to database, then ping to test
	connectDatabase()
	defer db.Close()
	pingDatabase()

	//create tables in database
	createUserTable()
	createSessionTable()
	createRestaurantTable()
	createFoodTable()
	createOrderTable()
	createOrderItemTable()

	//Start HTTP Server
	StartHTTPServer()
}
