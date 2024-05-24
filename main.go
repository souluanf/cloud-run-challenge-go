package main

import (
	"cloud-run-challenge-go/api"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Get("/:cep", api.HandleRequest)
	log.Fatal(app.Listen(":8080"))
}
