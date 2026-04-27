//go:build integration

package store

import (
	"context"
	"testing"
	"time"
)

func TestIntegrationAppleSearchDetailsReviews(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client := NewAppleClient(HTTPClient(30 * time.Second))

	search, err := client.Search(ctx, SearchRequest{Query: "spotify", Platform: "ios", Country: "us", Lang: "en", Limit: 1})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if search.Count == 0 {
		t.Fatal("Search() returned no results")
	}

	if _, err := client.Details(ctx, DetailsRequest{AppID: "324684580", Platform: "ios", Country: "us", Lang: "en"}); err != nil {
		t.Fatalf("Details() error = %v", err)
	}
	reviews, err := client.Reviews(ctx, ReviewsRequest{AppID: "324684580", Platform: "ios", Country: "us", Lang: "en", Sort: "newest", Limit: 1})
	if err != nil {
		t.Fatalf("Reviews() error = %v", err)
	}
	if reviews.Count == 0 {
		t.Fatal("Reviews() returned no results")
	}
}

func TestIntegrationGoogleSearchDetailsReviews(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client := NewGoogleClient(HTTPClient(30 * time.Second))

	search, err := client.Search(ctx, SearchRequest{Query: "spotify", Platform: "android", Country: "us", Lang: "en", Limit: 1})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if search.Count == 0 {
		t.Fatal("Search() returned no results")
	}

	if _, err := client.Details(ctx, DetailsRequest{AppID: "com.spotify.music", Platform: "android", Country: "us", Lang: "en"}); err != nil {
		t.Fatalf("Details() error = %v", err)
	}
	reviews, err := client.Reviews(ctx, ReviewsRequest{AppID: "com.spotify.music", Platform: "android", Country: "us", Lang: "en", Sort: "newest", Limit: 1})
	if err != nil {
		t.Fatalf("Reviews() error = %v", err)
	}
	if reviews.Count == 0 {
		t.Fatal("Reviews() returned no results")
	}
}
