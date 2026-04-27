package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/youngminz/appstore-scraper-cli/internal/output"
	"github.com/youngminz/appstore-scraper-cli/internal/store"
)

func newReviewsCommand(app *appContext) *cobra.Command {
	var limit int
	var sort string
	cmd := &cobra.Command{
		Use:   "reviews <app-id> --platform ios|android [flags]",
		Short: "Fetch raw app reviews",
		Long:  reviewsHelp,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("reviews requires exactly one app ID")
			}
			if strings.TrimSpace(args[0]) == "" {
				return errors.New("app ID cannot be empty")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if limit < 1 || limit > 1000 {
				return fmt.Errorf("review limit must be between 1 and 1000")
			}
			ctx, cancel, err := app.prepare()
			if err != nil {
				return err
			}
			defer cancel()
			if err := validateSort(app.globals.platform, sort); err != nil {
				return err
			}

			req := store.ReviewsRequest{
				AppID:    args[0],
				Platform: app.globals.platform,
				Country:  app.globals.country,
				Lang:     app.globals.lang,
				Sort:     sort,
				Limit:    limit,
			}
			res, err := app.client.Reviews(ctx, req)
			if err != nil {
				return err
			}
			if app.globals.format == "csv" {
				return output.WriteReviewsCSV(app.out, res)
			}
			return app.writeJSON(res)
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 100, "Maximum number of reviews")
	cmd.Flags().StringVar(&sort, "sort", "newest", "Review sort order: newest, rating, or helpfulness")
	return cmd
}

func validateSort(platform, sort string) error {
	switch sort {
	case "newest", "helpfulness":
		return nil
	case "rating":
		if platform == "android" {
			return nil
		}
		return errors.New("unsupported sort \"rating\" for ios: supported sorts are newest and helpfulness")
	default:
		return fmt.Errorf("unsupported sort %q: must be newest, rating, or helpfulness", sort)
	}
}

const reviewsHelp = `Fetch raw app reviews.

Usage:
  appstore-scraper reviews <app-id> --platform ios|android [flags]

Examples:
  appstore-scraper reviews 324684580 --platform ios --country us --limit 100
  appstore-scraper reviews com.spotify.music --platform android --country us --sort newest --output csv

Flags:
      --limit int     Maximum number of reviews (default 100, max 1000)
      --sort string   Review sort order: newest, rating, or helpfulness (default "newest")
  -h, --help          Show help

Sort support:
  ios      newest, helpfulness
  android  newest, rating, helpfulness`
