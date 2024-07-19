package router

import (
	"github.com/Farmer-Pete/HokeyPoke/server"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var sessionStore *session.Store

func RegisterHandlers(server server.Server) {
	server.App.Get("/", home)
}
