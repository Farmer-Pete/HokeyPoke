package db

import (
	"database/sql"
	"fmt"

	"github.com/Farmer-Pete/HokeyPoke/util"
)

func Connect(file string) *sql.DB {
	client, err := sql.Open(
		"sqlite3",
		fmt.Sprintf("file:%s?mode=rwc&cache=shared&_fk=1", file),
	)
	util.AssertNil(err)

	return client
}

var _DB_CLIENT *sql.DB = nil

func GetConnection() *sql.DB {
	util.AssertNotNil(_DB_CLIENT, "No DB client found")
	return _DB_CLIENT
}

func SetConnection(client *sql.DB) {
	util.AssertNotNil(client, "Client cannot be nil")
	_DB_CLIENT = client
}
