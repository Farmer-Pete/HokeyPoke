package main

import (
	"os"

	"github.com/Farmer-Pete/HokeyPoke/db"
	"github.com/Farmer-Pete/HokeyPoke/db/query/table"
	"github.com/Farmer-Pete/HokeyPoke/router"
	"github.com/Farmer-Pete/HokeyPoke/server"
	"github.com/Farmer-Pete/HokeyPoke/util"
	. "github.com/go-jet/jet/v2/sqlite"
)

func main() {
	address := os.Getenv("HOKEY_LISTEN_ADDRESS")
	if len(address) == 0 {
		address = "127.0.0.1:3000"
	}

	client := db.Connect(os.Getenv("HOKEY_DB_FILE"))
	defer client.Close()

	db.SetConnection(client)

	type CountResponse struct {
		Count int
	}

	stmt := SELECT(COUNT(table.Card.ID).AS("count_response.count")).FROM(table.Card)
	response := CountResponse{}
	err := stmt.Query(client, &response)
	util.AssertNil(err)

	if response.Count == 0 {
		println("No cards found, building DB...")
		db.BuildDB(client)
	} else {
		println("DB already built")
	}

	server := server.NewServer(address)
	router.RegisterHandlers(server)
	util.AssertNil(server.Start())
}
