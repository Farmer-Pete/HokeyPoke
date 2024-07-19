package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
)

func NewServer(address string) Server {
	engine := html.New("./templates", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	return Server{
		App:          app,
		SessionStore: session.New(),
		address:      address}
}

func (s Server) Start() error {
	return s.App.Listen(s.address)
}
