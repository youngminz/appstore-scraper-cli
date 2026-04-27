package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/youngminz/appstore-scraper-cli/internal/output"
	"github.com/youngminz/appstore-scraper-cli/internal/store"
)

type globalOptions struct {
	platform string
	country  string
	lang     string
	format   string
	timeout  time.Duration
}

type appContext struct {
	globals globalOptions
	out     io.Writer
	errOut  io.Writer
	client  store.Client
}

func Execute(args []string, stdout, stderr io.Writer) error {
	app := &appContext{out: stdout, errOut: stderr}
	root := newRootCommand(app)
	root.SetArgs(args)
	root.SetOut(stdout)
	root.SetErr(stderr)
	if err := root.Execute(); err != nil {
		fmt.Fprintln(stderr, err)
		return err
	}
	return nil
}

func newRootCommand(app *appContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "appstore-scraper <command> [flags]",
		Short:         "appstore-scraper retrieves public mobile app store data.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.Long = rootHelp

	cmd.PersistentFlags().StringVar(&app.globals.platform, "platform", "", "Store platform: ios or android")
	cmd.PersistentFlags().StringVar(&app.globals.country, "country", "us", "Two-letter store country code")
	cmd.PersistentFlags().StringVar(&app.globals.lang, "lang", "en", "Language code where supported")
	cmd.PersistentFlags().StringVar(&app.globals.format, "output", "json", "Output format: json or csv")
	cmd.PersistentFlags().DurationVar(&app.globals.timeout, "timeout", 30*time.Second, "Request timeout")

	cmd.AddCommand(newSearchCommand(app), newDetailsCommand(app), newReviewsCommand(app))
	return cmd
}

func (app *appContext) prepare() (context.Context, context.CancelFunc, error) {
	if err := validateGlobalOptions(app.globals); err != nil {
		return nil, nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), app.globals.timeout)
	if app.client == nil {
		app.client = store.NewRouter(store.HTTPClient(app.globals.timeout))
	}
	return ctx, cancel, nil
}

func (app *appContext) writeJSON(v any) error {
	return output.WriteJSON(app.out, v)
}

func validateGlobalOptions(opts globalOptions) error {
	switch opts.platform {
	case "ios", "android":
	default:
		return fmt.Errorf("unsupported platform %q: must be ios or android", opts.platform)
	}
	if strings.TrimSpace(opts.country) == "" {
		return errors.New("country is required")
	}
	if strings.TrimSpace(opts.lang) == "" {
		return errors.New("lang is required")
	}
	switch opts.format {
	case "json", "csv":
	default:
		return fmt.Errorf("unsupported output %q: must be json or csv", opts.format)
	}
	if opts.timeout <= 0 {
		return errors.New("timeout must be greater than 0")
	}
	return nil
}

const rootHelp = `appstore-scraper retrieves public mobile app store data.

Usage:
  appstore-scraper <command> [flags]

Commands:
  search    Search for apps by keyword
  details   Fetch app metadata by store ID, bundle ID, or package name
  reviews   Fetch raw app reviews

Flags:
      --platform string   Store platform: ios or android
      --country string    Two-letter store country code (default "us")
      --lang string       Language code where supported (default "en")
      --output string     Output format: json or csv (default "json")
      --timeout duration  Request timeout (default 30s)
  -h, --help              Show help`
