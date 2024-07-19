package pokeclient

import (
	"os"

	"github.com/Farmer-Pete/HokeyPoke/util"
	pokemontcg "github.com/PokemonTCG/pokemon-tcg-sdk-go-v2/pkg"
)

var client pokemontcg.Client = nil

func Client() pokemontcg.Client {
	if client == nil {
		token := os.Getenv("HOKEY_POKEMONTCG_TOKEN")
		util.AssertNotNill(token, "No token token found (`HOKEY_POKEMONTCG_TOKEN`)")
		client = pokemontcg.NewClient(token)
		util.AssertNotNill(client, "Unable to create client")
	}
	return client
}

func GetSuperTypes() map[string]string {
	result := map[string]string{}
	result["Energy"] = "square-plus"
	result["Pokémon"] = "circle-dot"
	result["Trainer"] = "graduation-cap"
	return result
}

func GetTypes() map[string]string {
	result := map[string]string{}
	result["Colorless"] = "star"
	result["Darkness"] = "moon"
	result["Dragon"] = "dragon"
	result["Fairy"] = "wand-sparkles"
	result["Fighting"] = "hand-fist"
	result["Fire"] = "fire"
	result["Grass"] = "leaf"
	result["Lightning"] = "bolt-lightning"
	result["Metal"] = "gear"
	result["Psychic"] = "cloud"
	result["Water"] = "droplet"
	return result
}
