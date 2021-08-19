package main

func main(){
	//Begin connection to database, then ping to test
	connectDatabase()
	defer db.Close()
	pingDatabase()

	//create tables in database
	createUserTable()
	createSessionTable()

	//Start HTTP Server
	StartHTTPServer()
}
