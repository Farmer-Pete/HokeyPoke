package data

import (
	"slices"
	"strings"
)

func NormalizeName(card CardMetadata) string {
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

	return new_name
}
