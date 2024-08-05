package db

import (
	"context"
	"database/sql"
	"encoding/csv"
	"os"
	"strconv"

	"github.com/Farmer-Pete/HokeyPoke/db/query/model"
	"github.com/Farmer-Pete/HokeyPoke/db/query/table"
	"github.com/Farmer-Pete/HokeyPoke/util"
	. "github.com/go-jet/jet/v2/sqlite"
)

func build_collections(client *sql.DB, ctx context.Context) {
	tx, err := client.BeginTx(ctx, nil)
	defer tx.Rollback()
	util.AssertNil(err)

	f, err := os.OpenFile("data/collection.csv", os.O_RDONLY, 0644)
	util.AssertNil(err)
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Read() // Skip header
	rawRecords, err := csvReader.ReadAll()
	util.AssertNil(err)

	stmt := table.Collection.INSERT(table.Collection.Name).
		VALUES("Peter").
		RETURNING(table.Collection.AllColumns)

	var collection model.Collection
	stmt.QueryContext(ctx, client, &collection)

	for _, record := range rawRecords {
		count, err := strconv.ParseInt(record[1], 10, 32)
		util.AssertNil(err)

		stmt := table.CardCollection.INSERT(
			table.CardCollection.CardID,
			table.CardCollection.CollectionID,
			table.CardCollection.Count).
			VALUES(
				SELECT(table.Card.ID).
					FROM(table.Card).
					WHERE(table.Card.PtcgID.EQ(String(record[0]))),
				collection.ID, count,
			)

		stmt.ExecContext(ctx, client)
	}

	util.AssertNil(tx.Commit())
}
