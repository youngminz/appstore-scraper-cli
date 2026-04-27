package store

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/youngminz/appstore-scraper-cli/internal/model"
)

type AppleClient struct {
	http *http.Client
}

func NewAppleClient(httpClient *http.Client) *AppleClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &AppleClient{http: httpClient}
}

func (c *AppleClient) Search(ctx context.Context, req SearchRequest) (model.SearchResponse, error) {
	now := time.Now().UTC()
	values := url.Values{
		"term":    {req.Query},
		"entity":  {"software"},
		"country": {req.Country},
		"lang":    {appleLang(req.Lang)},
		"limit":   {strconv.Itoa(req.Limit)},
	}
	var lookup appleLookupResponse
	if err := c.getJSON(ctx, "https://itunes.apple.com/search?"+values.Encode(), &lookup); err != nil {
		return model.SearchResponse{}, err
	}
	results := make([]model.App, 0, len(lookup.Results))
	for _, app := range lookup.Results {
		results = append(results, normalizeAppleApp(app))
	}
	return model.SearchResponse{
		Query: req.Query, Platform: req.Platform, Country: req.Country, Lang: req.Lang,
		Limit: req.Limit, Count: len(results), FetchedAt: now, Results: results,
	}, nil
}

func (c *AppleClient) Details(ctx context.Context, req DetailsRequest) (model.DetailsResponse, error) {
	now := time.Now().UTC()
	app, err := c.lookupApp(ctx, req.AppID, req.Country, req.Lang)
	if err != nil {
		return model.DetailsResponse{}, err
	}
	return model.DetailsResponse{
		AppID: req.AppID, Platform: req.Platform, Country: req.Country, Lang: req.Lang,
		FetchedAt: now, App: app,
	}, nil
}

func (c *AppleClient) Reviews(ctx context.Context, req ReviewsRequest) (model.ReviewsResponse, error) {
	now := time.Now().UTC()
	appID := req.AppID
	if _, err := strconv.ParseInt(appID, 10, 64); err != nil {
		app, lookupErr := c.lookupApp(ctx, appID, req.Country, req.Lang)
		if lookupErr != nil {
			return model.ReviewsResponse{}, lookupErr
		}
		appID = app.ID
	}

	reviews := make([]model.Review, 0, req.Limit)
	sort := "mostrecent"
	if req.Sort == "helpfulness" {
		sort = "mosthelpful"
	}
	for page := 1; len(reviews) < req.Limit && page <= 10; page++ {
		endpoint := fmt.Sprintf("https://itunes.apple.com/%s/rss/customerreviews/page=%d/id=%s/sortby=%s/json", strings.ToLower(req.Country), page, url.PathEscape(appID), sort)
		var feed appleReviewResponse
		if err := c.getJSON(ctx, endpoint, &feed); err != nil {
			return model.ReviewsResponse{}, err
		}
		pageReviews := normalizeAppleReviews(feed)
		if len(pageReviews) == 0 {
			break
		}
		for _, review := range pageReviews {
			if len(reviews) >= req.Limit {
				break
			}
			reviews = append(reviews, review)
		}
	}

	return model.ReviewsResponse{
		AppID: req.AppID, Platform: req.Platform, Country: req.Country, Lang: req.Lang,
		Sort: req.Sort, Limit: req.Limit, Count: len(reviews), FetchedAt: now, Reviews: reviews,
	}, nil
}

func (c *AppleClient) lookupApp(ctx context.Context, appID, country, lang string) (model.App, error) {
	values := url.Values{
		"country": {country},
		"lang":    {appleLang(lang)},
	}
	if _, err := strconv.ParseInt(appID, 10, 64); err == nil {
		values.Set("id", appID)
	} else {
		values.Set("bundleId", appID)
	}
	var lookup appleLookupResponse
	if err := c.getJSON(ctx, "https://itunes.apple.com/lookup?"+values.Encode(), &lookup); err != nil {
		return model.App{}, err
	}
	if len(lookup.Results) == 0 {
		return model.App{}, fmt.Errorf("ios app not found: %s", appID)
	}
	return normalizeAppleApp(lookup.Results[0]), nil
}

func (c *AppleClient) getJSON(ctx context.Context, endpoint string, dest any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "appstore-scraper-cli/0.1")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("apple request failed: %s", resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return err
	}
	return nil
}

type appleLookupResponse struct {
	ResultCount int              `json:"resultCount"`
	Results     []appleLookupApp `json:"results"`
}

type appleLookupApp struct {
	TrackID               int64    `json:"trackId"`
	BundleID              string   `json:"bundleId"`
	TrackName             string   `json:"trackName"`
	Description           string   `json:"description"`
	SellerName            string   `json:"sellerName"`
	ArtistID              int64    `json:"artistId"`
	ArtistName            string   `json:"artistName"`
	ArtworkURL512         string   `json:"artworkUrl512"`
	ScreenshotURLs        []string `json:"screenshotUrls"`
	IPadScreenshotURLs    []string `json:"ipadScreenshotUrls"`
	AverageUserRating     float64  `json:"averageUserRating"`
	UserRatingCount       int64    `json:"userRatingCount"`
	Price                 float64  `json:"price"`
	Currency              string   `json:"currency"`
	FormattedPrice        string   `json:"formattedPrice"`
	PrimaryGenreID        int64    `json:"primaryGenreId"`
	PrimaryGenreName      string   `json:"primaryGenreName"`
	ContentAdvisoryRating string   `json:"contentAdvisoryRating"`
	ReleaseDate           string   `json:"releaseDate"`
	CurrentVersionDate    string   `json:"currentVersionReleaseDate"`
	Version               string   `json:"version"`
	ReleaseNotes          string   `json:"releaseNotes"`
	TrackViewURL          string   `json:"trackViewUrl"`
}

