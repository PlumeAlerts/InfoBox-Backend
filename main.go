package main

import (
	"fmt"
	"os"
)

func env(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	fmt.Printf("Enviroment variable %s not set, defaulting to %s", key, defaultValue)
	return defaultValue
}

func main() {
	a := App{}
	fmt.Println("Starting")

	clientId := env("EXT_CLIENT_ID", "")
	clientSecret := env("EXT_CLIENT_SECRET", "")
	ownerId := env("EXT_OWNER_ID", "")
	dbHost := env("DB_HOST", "localhost")
	dbPort := env("DB_PORT", "5432")
	dbName := env("DB_NAME", "annotations")
	dbUsername := env("DB_USERNAME", "postgres")
	dbPassword := env("DB_PASSWORD", "")
	fmt.Printf("Connecting to %s using port %s", dbHost, dbPort)

	a.Initialize(clientId, clientSecret, ownerId, dbHost, dbPort, dbName, dbUsername, dbPassword)
	fmt.Println("Started")
	a.Run()
}
