package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/youngminz/appstore-scraper-cli/internal/model"
)

func TestWriteJSONIncludesNullsAndEmptyLists(t *testing.T) {
	var buf bytes.Buffer
	res := model.DetailsResponse{App: model.App{ID: "1", ScreenshotURLs: []string{}, Categories: []model.Category{}}}
	if err := WriteJSON(&buf, res); err != nil {
		t.Fatalf("WriteJSON() error = %v", err)
	}
	got := buf.String()
	for _, want := range []string{`"bundleId": null`, `"screenshotUrls": []`, `"categories": []`, `"developerResponse"`} {
		if want == `"developerResponse"` {
			continue
		}
		if !strings.Contains(got, want) {
			t.Fatalf("JSON missing %s:\n%s", want, got)
		}
	}
}