func normalizeAppleApp(app appleLookupApp) model.App {
	id := strconv.FormatInt(app.TrackID, 10)
	developerID := ""
	if app.ArtistID != 0 {
		developerID = strconv.FormatInt(app.ArtistID, 10)
	}
	free := app.Price == 0
	screenshots := app.ScreenshotURLs
	if len(screenshots) == 0 {
		screenshots = app.IPadScreenshotURLs
	}
	categories := []model.Category{}
	if app.PrimaryGenreName != "" || app.PrimaryGenreID != 0 {
		categoryID := ""
		if app.PrimaryGenreID != 0 {
			categoryID = strconv.FormatInt(app.PrimaryGenreID, 10)
		}
		categories = append(categories, model.Category{ID: model.StringPtr(categoryID), Name: model.StringPtr(app.PrimaryGenreName)})
	}
	return model.App{
		ID: id, BundleID: model.StringPtr(app.BundleID), PackageName: nil,
		Title: model.StringPtr(app.TrackName), Summary: nil, Description: model.StringPtr(app.Description),
		Developer: model.Developer{ID: model.StringPtr(developerID), Name: model.StringPtr(firstNonEmpty(app.ArtistName, app.SellerName))},
		IconURL:   model.StringPtr(app.ArtworkURL512), ScreenshotURLs: nonNilStrings(screenshots),
		Rating:     model.Rating{Score: model.FloatPtr(app.AverageUserRating), Count: model.Int64Ptr(app.UserRatingCount), Histogram: nil},
		Pricing:    model.Pricing{Price: model.FloatPtr(app.Price), Currency: model.StringPtr(app.Currency), FormattedPrice: model.StringPtr(app.FormattedPrice), Free: model.BoolPtr(free)},
		Categories: nonNilCategories(categories), ContentRating: model.StringPtr(app.ContentAdvisoryRating),
		ReleasedAt: parseTime(app.ReleaseDate), UpdatedAt: parseTime(app.CurrentVersionDate),
		Version: model.StringPtr(app.Version), ReleaseNotes: model.StringPtr(app.ReleaseNotes), StoreURL: model.StringPtr(app.TrackViewURL),
	}
}

type appleReviewResponse struct {
	Feed struct {
		Entry []appleReviewEntry `json:"entry"`
	} `json:"feed"`
}

type appleReviewEntry struct {
	ID     appleLabel `json:"id"`
	Author struct {
		Name appleLabel `json:"name"`
		URI  appleLabel `json:"uri"`
	} `json:"author"`
	Rating  appleLabel `json:"im:rating"`
	Version appleLabel `json:"im:version"`
	Title   appleLabel `json:"title"`
	Content appleLabel `json:"content"`
	Updated appleLabel `json:"updated"`
	Link    struct {
		Attributes struct {
			Href string `json:"href"`
		} `json:"attributes"`
	} `json:"link"`
}

type appleLabel struct {
	Label string `json:"label"`
}

func normalizeAppleReviews(feed appleReviewResponse) []model.Review {
	entries := feed.Feed.Entry
	if len(entries) == 0 {
		return nil
	}
	reviews := make([]model.Review, 0, len(entries))
	for _, entry := range entries {
		if entry.Rating.Label == "" && entry.Content.Label == "" {
			continue
		}
		rating, _ := strconv.Atoi(entry.Rating.Label)
		reviewedAt := parseTime(entry.Updated.Label)
		reviews = append(reviews, model.Review{
			ID: model.StringPtr(entry.ID.Label),
			User: model.User{
				Name: model.StringPtr(entry.Author.Name.Label),
				URL:  model.StringPtr(entry.Author.URI.Label),
			},
			Rating:       model.IntPtr(rating),
			Title:        model.StringPtr(entry.Title.Label),
			Text:         model.StringPtr(entry.Content.Label),
			ReviewedAt:   reviewedAt,
			Version:      model.StringPtr(entry.Version.Label),
			URL:          model.StringPtr(entry.Link.Attributes.Href),
			HelpfulCount: nil, DeveloperResponse: nil,
		})
	}
	return reviews
}

func appleLang(lang string) string {
	if strings.Contains(lang, "_") {
		return lang
	}
	return lang + "_us"
}

func parseTime(value string) *time.Time {
	if value == "" {
		return nil
	}
	layouts := []string{time.RFC3339, "2006-01-02T15:04:05Z07:00", "January 2, 2006", "Jan 2, 2006", "2006-01-02"}
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			utc := parsed.UTC()
			return &utc
		}
	}
	return nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
