package main

import (
	"context"
	"os"

	"github.com/Farmer-Pete/HokeyPoke/dal"
	"github.com/Farmer-Pete/HokeyPoke/router"
	"github.com/Farmer-Pete/HokeyPoke/server"
	"github.com/Farmer-Pete/HokeyPoke/util"
)

func main() {
	address := os.Getenv("HOKEY_LISTEN_ADDRESS")
	if len(address) == 0 {
		address = "127.0.0.1:3000"
	}

	client := dal.Connect(os.Getenv("HOKEY_DB_FILE"))
	defer client.Close()

	ctx := context.Background()
	if client.Card.Query().CountX(ctx) == 0 {
		println("No cards found, building DB...")
		dal.BuildDB(ctx, client)
	} else {
		println("DB already built")
	}

	dal.SetConnection(client, &ctx)

	server := server.NewServer(address)
	router.RegisterHandlers(server)
	util.AssertNill(server.Start())
}
