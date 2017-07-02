package main

import "os"

func main() {
	app := App{}
	app.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
		os.Getenv("APP_DB_HOST"),
		"5432",
		"require")

	port := os.Getenv("PORT")
	if port == nil {
		port = 5678
	}
	app.Run(":" + port)
}
