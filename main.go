package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	MountController(app)

	log.Println("aitime listening")

	app.Listen(":3000")
}
