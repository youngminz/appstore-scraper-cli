package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	gpapp "github.com/n0madic/google-play-scraper/pkg/app"
	gpreviews "github.com/n0madic/google-play-scraper/pkg/reviews"
	gpsearch "github.com/n0madic/google-play-scraper/pkg/search"
	"github.com/n0madic/google-play-scraper/pkg/store"
	"github.com/youngminz/appstore-scraper-cli/internal/model"
)

type GoogleClient struct{}

func NewGoogleClient() *GoogleClient {
	return &GoogleClient{}
}

func (c *GoogleClient) Search(ctx context.Context, req SearchRequest) (model.SearchResponse, error) {
	now := time.Now().UTC()
	scraper := gpsearch.NewQuery(req.Query, gpsearch.PriceAll, gpsearch.Options{
		Country: req.Country, Language: req.Lang, Number: req.Limit,
	})
	if err := runWithContext(ctx, scraper.Run); err != nil {
		return model.SearchResponse{}, fmt.Errorf("android search failed: %w", err)
	}
	_ = runWithContext(ctx, func() error {
		errs := scraper.LoadMoreDetails(4)
		if len(errs) > 0 {
			return errs[0]
		}
		return nil
	})
	results := make([]model.App, 0, len(scraper.Results))
	for _, app := range scraper.Results {
		if app != nil {
			results = append(results, normalizeGoogleApp(app))
		}
	}
	if len(results) > req.Limit {
		results = results[:req.Limit]
	}
	return model.SearchResponse{
		Query: req.Query, Platform: req.Platform, Country: req.Country, Lang: req.Lang,
		Limit: req.Limit, Count: len(results), FetchedAt: now, Results: results,
	}, nil
}

func (c *GoogleClient) Details(ctx context.Context, req DetailsRequest) (model.DetailsResponse, error) {
	now := time.Now().UTC()
	app := gpapp.New(req.AppID, gpapp.Options{Country: req.Country, Language: req.Lang})
	if err := runWithContext(ctx, app.LoadDetails); err != nil {
		return model.DetailsResponse{}, fmt.Errorf("android app not found or details failed: %s: %w", req.AppID, err)
	}
	return model.DetailsResponse{
		AppID: req.AppID, Platform: req.Platform, Country: req.Country, Lang: req.Lang,
		FetchedAt: now, App: normalizeGoogleApp(app),
	}, nil
}

func (c *GoogleClient) Reviews(ctx context.Context, req ReviewsRequest) (model.ReviewsResponse, error) {
	now := time.Now().UTC()
	reviewClient := gpreviews.New(req.AppID, gpreviews.Options{
		Country: req.Country, Language: req.Lang, Number: req.Limit, Sorting: googleSort(req.Sort),
	})
	if err := runWithContext(ctx, reviewClient.Run); err != nil {
		return model.ReviewsResponse{}, fmt.Errorf("android reviews failed: %w", err)
	}
	reviews := make([]model.Review, 0, len(reviewClient.Results))
	for _, review := range reviewClient.Results {
		reviews = append(reviews, normalizeGoogleReview(req.AppID, review))
		if len(reviews) >= req.Limit {
			break
		}
	}
	return model.ReviewsResponse{
		AppID: req.AppID, Platform: req.Platform, Country: req.Country, Lang: req.Lang,
		Sort: req.Sort, Limit: req.Limit, Count: len(reviews), FetchedAt: now, Reviews: reviews,
	}, nil
}

func normalizeGoogleApp(app *gpapp.App) model.App {
	free := app.Free
	formattedPrice := "Free"
	if !free && app.Price.Value > 0 {
		formattedPrice = fmt.Sprintf("%g %s", app.Price.Value, app.Price.Currency)
	}
	categories := []model.Category{}
	if app.GenreID != "" || app.Genre != "" {
		categories = append(categories, model.Category{ID: model.StringPtr(app.GenreID), Name: model.StringPtr(app.Genre)})
	}
	ratingCount := int64(app.Ratings)
	reviewCount := int64(app.ReviewsTotalCount)
	var updatedAt *time.Time
	if !app.Updated.IsZero() {
		updated := app.Updated.UTC()
		updatedAt = &updated
	}
	return model.App{
		ID: app.ID, BundleID: nil, PackageName: model.StringPtr(app.ID),
		Title: model.StringPtr(app.Title), Summary: model.StringPtr(app.Summary), Description: model.StringPtr(app.Description),
		Developer: model.Developer{
			ID:   model.StringPtr(firstNonEmpty(app.DeveloperID, app.DeveloperInternalID, app.Developer)),
			Name: model.StringPtr(app.Developer), Website: model.StringPtr(firstNonEmpty(app.DeveloperWebsite, app.DeveloperURL)),
			Email: model.StringPtr(app.DeveloperEmail),
		},
		IconURL: model.StringPtr(app.Icon), ScreenshotURLs: app.Screenshots,
		Rating:     model.Rating{Score: model.FloatPtr(app.Score), Count: model.Int64Ptr(ratingCount), ReviewCount: model.Int64Ptr(reviewCount), Histogram: intHistogram(app.RatingsHistogram)},
		Pricing:    model.Pricing{Price: model.FloatPtr(app.Price.Value), Currency: model.StringPtr(app.Price.Currency), FormattedPrice: model.StringPtr(formattedPrice), Free: model.BoolPtr(free)},
		Categories: categories, ContentRating: model.StringPtr(app.ContentRating), ReleasedAt: parseTime(app.Released), UpdatedAt: updatedAt,
		Version: model.StringPtr(app.Version), ReleaseNotes: model.StringPtr(app.RecentChanges), StoreURL: model.StringPtr(app.URL),
	}
}

func normalizeGoogleReview(appID string, review *gpreviews.Review) model.Review {
	if review == nil {
		return model.Review{}
	}
	var replied *model.DeveloperResponse
	if review.Reply != "" || !review.ReplyTimestamp.IsZero() {
		var repliedAt *time.Time
		if !review.ReplyTimestamp.IsZero() {
			t := review.ReplyTimestamp.UTC()
			repliedAt = &t
		}
		replied = &model.DeveloperResponse{Text: model.StringPtr(review.Reply), RespondedAt: repliedAt}
	}
	var reviewedAt *time.Time
	if !review.Timestamp.IsZero() {
		t := review.Timestamp.UTC()
		reviewedAt = &t
	}
	helpful := int64(review.Useful)
	return model.Review{
		ID:     model.StringPtr(review.ID),
		User:   model.User{Name: model.StringPtr(review.Reviewer), ImageURL: model.StringPtr(review.Avatar)},
		Rating: model.IntPtr(review.Score), Title: nil, Text: model.StringPtr(review.Text),
		ReviewedAt: reviewedAt, Version: model.StringPtr(review.Version), URL: model.StringPtr(review.URL(appID)),
		HelpfulCount: model.Int64Ptr(helpful), DeveloperResponse: replied,
	}
}

func googleSort(sort string) store.Sort {
	switch sort {
	case "rating":
		return store.SortRating
	case "helpfulness":
		return store.SortHelpfulness
	default:
		return store.SortNewest
	}
}

func intHistogram(in map[int]int) map[int]int64 {
	if len(in) == 0 {
		return nil
	}
	out := make(map[int]int64, len(in))
	for k, v := range in {
		out[k] = int64(v)
	}
	return out
}

func runWithContext(ctx context.Context, fn func() error) error {
	done := make(chan error, 1)
	go func() {
		done <- fn()
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		if err != nil && strings.Contains(err.Error(), "status code: 404") {
			return fmt.Errorf("app not found: %w", err)
		}
		return err
	}
}
