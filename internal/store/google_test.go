package store

import (
	"testing"

	gpapp "github.com/n0madic/google-play-scraper/pkg/app"
	gpreviews "github.com/n0madic/google-play-scraper/pkg/reviews"
	gpstore "github.com/n0madic/google-play-scraper/pkg/store"
)

func TestNormalizeGoogleAppUsesCommonShape(t *testing.T) {
	app := normalizeGoogleApp(&gpapp.App{
		ID: "com.spotify.music", Title: "Spotify", Developer: "Spotify AB", Free: true,
		GenreID: "MUSIC_AND_AUDIO", Genre: "Music & Audio",
		Price: gpapp.Price{Currency: "USD", Value: 0}, Ratings: 12, ReviewsTotalCount: 3,
	})
	if app.ID != "com.spotify.music" {
		t.Fatalf("ID = %q", app.ID)
	}
	if app.BundleID != nil {
		t.Fatalf("BundleID = %v, want nil", *app.BundleID)
	}
	if app.PackageName == nil || *app.PackageName != "com.spotify.music" {
		t.Fatalf("PackageName = %v", app.PackageName)
	}
	if app.ScreenshotURLs == nil || app.Categories == nil {
		t.Fatalf("slices must be non-nil")
	}
}

func TestNormalizeGoogleReviewIncludesDeveloperResponse(t *testing.T) {
	review := normalizeGoogleReview("com.example", &gpreviews.Review{
		ID: "r1", Reviewer: "Alice", Score: 4, Text: "Good", Useful: 2, Reply: "Thanks",
	})
	if review.DeveloperResponse == nil {
		t.Fatal("DeveloperResponse = nil")
	}
	if review.HelpfulCount == nil || *review.HelpfulCount != 2 {
		t.Fatalf("HelpfulCount = %v", review.HelpfulCount)
	}
}

func TestGoogleSort(t *testing.T) {
	tests := map[string]gpstore.Sort{
		"newest":      gpstore.SortNewest,
		"rating":      gpstore.SortRating,
		"helpfulness": gpstore.SortHelpfulness,
	}
	for input, want := range tests {
		if got := googleSort(input); got != want {
			t.Fatalf("googleSort(%q) = %v, want %v", input, got, want)
		}
	}
}
