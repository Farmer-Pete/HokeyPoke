package dal

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Farmer-Pete/HokeyPoke/dal/db"
	"github.com/Farmer-Pete/HokeyPoke/dal/db/group"
	"github.com/Farmer-Pete/HokeyPoke/dal/db/supertype"
	"github.com/Farmer-Pete/HokeyPoke/dal/db/typ3"
	"github.com/Farmer-Pete/HokeyPoke/data"
	"github.com/Farmer-Pete/HokeyPoke/util"
	_ "github.com/mattn/go-sqlite3"
)

func BuildDB(ctx context.Context, client *db.Client) {
	tx, err := client.Tx(ctx)
	util.AssertNill(err)

	files, err := filepath.Glob("data/pokemon-tcg-data/cards/en/*.json")
	util.AssertNill(err)

	for _, file := range files {
		fmt.Printf("Processing %s ", file)

		bytes, err := os.ReadFile(file)
		util.AssertNill(err)

		var cards []data.CardMetadata

		util.AssertNill(json.Unmarshal(bytes, &cards))

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

			Supertype := GetSupertype(ctx, client, card)
			Types := GetTypes(ctx, client, card)
			Group := GetGroup(ctx, client, card)

			err := client.Card.
				Create().
				SetMetadata(card).
				SetName(card.Name).
				SetPtcgID(card.ID).
				SetSupertype(Supertype).
				AddTypes(Types...).
				SetGroup(Group).
				OnConflict().
				DoNothing().
				Exec(ctx)

			if err != nil && err != sql.ErrNoRows {
				// It's okay if nothing was inserted (because it's already there)
				panic(err)
			}

		}

		fmt.Println()
	}
	util.AssertNill(tx.Commit())
}

func GetSupertype(ctx context.Context, client *db.Client, card data.CardMetadata) *db.Supertype {
	s, err := client.Supertype.
		Query().
		Where(supertype.Name(card.Supertype)).
		First(ctx)

	if err != nil {
		if db.IsNotFound(err) {
			s, err = client.Supertype.
				Create().
				SetName(card.Supertype).
				Save(ctx)
			util.AssertNill(err)
		} else {
			util.AssertNill(err)
		}
	}

	return s
}

func GetTypes(ctx context.Context, client *db.Client, card data.CardMetadata) []*db.Typ3 {
	result := []*db.Typ3{}

	for _, name := range card.Types {
		t, err := client.Typ3.
			Query().
			Where(typ3.Name(name)).
			First(ctx)

		if err != nil {
			if db.IsNotFound(err) {
				t, err = client.Typ3.
					Create().
					SetName(name).
					Save(ctx)
				util.AssertNill(err)
			} else {
				util.AssertNill(err)
			}
		}

		result = append(result, t)
	}

	return result
}

func GetGroup(ctx context.Context, client *db.Client, card data.CardMetadata) *db.Group {
	name := data.NormalizeName(card)

	g, err := client.Group.
		Query().
		Where(group.Name(name)).
		WithSupertype().
		First(ctx)

	if err != nil {
		if db.IsNotFound(err) {
			g, err = client.Group.
				Create().
				SetName(name).
				SetSupertype(GetSupertype(ctx, client, card)).
				Save(ctx)
			util.AssertNill(err)
		} else {
			panic(err)
		}
	}

	g.Update().AddTypes(GetTypes(ctx, client, card)...).Save(ctx)

	return g
}
