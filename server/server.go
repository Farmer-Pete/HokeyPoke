package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type Server struct {
	App          *fiber.App
	SessionStore *session.Store
	address      string
}
