package router

import (
	"slices"
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

type urlBuilder struct {
	Current pokemonQuery
}

func (ub urlBuilder) CategoryExists(category string) bool {
	return slices.Contains(ub.Current.Categories, category)
}

func (ub urlBuilder) GetToggledCategoryJSON(category string) string {
	query := ub.Current.copy()
	query.Categories = toggleValue(query.Categories, category)
	return query.Encode()
}

func (ub urlBuilder) TypeExists(t string) bool {
	return slices.Contains(ub.Current.Types, t)
}

func (ub urlBuilder) GetToggledTypeJSON(t string) string {
	query := ub.Current.copy()
	query.Types = toggleValue(query.Types, t)
	return query.Encode()
}

func (ub urlBuilder) GetGroupJSON(group string) string {
	query := ub.Current.copy()
	query.Group = group
	return query.Encode()
}

func (ub urlBuilder) GetToggledCollectionOnlyJSON() string {
	query := ub.Current.copy()
	query.CollectionOnly = !ub.Current.CollectionOnly
	return query.Encode()
}
