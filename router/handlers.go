package router

import (
	"encoding/base64"
	"encoding/json"
	"slices"
	"strings"

	"github.com/Farmer-Pete/HokeyPoke/data"
	"github.com/Farmer-Pete/HokeyPoke/db"
	"github.com/Farmer-Pete/HokeyPoke/db/query/model"
	"github.com/Farmer-Pete/HokeyPoke/db/query/table"
	"github.com/Farmer-Pete/HokeyPoke/util"
	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/k0kubun/pp/v3"
)

func toggleValue(slice []string, value string) []string {
	if slices.Contains(slice, value) {
		slice = slices.DeleteFunc(
			slice,
			func(v string) bool { return v == value },
		)
	} else {
		slice = append(slice, value)
	}
	return slice
}

type pokemonQuery struct {
	Categories []string `url:"categories"`
	Types      []string `url:"types"`
	Group      string   `url:"group"`
}

func (pq *pokemonQuery) copy() *pokemonQuery {
	new := pokemonQuery{}
	new.Categories = slices.Clone(pq.Categories)
	new.Types = slices.Clone(pq.Types)
	return &new
}

type urlBuilder struct {
	Current pokemonQuery
}

func (ub urlBuilder) CategoryExists(category string) bool {
	return slices.Contains(ub.Current.Categories, category)
}
func (ub urlBuilder) GetToggledCategoryJSON(category string) string {
	query := ub.Current.copy()
	query.Categories = toggleValue(query.Categories, category)
	result, err := json.Marshal(query)
	util.AssertNil(err)
	return base64.StdEncoding.EncodeToString(result)
}
func (ub urlBuilder) TypeExists(t string) bool {
	return slices.Contains(ub.Current.Types, t)
}
func (ub urlBuilder) GetToggledTypeJSON(t string) string {
	query := ub.Current.copy()
	query.Types = toggleValue(query.Types, t)
	result, err := json.Marshal(query)
	util.AssertNil(err)
	return base64.StdEncoding.EncodeToString(result)
}
func (ub urlBuilder) GetGroupJSON(group string) string {
	query := ub.Current.copy()
	query.Group = group
	result, err := json.Marshal(query)
	util.AssertNil(err)
	return base64.StdEncoding.EncodeToString(result)
}

func home(ctx *fiber.Ctx) error {
	var query pokemonQuery
	qString := ctx.Queries()
	qBytes, err := base64.StdEncoding.DecodeString(qString["q"])
	util.AssertNil(err)
	json.Unmarshal(qBytes, &query)

	pp.Println(query)

	db := db.GetConnection()

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

	type Group struct {
		model.Group
		Types []model.Type
	}

	filteredGroups := func() []Group {

		stmt := SELECT(
			table.Group.AllColumns,
			table.GroupType.AllColumns,
			table.Type.AllColumns,
			table.Supertype.AllColumns).
			FROM(
				table.Group.
					LEFT_JOIN(table.GroupType, table.GroupType.GroupID.EQ(table.Group.ID)).
					LEFT_JOIN(table.Type, table.Type.ID.EQ(table.GroupType.TypeID)).
					INNER_JOIN(table.Supertype, table.Supertype.ID.EQ(table.Group.SupertypeID)),
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

		stmt = stmt.WHERE(groupFilterStatement)

		result := []Group{}
		print(stmt.DebugSql())
		util.AssertNil(stmt.Query(db, &result))

		return result
	}()

	type Card struct {
		model.Card
		MetadataObj data.CardMetadata
	}

	filteredCards := func() []Card {
		stmt := SELECT(table.Card.AllColumns).
			FROM(table.Card).
			WHERE(table.Card.GroupID.IN(
				SELECT(table.Group.ID).
					FROM(table.Group).
					WHERE(table.Group.Name.EQ(String(query.Group))),
			))

		cards := []model.Card{}
		util.AssertNil(stmt.Query(db, &cards))

		result := []Card{}
		for idx := range cards {
			var metadataObj data.CardMetadata
			util.AssertNil(
				json.Unmarshal([]byte(*cards[idx].Metadata), &metadataObj),
			)
			card := Card{cards[idx], metadataObj}
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
		},
	)
}
