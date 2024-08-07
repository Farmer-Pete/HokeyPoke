package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/Farmer-Pete/HokeyPoke/data"
	"github.com/Farmer-Pete/HokeyPoke/db/query/model"
	"github.com/Farmer-Pete/HokeyPoke/db/query/table"
	"github.com/Farmer-Pete/HokeyPoke/util"
	mapset "github.com/deckarep/golang-set/v2"
	. "github.com/go-jet/jet/v2/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func test_card_parsing() {
	files, err := filepath.Glob("data/pokemon-tcg-data/cards/en/*.json")
	util.AssertNil(err)

	c := data.CardMetadata{}
	t := reflect.TypeOf(c)
	fields := mapset.NewSet[string]()
	for i := 0; i < t.NumField(); i++ {
		name := strings.Split(
			t.Field(i).Tag.Get("json"), ",",
		)[0]
		fields.Add(name)
	}

	cards := []map[string]interface{}{}
	for _, file := range files {

		bytes, err := os.ReadFile(file)
		util.AssertNil(err)
		util.AssertNil(json.Unmarshal(bytes, &cards))

		for _, card := range cards {
			for key := range card {
				if !fields.ContainsOne(key) {
					panic(fmt.Sprintf("File %s has unknown field %s", file, key))
				}
			}
		}
	}
}

func build_cards(client *sql.DB, ctx context.Context) {
	tx, err := client.BeginTx(ctx, nil)
	defer tx.Rollback()
	util.AssertNil(err)

	files, err := filepath.Glob("data/pokemon-tcg-data/cards/en/*.json")
	util.AssertNil(err)

	for _, file := range files {
		fmt.Printf("Processing %s ", file)

		bytes, err := os.ReadFile(file)
		util.AssertNil(err)

		var cards []data.CardMetadata

		util.AssertNil(json.Unmarshal(bytes, &cards))

		for _, card := range cards {
			switch card.Supertype {
			case "Pokémon":
				fmt.Print(".")
			case "Energy":
				fmt.Print("*")
			case "Trainer":
				fmt.Print("+")
			default:
				fmt.Print("?")
			}

			Supertype := get_supertype(ctx, client, card)
			Types := get_types(ctx, client, card)
			Group := get_group(ctx, client, card)

			cardJson, err := json.Marshal(card)
			util.AssertNil(err)
			cardJsonStr := string(cardJson)

			stmt := table.Card.INSERT(table.Card.MutableColumns).
				MODEL(model.Card{
					Metadata:    &cardJsonStr,
					Name:        card.Name,
					PtcgID:      card.ID,
					SupertypeID: *Supertype.ID,
					GroupID:     *Group.ID,
				}).
				ON_CONFLICT(table.Card.PtcgID).
				DO_UPDATE(
					SET(
						table.Card.Metadata.SET(table.Card.EXCLUDED.Metadata),
						table.Card.Name.SET(table.Card.EXCLUDED.Name),
						table.Card.SupertypeID.SET(table.Card.EXCLUDED.SupertypeID),
						table.Card.GroupID.SET(table.Card.EXCLUDED.GroupID),
					),
				).
				RETURNING(table.Card.AllColumns)

			var result model.Card
			stmt.QueryContext(ctx, client, &result)
			util.AssertNil(err, stmt.DebugSql())

			for _, t := range Types {
				stmt = table.CardType.INSERT(table.CardType.MutableColumns).
					MODEL(model.CardType{CardID: *result.ID, TypeID: *t.ID}).
					ON_CONFLICT(table.CardType.CardID, table.CardType.TypeID).
					DO_NOTHING()

				_, err = stmt.ExecContext(ctx, client)
				util.AssertNil(err, stmt.DebugSql())
			}
		}

		fmt.Println()
	}

	util.AssertNil(tx.Commit())
}

func get_supertype(ctx context.Context, client *sql.DB, card data.CardMetadata) *model.Supertype {
	stmt := table.Supertype.INSERT(table.Supertype.MutableColumns).
		MODEL(model.Supertype{Name: card.Supertype}).
		ON_CONFLICT(table.Supertype.Name).
		DO_UPDATE(
			SET(
				table.Supertype.Name.SET(table.Supertype.EXCLUDED.Name),
			),
		).
		RETURNING(table.Supertype.AllColumns)

	var result model.Supertype
	err := stmt.QueryContext(ctx, client, &result)
	util.AssertNil(err, stmt.DebugSql())
	util.AssertNotNil(result, "Supertype not created properly")

	return &result
}

func get_types(ctx context.Context, client *sql.DB, card data.CardMetadata) []model.Type {
	if len(card.Types) == 0 {
		return []model.Type{}
	}

	stmt := table.Type.INSERT(table.Type.MutableColumns)

	for _, name := range card.Types {
		stmt = stmt.MODEL(model.Type{Name: name})
	}

	stmt = stmt.
		ON_CONFLICT(table.Type.Name).
		DO_UPDATE(
			SET(
				table.Type.Name.SET(table.Type.EXCLUDED.Name),
			),
		).
		RETURNING(table.Type.AllColumns)

	var result []model.Type
	err := stmt.QueryContext(ctx, client, &result)
	util.AssertNil(err, stmt.DebugSql())
	util.AssertTrue(len(result) == len(card.Types), "Type not created properly")

	return result
}

func get_group(ctx context.Context, client *sql.DB, card data.CardMetadata) *model.Group {
	name := data.NormalizeName(card)

	supertype := get_supertype(ctx, client, card)

	stmt := table.Group.INSERT(table.Group.MutableColumns).
		MODEL(model.Group{Name: name, SupertypeID: *supertype.ID}).
		ON_CONFLICT(table.Group.Name).
		DO_UPDATE(
			SET(
				table.Group.Name.SET(table.Group.EXCLUDED.Name),
			),
		).
		RETURNING(table.Group.AllColumns)

	var result model.Group
	err := stmt.QueryContext(ctx, client, &result)
	util.AssertNil(err, stmt.DebugSql())
	util.AssertNotNil(result, "Group not created properly")

	types := get_types(ctx, client, card)

	stmt = table.GroupType.INSERT(table.GroupType.MutableColumns)

	for _, t := range types {
		stmt = stmt.MODEL(model.GroupType{GroupID: *result.ID, TypeID: *t.ID})
	}

	stmt = stmt.
		ON_CONFLICT(table.GroupType.GroupID, table.GroupType.TypeID).
		DO_NOTHING()

	stmt.ExecContext(ctx, client)
	util.AssertNil(err, stmt.DebugSql())

	return &result
}
