package main

import (
	"os"
)

func env(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	a := App{}

	clientId := env("EXT_CLIENT_ID", "")
	clientSecret := env("EXT_CLIENT_SECRET", "")
	ownerId := env("EXT_OWNER_ID", "")
	dbHost := env("DB_HOST", "localhost")
	dbPort := env("DB_PORT", "5432")
	dbName := env("DB_NAME", "annotations")
	dbUsername := env("DB_USERNAME", "postgres")
	dbPassword := env("DB_PASSWORD", "")

	a.Initialize(clientId, clientSecret, ownerId, dbHost, dbPort, dbName, dbUsername, dbPassword)
	a.Run()
}
