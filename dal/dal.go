package dal

import (
	"context"
	"fmt"

	"github.com/Farmer-Pete/HokeyPoke/dal/db"
	"github.com/Farmer-Pete/HokeyPoke/util"
)

func Connect(file string) *db.Client {
	client, err := db.Open(
		"sqlite3",
		fmt.Sprintf("file:%s?mode=rwc&cache=shared&_fk=1", file),
	)
	util.AssertNill(err)

	// Run auto migrations
	err = client.Schema.Create(context.Background())
	util.AssertNill(err)

	return client
}

var _DB_CLIENT *db.Client = nil
var _DB_CTX *context.Context = nil

func GetConnection() (*db.Client, *context.Context) {
	util.AssertNotNill(_DB_CLIENT, "No DB client found")
	util.AssertNotNill(_DB_CTX, "No DB context found")
	return _DB_CLIENT, _DB_CTX
}

func SetConnection(client *db.Client, ctx *context.Context) {
	util.AssertNotNill(client, "Client cannot be nil")
	util.AssertNotNill(ctx, "Context cannot be nil")
	_DB_CLIENT = client
	_DB_CTX = ctx
}
