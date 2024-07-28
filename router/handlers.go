package router

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/Farmer-Pete/HokeyPoke/data"
	"github.com/Farmer-Pete/HokeyPoke/db"
	"github.com/Farmer-Pete/HokeyPoke/db/query/model"
	"github.com/Farmer-Pete/HokeyPoke/db/query/table"
	"github.com/Farmer-Pete/HokeyPoke/util"
	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/k0kubun/pp/v3"
	sqlite "github.com/mattn/go-sqlite3"
)

type pokemonQuery struct {
	Categories     []string
	Types          []string
	Group          string
	CollectionOnly bool
}

func (pq *pokemonQuery) load(encoded string) {
	qBytes, err := base64.StdEncoding.DecodeString(encoded)
	util.AssertNil(err)
	json.Unmarshal(qBytes, pq)
}

func (pq *pokemonQuery) copy() *pokemonQuery {
	new := pokemonQuery{}
	new.Categories = slices.Clone(pq.Categories)
	new.Types = slices.Clone(pq.Types)
	return &new
}

func (pq *pokemonQuery) Encode() string {
	result, err := json.Marshal(pq)
	util.AssertNil(err)
	return base64.StdEncoding.EncodeToString(result)
}

var COLLECTION_NAME = "Peter"

func home(ctx *fiber.Ctx) error {
	var query pokemonQuery
	qStrings := ctx.Queries()
	query.load(qStrings["q"])

	pp.Println(query)
	pp.Println(qStrings)

	action := qStrings["action"]
	card_PtcgID := qStrings["card_PtcgID"]

	db := db.GetConnection()

	/*******************************************************************************
	 * CRUD Operations
	 *******************************************************************************/

	if action == "add" && len(card_PtcgID) > 0 {
		card := CTE("CardCTE")
		collection := CTE("CollectionCTE")

		stmt := WITH(
			card.AS(
				SELECT(table.Card.ID).
					FROM(table.Card).
					WHERE(table.Card.PtcgID.EQ(String(card_PtcgID))),
			),
			collection.AS(
				SELECT(table.Collection.ID).
					FROM(table.Collection).
					WHERE(table.Collection.Name.EQ(String(COLLECTION_NAME))),
			),
		)(
			table.CardCollection.INSERT(table.CardCollection.CardID, table.CardCollection.CollectionID).
				QUERY(
					SELECT(card.AllColumns(), collection.AllColumns()).
						FROM(card, collection),
				),
		)

		result, err := stmt.Exec(db)

		if sqlErr, _ := err.(sqlite.Error); sqlErr.Code == sqlite.ErrConstraint {
			// Do nothing if the card already exists
		} else if err != nil {
			// Unknown error
			panic(fmt.Sprintf("%s: %s", stmt.DebugSql(), err))
		} else {
			// Everything looks good, but check that the card was actually added
			count, err := result.RowsAffected()
			util.AssertNil(err)
			util.AssertTrue(count == 1, fmt.Sprintf("Expected 1 row to be affected, got %d", count))
		}

		print(stmt.DebugSql())
		pp.Println(err)

	} else if action == "remove" && len(card_PtcgID) > 0 {
		stmt := table.CardCollection.DELETE().WHERE(
			table.CardCollection.CardID.IN(
				SELECT(table.Card.ID).
					FROM(table.Card).
					WHERE(table.Card.PtcgID.EQ(String(card_PtcgID))),
			),
		)

		result, err := stmt.Exec(db)
		util.AssertNil(err, stmt.DebugSql())

		count, err := result.RowsAffected()
		util.AssertNil(err)
		util.AssertTrue(count == 1, fmt.Sprintf("Expected 1 row to be affected, got %d", count))

		print(stmt.DebugSql())
	} else if action == "increment" && len(card_PtcgID) > 0 {
		stmt := table.CardCollection.UPDATE(table.CardCollection.Count).
			SET(table.CardCollection.Count.SET(table.CardCollection.Count.ADD(Int(1)))).
			WHERE(table.CardCollection.CardID.IN(
				SELECT(table.Card.ID).
					FROM(table.Card).
					WHERE(table.Card.PtcgID.EQ(String(card_PtcgID))),
			))

		result, err := stmt.Exec(db)
		util.AssertNil(err, stmt.DebugSql())

		count, err := result.RowsAffected()
		util.AssertNil(err)
		util.AssertTrue(count == 1, fmt.Sprintf("Expected 1 row to be affected, got %d", count))

		print(stmt.DebugSql())
	}

	/*******************************************************************************
	 * Fetch Data
	 *******************************************************************************/

	var strCategories []string
	util.AssertNil(
		SELECT(table.Supertype.Name).
			FROM(table.Supertype).
			ORDER_BY(table.Supertype.Name).
			Query(db, &strCategories))

	var strTypes []string
	util.AssertNil(
		SELECT(table.Type.Name).
			FROM(table.Type).
			ORDER_BY(table.Type.Name).
			Query(db, &strTypes))

	type ModelGroup struct {
		model.Group
		Types      []model.Type
		Collection []model.CardCollection
	}
	type Group struct {
		ModelGroup
		CardCollectionCount int32
	}

	filteredGroups := func() []Group {

		stmt := SELECT(
			table.Group.AllColumns,
			table.Type.AllColumns,
			table.CardCollection.AllColumns).
			FROM(
				table.Group.
					LEFT_JOIN(table.GroupType, table.GroupType.GroupID.EQ(table.Group.ID)).
					LEFT_JOIN(table.Type, table.Type.ID.EQ(table.GroupType.TypeID)).
					INNER_JOIN(table.Supertype, table.Supertype.ID.EQ(table.Group.SupertypeID)).
					INNER_JOIN(table.Card, table.Card.GroupID.EQ(table.Group.ID)).
					LEFT_JOIN(table.CardCollection, table.CardCollection.CardID.EQ(table.Card.ID)),
			)

		groupFilterStatement := Bool(true)

		if len(query.Categories) > 0 {
			sqlCategories := []Expression{}
			for _, category := range query.Categories {
				sqlCategories = append(sqlCategories, String(category))
			}
			groupFilterStatement = groupFilterStatement.AND(table.Supertype.Name.IN(sqlCategories...))
		}

		if len(query.Types) > 0 {
			sqlTypes := []Expression{}
			for _, typ := range query.Types {
				sqlTypes = append(sqlTypes, String(typ))
			}
			groupFilterStatement = groupFilterStatement.AND(
				table.Group.ID.IN(
					SELECT(table.GroupType.GroupID).
						FROM(table.GroupType.INNER_JOIN(table.Type, table.Type.ID.EQ(table.GroupType.TypeID))).
						WHERE(table.Type.Name.IN(sqlTypes...)),
				),
			)
		}

		if len(query.Group) > 0 {
			groupFilterStatement = groupFilterStatement.AND(table.Group.Name.EQ(String(query.Group)))
		}

		if query.CollectionOnly {
			groupFilterStatement = groupFilterStatement.AND(table.CardCollection.Count.GT(Int(0)))
		}

		stmt = stmt.WHERE(groupFilterStatement)

		groups := []ModelGroup{}
		util.AssertNil(stmt.Query(db, &groups), stmt.DebugSql())

		result := []Group{}
		for groupIdx := range groups {
			var collectionCount int32 = 0
			for collectionIdx := range groups[groupIdx].Collection {
				collectionCount += groups[groupIdx].Collection[collectionIdx].Count
			}

			group := Group{groups[groupIdx], collectionCount}
			result = append(result, group)
		}

		return result
	}()

	type ModelCard struct {
		model.Card
		CardCollection model.CardCollection
	}
	type Card struct {
		ModelCard
		MetadataObj         data.CardMetadata
		CardCollectionCount int32
	}

	filteredCards := func() []Card {
		stmt := SELECT(table.Card.AllColumns, table.CardCollection.AllColumns).
			FROM(table.Card.
				LEFT_JOIN(table.CardCollection, table.CardCollection.CardID.EQ(table.Card.ID)).
				LEFT_JOIN(table.Collection, table.Collection.ID.EQ(table.CardCollection.CollectionID))).
			WHERE(
				AND(
					table.Card.GroupID.IN(
						SELECT(table.Group.ID).
							FROM(table.Group).
							WHERE(table.Group.Name.EQ(String(query.Group))),
					),
					OR(
						table.Collection.Name.EQ(String(COLLECTION_NAME)),
						table.Collection.Name.IS_NULL(),
					),
				),
			)

		cards := []ModelCard{}
		util.AssertNil(stmt.Query(db, &cards), stmt.DebugSql())

		result := []Card{}
		for idx := range cards {
			var metadataObj data.CardMetadata
			util.AssertNil(
				json.Unmarshal([]byte(*cards[idx].Metadata), &metadataObj),
			)
			if cards[idx].PtcgID == card_PtcgID {
				pp.Print(cards[idx])
			}
			card := Card{cards[idx], metadataObj, cards[idx].CardCollection.Count}
			result = append(result, card)
		}

		return result
	}()

	return ctx.Render(
		"cards",
		fiber.Map{
			"Categories":     strCategories,
			"Types":          strTypes,
			"FilteredCards":  filteredCards,
			"FilteredGroups": filteredGroups,
			"URL":            urlBuilder{query},
			"Alphabet":       strings.Split("ABCDEFGHIJKLMNOPQRSTUVWXYZ", ""),
			"CurrentURL":     query.Encode(),
			"RequestID":      uuid.NewString(),
		},
	)
}
