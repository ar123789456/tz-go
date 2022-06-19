package main

import (
	"log"
	"tz/server"
)

func main() {
	app := server.NewApp()
	if err := app.Run("8080"); err != nil {
		log.Fatalln(err)
	}
}
