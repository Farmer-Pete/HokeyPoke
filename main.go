package main

import (
	"os"

	"github.com/Farmer-Pete/HokeyPoke/pokeclient"
	"github.com/Farmer-Pete/HokeyPoke/router"
	"github.com/Farmer-Pete/HokeyPoke/server"
	"github.com/Farmer-Pete/HokeyPoke/util"
)

func main() {
	address := os.Getenv("HOKEY_LISTEN_ADDRESS")
	if len(address) == 0 {
		address = "127.0.0.1:3000"
	}

	pokeclient.BuildDB()

	server := server.NewServer(address)
	router.RegisterHandlers(server)
	util.AssertNill(server.Start())
}
