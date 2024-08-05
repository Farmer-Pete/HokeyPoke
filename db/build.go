package db

import (
	"context"
	"database/sql"
)

func BuildDB(client *sql.DB) {
	ctx := context.Background()
	test_card_parsing()
	build_cards(client, ctx)
	build_collections(client, ctx)
}
