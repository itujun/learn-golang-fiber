package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		IdleTimeout: time.Second * 5, // Set idle timeout to 5 seconds
		ReadTimeout: time.Second * 5, // Set read timeout to 5 seconds
		WriteTimeout: time.Second * 5, // Set write timeout to 5 seconds
		Prefork: true, // Enable pre-forking
	})

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, World!")
	})

	if fiber.IsChild(){
		fmt.Println("I'am child process");
	}else{
		fmt.Println("I'am parent process");
	}
	// jalankan server, lalu cek task manager
	// maka akan menemunkan total 5 process dari 1 proses parent dan 4 proses child

	err := app.Listen("localhost:3000")
	if err != nil {
		panic(err)
	}
}