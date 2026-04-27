package output

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"
	"time"

	"github.com/youngminz/appstore-scraper-cli/internal/model"
)

func TestWriteReviewsCSVPreservesEscapedText(t *testing.T) {
	fetchedAt := time.Date(2026, 4, 27, 11, 30, 0, 0, time.UTC)
	reviewedAt := time.Date(2026, 4, 27, 0, 0, 0, 0, time.UTC)
	res := model.ReviewsResponse{
		AppID: "com.example", Platform: "android", Country: "us", Lang: "en", Sort: "newest", Limit: 1, FetchedAt: fetchedAt,
		Reviews: []model.Review{{
			ID: model.StringPtr("r1"), User: model.User{Name: model.StringPtr("Alice")}, Rating: model.IntPtr(5),
			Title: model.StringPtr("Great, really"), Text: model.StringPtr("line one\n\"line two\""), ReviewedAt: &reviewedAt,
		}},
	}
	var buf bytes.Buffer
	if err := WriteReviewsCSV(&buf, res); err != nil {
		t.Fatalf("WriteReviewsCSV() error = %v", err)
	}
	records, err := csv.NewReader(strings.NewReader(buf.String())).ReadAll()
	if err != nil {
		t.Fatalf("csv parse error = %v\n%s", err, buf.String())
	}
	if len(records) != 2 {
		t.Fatalf("records len = %d, want 2", len(records))
	}
	if got := records[1][11]; got != "line one\n\"line two\"" {
		t.Fatalf("text = %q", got)
	}
}

func TestWriteDetailsCSVJoinsLists(t *testing.T) {
	fetchedAt := time.Date(2026, 4, 27, 11, 30, 0, 0, time.UTC)
	res := model.DetailsResponse{
		AppID: "1", Platform: "ios", Country: "us", Lang: "en", FetchedAt: fetchedAt,
		App: model.App{
			ID: "1", ScreenshotURLs: []string{"a", "b"},
			Categories: []model.Category{{ID: model.StringPtr("6011"), Name: model.StringPtr("Music")}},
		},
	}
	var buf bytes.Buffer
	if err := WriteDetailsCSV(&buf, res); err != nil {
		t.Fatalf("WriteDetailsCSV() error = %v", err)
	}
	records, err := csv.NewReader(strings.NewReader(buf.String())).ReadAll()
	if err != nil {
		t.Fatalf("csv parse error = %v", err)
	}
	if got := records[1][16]; got != "a|b" {
		t.Fatalf("screenshot_urls = %q", got)
	}
	if got := records[1][24]; got != "Music" {
		t.Fatalf("categories = %q", got)
	}
}
