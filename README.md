# appstore-scraper-cli

[![CI](https://github.com/youngminz/appstore-scraper-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/youngminz/appstore-scraper-cli/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/go-1.22%2B-00ADD8.svg)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`appstore-scraper-cli` retrieves public mobile app store data from the Apple App Store and Google Play Store.

It is designed for local scripts and automation jobs that need structured app search results, app metadata, or raw reviews for downstream analysis.

No App Store Connect, Google Play Console, or private API credentials are required.

## Quick Start

```bash
brew install youngminz/tap/appstore-scraper

appstore-scraper search spotify --platform ios --limit 3
appstore-scraper details com.spotify.music --platform android
appstore-scraper reviews 324684580 --platform ios --limit 10
```

## Installation

Requires Go 1.22 or newer.

Homebrew:

```bash
brew install youngminz/tap/appstore-scraper
```

Go:

```bash
go install github.com/youngminz/appstore-scraper-cli@latest
```

From this repository:

```bash
go build -o appstore-scraper .
```

## Usage

```bash
appstore-scraper <command> --platform ios|android [flags]
```

Commands:

- `search <term>` searches apps by keyword.
- `details <app-id>` fetches app metadata by numeric App Store ID, iOS bundle ID, or Android package name.
- `reviews <app-id>` fetches raw app reviews.

Global flags:

- `--platform ios|android` is required.
- `--country <code>` defaults to `us`.
- `--lang <code>` defaults to `en`.
- `--output json|csv` defaults to `json`.
- `--timeout <duration>` defaults to `30s`.

## Examples

Search iOS apps:

```bash
appstore-scraper search "spotify" --platform ios --country us --limit 5
```

Search Android apps:

```bash
appstore-scraper search "photo editor" --platform android --country us --limit 10
```

Fetch iOS details by numeric App Store ID:

```bash
appstore-scraper details 324684580 --platform ios --country us
```

Fetch iOS details by bundle ID:

```bash
appstore-scraper details com.spotify.client --platform ios --country us
```

Fetch Android details:

```bash
appstore-scraper details com.spotify.music --platform android --country us
```

Fetch reviews as JSON:

```bash
appstore-scraper reviews com.spotify.music --platform android --country us --sort newest --limit 200
```

Fetch reviews as CSV:

```bash
appstore-scraper reviews 324684580 --platform ios --country us --limit 100 --output csv
```

## Common Workflows

Find the right store ID before fetching details:

```bash
appstore-scraper search "spotify" --platform ios --country us --limit 5
appstore-scraper search "spotify" --platform android --country us --limit 5
```

Save metadata for later processing:

```bash
appstore-scraper details com.spotify.music --platform android > spotify-android.json
```

Fetch review text for an LLM or labeling workflow:

```bash
appstore-scraper reviews com.spotify.music --platform android --limit 200 \
  | jq '.reviews[] | {rating, text, reviewedAt}'
```

Export reviews to a spreadsheet-friendly CSV:

```bash
appstore-scraper reviews 324684580 --platform ios --limit 500 --output csv > spotify-ios-reviews.csv
```

## Output

JSON is the default output and uses `camelCase` fields. CSV output uses `snake_case` columns and is written with Go's standard CSV writer.

Successful command output is written to stdout. Errors are written to stderr, and invalid input, network failures, scraper failures, and unsupported platform behavior return non-zero exit codes.

Unavailable scalar fields are returned as `null` in JSON. Unavailable list fields are returned as `[]`.

## Command Support

| Command | iOS | Android | Output |
| --- | --- | --- | --- |
| `search <term>` | yes | yes | json, csv |
| `details <app-id>` | yes | yes | json, csv |
| `reviews <app-id>` | yes | yes | json, csv |

## Review Sorting

Supported review sort values:

- iOS: `newest`, `helpfulness`
- Android: `newest`, `rating`, `helpfulness`

Unsupported platform and sort combinations are rejected before making store requests.

## Development

Run unit tests:

```bash
go test ./...
```

Run optional integration tests against public store endpoints:

```bash
go test -tags=integration ./internal/store
```

Build the CLI:

```bash
go build -o appstore-scraper .
```

## Release

Create a tag to publish a GitHub release:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The release workflow builds platform binaries and can update the Homebrew tap at `youngminz/homebrew-tap`.

To enable automatic tap updates, add a repository secret named `HOMEBREW_TAP_TOKEN` with permission to push to `youngminz/homebrew-tap`.

For a manual tap update, compute the source archive SHA-256 and render the Formula:

```bash
curl -fsSL -o source.tar.gz \
  https://github.com/youngminz/appstore-scraper-cli/archive/refs/tags/v0.1.0.tar.gz

shasum -a 256 source.tar.gz

./scripts/update-homebrew-formula.sh v0.1.0 <sha256> ../homebrew-tap/Formula/appstore-scraper.rb
```

## Troubleshooting

If a command is slow or the store is temporarily unavailable, retry with a larger timeout:

```bash
appstore-scraper reviews com.spotify.music --platform android --limit 100 --timeout 60s
```

If Android search returns unexpected results, try setting both country and language explicitly:

```bash
appstore-scraper search "spotify" --platform android --country us --lang en
```

Unsupported review sort combinations are rejected before network calls:

```bash
appstore-scraper reviews 324684580 --platform ios --sort rating
```

## Limitations

- This tool uses public store endpoints and scraping. Store markup and endpoint behavior can change without notice.
- iOS review pagination is limited by public App Store RSS endpoints.
- iOS developer responses are not available through the public review endpoint.
- Google Play sorting and pagination can differ from the public store UI.
- The CLI is not an App Store Connect or Google Play Console client.
- The CLI does not perform ASO scoring, keyword suggestion, LLM review analysis, authentication, private analytics, or historical storage.
