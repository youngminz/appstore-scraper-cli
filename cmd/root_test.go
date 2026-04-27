package cmd

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/youngminz/appstore-scraper-cli/internal/model"
	"github.com/youngminz/appstore-scraper-cli/internal/store"
)

type fakeClient struct {
	searchCalled bool
}

func (f *fakeClient) Search(ctx context.Context, req store.SearchRequest) (model.SearchResponse, error) {
	f.searchCalled = true
	return model.SearchResponse{Query: req.Query, Platform: req.Platform, Country: req.Country, Lang: req.Lang, Limit: req.Limit, Results: []model.App{}, Count: 0}, nil
}

func (f *fakeClient) Details(ctx context.Context, req store.DetailsRequest) (model.DetailsResponse, error) {
	return model.DetailsResponse{AppID: req.AppID, Platform: req.Platform, Country: req.Country, Lang: req.Lang, App: model.App{ID: req.AppID, ScreenshotURLs: []string{}, Categories: []model.Category{}}}, nil
}

func (f *fakeClient) Reviews(ctx context.Context, req store.ReviewsRequest) (model.ReviewsResponse, error) {
	return model.ReviewsResponse{AppID: req.AppID, Platform: req.Platform, Country: req.Country, Lang: req.Lang, Sort: req.Sort, Limit: req.Limit, Reviews: []model.Review{}}, nil
}

func TestSearchValidationRejectsMissingPlatformBeforeNetwork(t *testing.T) {
	var stdout, stderr bytes.Buffer
	fake := &fakeClient{}
	app := &appContext{out: &stdout, errOut: &stderr, client: fake}
	root := newRootCommand(app)
	root.SetArgs([]string{"search", "spotify"})
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error")
	}
	if fake.searchCalled {
		t.Fatal("network client was called")
	}
	if !strings.Contains(err.Error(), "unsupported platform") {
		t.Fatalf("error = %v", err)
	}
}

func TestReviewsRejectsUnsupportedIOSSortBeforeNetwork(t *testing.T) {
	var stdout, stderr bytes.Buffer
	fake := &fakeClient{}
	app := &appContext{out: &stdout, errOut: &stderr, client: fake}
	root := newRootCommand(app)
	root.SetArgs([]string{"reviews", "324684580", "--platform", "ios", "--sort", "rating"})
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unsupported sort") {
		t.Fatalf("error = %v", err)
	}
}
