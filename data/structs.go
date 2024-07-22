package data

type CardMetadata struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Supertype   string   `json:"supertype,omitempty"`
	Subtypes    []string `json:"subtypes,omitempty"`
	HP          string   `json:"hp,omitempty"`
	Types       []string `json:"types,omitempty"`
	EvolvesFrom string   `json:"evolvesFrom,omitempty"`
	Attacks     []struct {
		Name                string   `json:"name,omitempty"`
		Cost                []string `json:"cost,omitempty"`
		ConvertedEnergyCost int      `json:"convertedEnergyCost,omitempty"`
		Damage              string   `json:"damage,omitempty"`
		Text                string   `json:"text,omitempty"`
	} `json:"attacks,omitempty"`
	Weaknesses []struct {
		Type  string `json:"type,omitempty"`
		Value string `json:"value,omitempty"`
	} `json:"weaknesses,omitempty"`
	RetreatCost                []string          `json:"retreatCost,omitempty"`
	ConvertedRetreatCost       int               `json:"convertedRetreatCost,omitempty"`
	Number                     string            `json:"number,omitempty"`
	Artist                     string            `json:"artist,omitempty"`
	Rarity                     string            `json:"rarity,omitempty"`
	NationtionalPokedexNumbers []int             `json:"nationalPokedexNumbers,omitempty"`
	Legalities                 map[string]string `json:"legalities,omitempty"`
	Images                     struct {
		Small string `json:"small,omitempty"`
		Large string `json:"large,omitempty"`
	} `json:"images,omitempty"`
}
