package output

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/youngminz/appstore-scraper-cli/internal/model"
)

var searchHeader = []string{"platform", "country", "lang", "query", "limit", "fetched_at", "id", "bundle_id", "package_name", "title", "developer_id", "developer_name", "icon_url", "rating_score", "rating_count", "price", "currency", "formatted_price", "free", "store_url"}
var detailsHeader = []string{"platform", "country", "lang", "app_id", "fetched_at", "id", "bundle_id", "package_name", "title", "summary", "description", "developer_id", "developer_name", "developer_website", "developer_email", "icon_url", "screenshot_urls", "rating_score", "rating_count", "review_count", "price", "currency", "formatted_price", "free", "categories", "content_rating", "released_at", "updated_at", "version", "release_notes", "store_url"}
var reviewsHeader = []string{"platform", "country", "lang", "app_id", "sort", "limit", "fetched_at", "review_id", "user_name", "rating", "title", "text", "reviewed_at", "version", "url", "helpful_count", "developer_response_text", "developer_responded_at"}

func WriteSearchCSV(w io.Writer, res model.SearchResponse) error {
	cw := csv.NewWriter(w)
	if err := cw.Write(searchHeader); err != nil {
		return err
	}
	for _, app := range res.Results {
		if err := cw.Write([]string{
			res.Platform, res.Country, res.Lang, res.Query, fmt.Sprint(res.Limit), timeString(&res.FetchedAt),
			app.ID, str(app.BundleID), str(app.PackageName), str(app.Title), str(app.Developer.ID), str(app.Developer.Name),
			str(app.IconURL), float(app.Rating.Score), int64str(app.Rating.Count), float(app.Pricing.Price),
			str(app.Pricing.Currency), str(app.Pricing.FormattedPrice), boolstr(app.Pricing.Free), str(app.StoreURL),
		}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func WriteDetailsCSV(w io.Writer, res model.DetailsResponse) error {
	cw := csv.NewWriter(w)
	if err := cw.Write(detailsHeader); err != nil {
		return err
	}
	app := res.App
	row := []string{
		res.Platform, res.Country, res.Lang, res.AppID, timeString(&res.FetchedAt),
		app.ID, str(app.BundleID), str(app.PackageName), str(app.Title), str(app.Summary), str(app.Description),
		str(app.Developer.ID), str(app.Developer.Name), str(app.Developer.Website), str(app.Developer.Email),
		str(app.IconURL), strings.Join(app.ScreenshotURLs, "|"), float(app.Rating.Score), int64str(app.Rating.Count),
		int64str(app.Rating.ReviewCount), float(app.Pricing.Price), str(app.Pricing.Currency), str(app.Pricing.FormattedPrice),
		boolstr(app.Pricing.Free), joinCategories(app.Categories), str(app.ContentRating), timeString(app.ReleasedAt),
		timeString(app.UpdatedAt), str(app.Version), str(app.ReleaseNotes), str(app.StoreURL),
	}
	if err := cw.Write(row); err != nil {
		return err
	}
	cw.Flush()
	return cw.Error()
}

func WriteReviewsCSV(w io.Writer, res model.ReviewsResponse) error {
	cw := csv.NewWriter(w)
	if err := cw.Write(reviewsHeader); err != nil {
		return err
	}
	for _, review := range res.Reviews {
		devText := ""
		devAt := ""
		if review.DeveloperResponse != nil {
			devText = str(review.DeveloperResponse.Text)
			devAt = timeString(review.DeveloperResponse.RespondedAt)
		}
		row := []string{
			res.Platform, res.Country, res.Lang, res.AppID, res.Sort, fmt.Sprint(res.Limit), timeString(&res.FetchedAt),
			str(review.ID), str(review.User.Name), intstr(review.Rating), str(review.Title), str(review.Text),
			timeString(review.ReviewedAt), str(review.Version), str(review.URL), int64str(review.HelpfulCount), devText, devAt,
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func str(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func float(v *float64) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%g", *v)
}

func intstr(v *int) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(*v)
}

func int64str(v *int64) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(*v)
}

func boolstr(v *bool) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(*v)
}

func timeString(v *time.Time) string {
	if v == nil || v.IsZero() {
		return ""
	}
	return v.UTC().Format(time.RFC3339)
}

func joinCategories(categories []model.Category) string {
	parts := make([]string, 0, len(categories))
	for _, category := range categories {
		name := str(category.Name)
		if name == "" {
			name = str(category.ID)
		}
		if name != "" {
			parts = append(parts, name)
		}
	}
	return strings.Join(parts, "|")
}
