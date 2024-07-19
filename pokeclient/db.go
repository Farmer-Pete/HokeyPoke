package pokeclient

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/Farmer-Pete/HokeyPoke/util"
	pokemontcgv2 "github.com/PokemonTCG/pokemon-tcg-sdk-go-v2/pkg"
	"github.com/PokemonTCG/pokemon-tcg-sdk-go-v2/pkg/request"
)

const DB_FILE = "cards.json"

type Card struct {
	Supertype string            `json:"supertype"`
	Types     map[string]bool   `json:"type"`
	Name      string            `json:"name"`
	Images    map[string]Images `json:"IDs"`
	Count     int               `json:"count"`
}

type Images struct {
	Large string `json:"large"`
	Small string `json:"small"`
}

func BuildDB() {
	if _, err := os.Stat(DB_FILE); err == nil {
		fmt.Printf("Not rebuilding `%s` as DB file exists already\n", DB_FILE)
		return
	}

	db := map[string]Card{}
	page := 0

	for {
		page++
		fmt.Printf("Fetching: page=%d... ", page)

		cards, err := Client().GetCards(
			request.Query(),
			request.Page(page),
		)
		util.AssertNill(err)

		fmt.Printf("Got: size=%d... ", len(cards))
		if len(cards) == 0 {
			break
		}

		for _, card := range cards {
			var record Card
			var ok bool

			card.Name = normalizeName(*card)

			if record, ok = db[card.Name]; !ok {
				record = Card{
					Supertype: card.Supertype,
					Types:     map[string]bool{},
					Name:      card.Name,
					Images:    map[string]Images{},
					Count:     0,
				}
			}

			record.Count += 1
			record.Images[card.ID] = Images{
				Small: card.Images.Small,
				Large: card.Images.Large,
			}
			for _, t := range card.Types {
				record.Types[t] = true
			}
			db[card.Name] = record
		}

		fmt.Printf("Pokemon Cards: count=%d...\n", len(db))

	}

	f, err := os.Create(DB_FILE)
	util.AssertNill(err)

	db_cards := []Card{}
	for _, c := range db {
		db_cards = append(db_cards, c)
	}

	j := json.NewEncoder(f)
	util.AssertNill(j.Encode(db_cards))

	fmt.Println("saved!")
}

func normalizeName(card pokemontcgv2.PokemonCard) string {
	if card.Supertype != "Pokémon" {
		return card.Name
	}

	parts := strings.FieldsFunc(
		card.Name,
		func(r rune) bool {
			return r == ' ' || r == '-'
		},
	)

	junk := []string{
		"ALOLAN", "BREAK", "GALARIAN", "HISUIAN",
		"RADIANT", "SHINING", "TEAM", "UNOWN",
		"LEGEND", "UNION", "DARK", "LIGHT",
		"V", "VMAX", "VSTAR", "C",
		"E4", "EX", "FB", "G",
		"GL", "GX", "LT.", "LV.X", "M",
		"δ", "◇", "★", "♀", "♂"}
	good := []string{}

	for _, part := range parts {
		if !slices.Contains(junk, part) && !slices.Contains(junk, strings.ToUpper(part)) && !strings.Contains(part, `'s`) {
			good = append(good, part)
		}
	}

	new_name := strings.Join(good, " ")
	if len(strings.Trim(new_name, " ")) == 0 {
		new_name = card.Name
	}

	if card.Name != new_name {
		fmt.Printf("[%s => %s]", card.Name, new_name)
	}
	return new_name
}

var db []Card

func DB() []Card {
	if len(db) == 0 {
		f, err := os.Open(DB_FILE)
		util.AssertNill(err)

		j := json.NewDecoder(f)
		util.AssertNill(j.Decode(&db))
	}

	return db
}

type pair struct {
	Name  string
	Types map[string]string
}

func GetCards(supertypes []string, types []string) []pair {
	cards := DB()
	result := []pair{}
	type_icons := GetTypes()

	for _, card := range cards {
		supertype_ok := slices.Contains(supertypes, card.Supertype)
		type_ok := false

		for t := range card.Types {
			if slices.Contains(types, t) {
				type_ok = true
			}
		}

		supertype_ok = supertype_ok || supertypes == nil || len(supertypes) == 0
		type_ok = type_ok || types == nil || len(types) == 0

		if supertype_ok && type_ok {
			p := pair{
				Name:  card.Name,
				Types: map[string]string{},
			}
			for icon := range card.Types {
				p.Types[icon] = type_icons[icon]
			}
			result = append(result, p)
		}
	}

	slices.SortFunc(
		result,
		func(a pair, b pair) int { return strings.Compare(a.Name, b.Name) },
	)
	return result
}

func GetCardImages(name string) map[string]Images {
	cards := DB()

	if len(name) == 0 {
		return nil
	}

	for _, card := range cards {
		if card.Name == name {
			return card.Images
		}
	}

	return nil

}
