package store

import "github.com/youngminz/appstore-scraper-cli/internal/model"

func nonNilStrings(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}

func nonNilCategories(values []model.Category) []model.Category {
	if values == nil {
		return []model.Category{}
	}
	return values
}
