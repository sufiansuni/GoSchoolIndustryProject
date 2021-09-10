package env

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// use godot package to load/read the .env file and
// return the value of the key
func Get(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println(err)
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
