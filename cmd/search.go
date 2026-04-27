package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/youngminz/appstore-scraper-cli/internal/output"
	"github.com/youngminz/appstore-scraper-cli/internal/store"
)

func newSearchCommand(app *appContext) *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:   "search <term> --platform ios|android [flags]",
		Short: "Search for apps by keyword",
		Long:  searchHelp,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("search requires exactly one term")
			}
			if strings.TrimSpace(args[0]) == "" {
				return errors.New("search term cannot be empty")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if limit < 1 || limit > 250 {
				return fmt.Errorf("search limit must be between 1 and 250")
			}
			ctx, cancel, err := app.prepare()
			if err != nil {
				return err
			}
			defer cancel()

			req := store.SearchRequest{
				Query:    args[0],
				Platform: app.globals.platform,
				Country:  app.globals.country,
				Lang:     app.globals.lang,
				Limit:    limit,
			}
			res, err := app.client.Search(ctx, req)
			if err != nil {
				return err
			}
			if app.globals.format == "csv" {
				return output.WriteSearchCSV(app.out, res)
			}
			return app.writeJSON(res)
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of search results")
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(cmd.OutOrStdout(), searchHelp)
	})
	return cmd
}

const searchHelp = `Search for apps by keyword.

Usage:
  appstore-scraper search <term> --platform ios|android [flags]

Examples:
  appstore-scraper search "spotify" --platform ios --country us --limit 5
  appstore-scraper search "photo editor" --platform android --country us

Flags:
      --limit int  Maximum number of search results (default 10, max 250)
  -h, --help       Show help

Global Flags:
      --platform string   Store platform: ios or android
      --country string    Two-letter store country code (default "us")
      --lang string       Language code where supported (default "en")
      --output string     Output format: json or csv (default "json")
      --timeout duration  Request timeout (default 30s)`
