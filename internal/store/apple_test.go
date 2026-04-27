package store

import (
	"encoding/json"
	"testing"
)

func TestNormalizeAppleAppUsesCommonShape(t *testing.T) {
	app := normalizeAppleApp(appleLookupApp{
		TrackID: 324684580, BundleID: "com.spotify.client", TrackName: "Spotify",
		ArtistID: 1, ArtistName: "Spotify", AverageUserRating: 4.8, UserRatingCount: 10,
		Price: 0, Currency: "USD", FormattedPrice: "Free", PrimaryGenreID: 6011, PrimaryGenreName: "Music",
	})
	if app.ID != "324684580" {
		t.Fatalf("ID = %q", app.ID)
	}
	if app.PackageName != nil {
		t.Fatalf("PackageName = %v, want nil", *app.PackageName)
	}
	if app.ScreenshotURLs == nil || app.Categories == nil {
		t.Fatalf("slices must be non-nil")
	}
	raw, err := json.Marshal(app)
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(raw) {
		t.Fatalf("invalid JSON")
	}
}

func TestNormalizeAppleReviewsSkipsFeedMetadataEntry(t *testing.T) {
	feed := appleReviewResponse{}
	feed.Feed.Entry = []appleReviewEntry{{Title: appleLabel{Label: "App name"}}, {ID: appleLabel{Label: "r1"}, Rating: appleLabel{Label: "5"}, Content: appleLabel{Label: "Nice"}}}
	reviews := normalizeAppleReviews(feed)
	if len(reviews) != 1 {
		t.Fatalf("reviews len = %d, want 1", len(reviews))
	}
	if got := *reviews[0].ID; got != "r1" {
		t.Fatalf("review ID = %q", got)
	}
}
