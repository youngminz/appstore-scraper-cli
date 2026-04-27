package store

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/youngminz/appstore-scraper-cli/internal/model"
)

type SearchRequest struct {
	Query    string
	Platform string
	Country  string
	Lang     string
	Limit    int
}

type DetailsRequest struct {
	AppID    string
	Platform string
	Country  string
	Lang     string
}

type ReviewsRequest struct {
	AppID    string
	Platform string
	Country  string
	Lang     string
	Sort     string
	Limit    int
}

type Client interface {
	Search(context.Context, SearchRequest) (model.SearchResponse, error)
	Details(context.Context, DetailsRequest) (model.DetailsResponse, error)
	Reviews(context.Context, ReviewsRequest) (model.ReviewsResponse, error)
}

type platformClient interface {
	Search(context.Context, SearchRequest) (model.SearchResponse, error)
	Details(context.Context, DetailsRequest) (model.DetailsResponse, error)
	Reviews(context.Context, ReviewsRequest) (model.ReviewsResponse, error)
}

type Router struct {
	ios     platformClient
	android platformClient
}

func NewRouter(httpClient *http.Client) *Router {
	return &Router{
		ios:     NewAppleClient(httpClient),
		android: NewGoogleClient(httpClient),
	}
}

func HTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout}
}

func (r *Router) Search(ctx context.Context, req SearchRequest) (model.SearchResponse, error) {
	client, err := r.client(req.Platform)
	if err != nil {
		return model.SearchResponse{}, err
	}
	return client.Search(ctx, req)
}

func (r *Router) Details(ctx context.Context, req DetailsRequest) (model.DetailsResponse, error) {
	client, err := r.client(req.Platform)
	if err != nil {
		return model.DetailsResponse{}, err
	}
	return client.Details(ctx, req)
}

func (r *Router) Reviews(ctx context.Context, req ReviewsRequest) (model.ReviewsResponse, error) {
	client, err := r.client(req.Platform)
	if err != nil {
		return model.ReviewsResponse{}, err
	}
	return client.Reviews(ctx, req)
}

func (r *Router) client(platform string) (platformClient, error) {
	switch platform {
	case "ios":
		return r.ios, nil
	case "android":
		return r.android, nil
	default:
		return nil, fmt.Errorf("unsupported platform %q", platform)
	}
}
