# appstore-scraper-cli

`appstore-scraper-cli` retrieves public mobile app store data from the Apple App Store and Google Play Store.

It is designed for local scripts and automation jobs that need structured app search results, app metadata, or raw reviews for downstream analysis.

## Installation

Requires Go 1.22 or newer.

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

## Output

JSON is the default output and uses `camelCase` fields. CSV output uses `snake_case` columns and is written with Go's standard CSV writer.

Successful command output is written to stdout. Errors are written to stderr, and invalid input, network failures, scraper failures, and unsupported platform behavior return non-zero exit codes.

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

## Limitations

- This tool uses public store endpoints and scraping. Store markup and endpoint behavior can change without notice.
- iOS review pagination is limited by public App Store RSS endpoints.
- iOS developer responses are not available through the public review endpoint.
- Google Play sorting and pagination can differ from the public store UI.
- The CLI does not perform ASO scoring, keyword suggestion, LLM review analysis, authentication, or historical storage.
