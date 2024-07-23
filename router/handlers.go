package router

import (
	"encoding/base64"
	"encoding/json"
	"slices"
	"strings"

	"github.com/Farmer-Pete/HokeyPoke/dal"
	"github.com/Farmer-Pete/HokeyPoke/dal/db/card"
	"github.com/Farmer-Pete/HokeyPoke/dal/db/group"
	"github.com/Farmer-Pete/HokeyPoke/dal/db/predicate"
	"github.com/Farmer-Pete/HokeyPoke/dal/db/supertype"
	"github.com/Farmer-Pete/HokeyPoke/dal/db/typ3"
	"github.com/Farmer-Pete/HokeyPoke/util"
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
	util.AssertNill(err)
	return base64.StdEncoding.EncodeToString(result)
}
func (ub urlBuilder) TypeExists(t string) bool {
	return slices.Contains(ub.Current.Types, t)
}
func (ub urlBuilder) GetToggledTypeJSON(t string) string {
	query := ub.Current.copy()
	query.Types = toggleValue(query.Types, t)
	result, err := json.Marshal(query)
	util.AssertNill(err)
	return base64.StdEncoding.EncodeToString(result)
}
func (ub urlBuilder) GetGroupJSON(group string) string {
	query := ub.Current.copy()
	query.Group = group
	result, err := json.Marshal(query)
	util.AssertNill(err)
	return base64.StdEncoding.EncodeToString(result)
}

func home(ctx *fiber.Ctx) error {
	var query pokemonQuery
	qString := ctx.Queries()
	qBytes, err := base64.StdEncoding.DecodeString(qString["q"])
	util.AssertNill(err)
	json.Unmarshal(qBytes, &query)

	pp.Println(query)

	db, dbCtx := dal.GetConnection()
	db = db.Debug()
	categories := db.Supertype.Query().Select("name").StringsX(*dbCtx)
	types := db.Typ3.Query().Select("name").StringsX(*dbCtx)

	var cardPredicates []predicate.Card
	var groupPredicates []predicate.Group

	if len(query.Categories) > 0 {
		groupPredicates = append(
			groupPredicates,
			group.HasSupertypeWith(supertype.NameIn(query.Categories...)),
		)
	}
	if len(query.Types) > 0 {
		groupPredicates = append(
			groupPredicates,
			group.HasTypesWith(typ3.NameIn(query.Types...)),
		)
	}
	if len(query.Group) > 0 {
		cardPredicates = append(
			cardPredicates,
			card.HasGroupWith(group.Name(query.Group)),
		)
		groupPredicates = append(
			groupPredicates,
			group.NameEQ(query.Group),
		)
	}

	filteredCards := db.Card.Query().
		Where(cardPredicates...).
		WithTypes().
		AllX(*dbCtx)
	filteredGroups := db.Group.Query().
		Where(groupPredicates...).
		WithTypes().
		AllX(*dbCtx)

	return ctx.Render(
		"cards",
		fiber.Map{
			"Categories":     categories,
			"Types":          types,
			"FilteredCards":  filteredCards,
			"FilteredGroups": filteredGroups,
			"URL":            urlBuilder{query},
			"Alphabet":       strings.Split("ABCDEFGHIJKLMNOPQRSTUVWXYZ", ""),
		},
	)
}
