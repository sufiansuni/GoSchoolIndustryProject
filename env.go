package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load()

	if err != nil {
		fmt.Println(err)
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}