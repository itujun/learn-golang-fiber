package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		IdleTimeout: time.Second * 5, // Set idle timeout to 5 seconds
		ReadTimeout: time.Second * 5, // Set read timeout to 5 seconds
		WriteTimeout: time.Second * 5, // Set write timeout to 5 seconds
	})

	err := app.Listen("localhost:3000")
	if err != nil {
		panic(err)
	}
}