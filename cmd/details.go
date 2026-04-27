package cmd

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"github.com/youngminz/appstore-scraper-cli/internal/output"
	"github.com/youngminz/appstore-scraper-cli/internal/store"
)

func newDetailsCommand(app *appContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "details <app-id> --platform ios|android [flags]",
		Short: "Fetch app metadata by store ID, bundle ID, or package name",
		Long:  detailsHelp,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("details requires exactly one app ID")
			}
			if strings.TrimSpace(args[0]) == "" {
				return errors.New("app ID cannot be empty")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel, err := app.prepare()
			if err != nil {
				return err
			}
			defer cancel()

			req := store.DetailsRequest{
				AppID:    args[0],
				Platform: app.globals.platform,
				Country:  app.globals.country,
				Lang:     app.globals.lang,
			}
			res, err := app.client.Details(ctx, req)
			if err != nil {
				return err
			}
			if app.globals.format == "csv" {
				return output.WriteDetailsCSV(app.out, res)
			}
			return app.writeJSON(res)
		},
	}
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Println(detailsHelp)
	})
	return cmd
}

const detailsHelp = `Fetch app metadata by store ID, bundle ID, or package name.

Usage:
  appstore-scraper details <app-id> --platform ios|android [flags]

Examples:
  appstore-scraper details 324684580 --platform ios --country us
  appstore-scraper details com.spotify.client --platform ios --country us
  appstore-scraper details com.spotify.music --platform android --country us

Flags:
  -h, --help  Show help

Global Flags:
      --platform string   Store platform: ios or android
      --country string    Two-letter store country code (default "us")
      --lang string       Language code where supported (default "en")
      --output string     Output format: json or csv (default "json")
      --timeout duration  Request timeout (default 30s)`
