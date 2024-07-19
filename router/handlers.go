package router

import (
	"encoding/base64"
	"encoding/json"
	"slices"
	"strings"

	"github.com/Farmer-Pete/HokeyPoke/pokeclient"
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
	Card       string   `url:"card"`
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
func (ub urlBuilder) GetCardJSON(card string) string {
	query := ub.Current.copy()
	query.Card = card
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

	cards := pokeclient.GetCards(query.Categories, query.Types)

	return ctx.Render("search", fiber.Map{
		"Categories":         pokeclient.GetSuperTypes(),
		"Types":              pokeclient.GetTypes(),
		"Cards":              cards,
		"SelectedCardName":   query.Card,
		"SelectedCardImages": pokeclient.GetCardImages(query.Card),
		"URL":                urlBuilder{query},
		"Alphabet":           strings.Split("ABCDEFGHIJKLMNOPQRSTUVWXYZ", ""),
	})
}
