package main

import (
	"bible-users-api/app"
)

func main() {
	server := app.NewApp()
	server.InitializeApp()
	server.InitializeDatabase()
	server.SetupRoutes()
	server.Start()
}
