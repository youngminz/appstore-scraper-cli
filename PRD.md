# Product Requirements Document: appstore-scraper-cli

## Overview

`appstore-scraper-cli` is a Go command-line tool for retrieving mobile app store data from the Apple App Store and Google Play Store. The initial version focuses on three reliable workflows:

- Search for apps by keyword.
- Fetch detailed app metadata.
- Fetch raw user reviews for downstream LLM-based analysis.

The tool should be simple, scriptable, and suitable for use in local workflows and automation jobs.

## Problem

LLM-based review analysis needs clean, structured app store data. For this project, the needed value is narrow:

- Find an app reliably across iOS and Android stores.
- Retrieve app metadata needed for context.
- Collect reviews in structured JSON for later analysis by an LLM.

## Goals

- Provide a fast Go CLI for app search, app details, and review fetching.
- Support both Apple App Store and Google Play Store.
- Return stable JSON output by default for easy piping into scripts and LLM workflows.
- Support CSV output for large exports.
- Keep command and package boundaries clear enough for maintainable CLI development.

## Non-Goals

- Do not implement ASO keyword scoring or keyword suggestion strategies.
- Do not perform LLM review analysis inside the CLI.
- Do not store historical app data or maintain a database.
- Do not support desktop app stores, Steam, Amazon Appstore, or other marketplaces in v1.
- Do not attempt authenticated App Store Connect or Google Play Console APIs.

## Primary Use Cases

### 1. Search for an app

User searches by keyword and platform to identify the correct app ID.

Example:

```bash
appstore-scraper search "spotify" --platform ios --country us --limit 5
```

Expected output:

```json
{
  "query": "spotify",
  "platform": "ios",
  "country": "us",
  "lang": "en",
  "limit": 5,
  "count": 5,
  "fetchedAt": "2026-04-27T11:30:00Z",
  "results": [
    {
      "id": "324684580",
      "bundleId": "com.spotify.client",
      "packageName": null,
      "title": "Spotify - Music and Podcasts",
      "developer": {
        "id": "324684580",
        "name": "Spotify"
      },
      "iconUrl": "https://...",
      "rating": {
        "score": 4.8,
        "count": null
      },
      "pricing": {
        "price": 0,
        "currency": "USD",
        "formattedPrice": "Free",
        "free": true
      },
      "storeUrl": "https://apps.apple.com/..."
    }
  ]
}
```

### 2. Fetch app details

User fetches store metadata for a known app.

Example:

```bash
appstore-scraper details com.spotify.client --platform ios --country us --lang en
appstore-scraper details com.spotify.music --platform android --country us --lang en
```

Expected output:

```json
{
  "appId": "com.spotify.music",
  "platform": "android",
  "country": "us",
  "lang": "en",
  "fetchedAt": "2026-04-27T11:30:00Z",
  "app": {
    "id": "com.spotify.music",
    "bundleId": null,
    "packageName": "com.spotify.music",
    "title": "Spotify: Music and Podcasts",
    "summary": "Listen to songs, podcasts, and playlists...",
    "description": "Listen to songs, podcasts, and playlists...",
    "developer": {
      "id": "Spotify AB",
      "name": "Spotify AB",
      "website": "https://www.spotify.com/",
      "email": "android-support@spotify.com"
    },
    "iconUrl": "https://...",
    "screenshotUrls": [
      "https://..."
    ],
    "rating": {
      "score": 4.3,
      "count": 31234567,
      "reviewCount": 1234567,
      "histogram": null
    },
    "pricing": {
      "price": 0,
      "currency": "USD",
      "formattedPrice": "Free",
      "free": true
    },
    "categories": [
      {
        "id": "MUSIC_AND_AUDIO",
        "name": "Music & Audio"
      }
    ],
    "contentRating": "Teen",
    "releasedAt": "2014-05-27T00:00:00Z",
    "updatedAt": "2026-04-20T00:00:00Z",
    "version": "9.0.0.0",
    "releaseNotes": "Bug fixes and improvements.",
    "storeUrl": "https://play.google.com/store/apps/details?id=com.spotify.music"
  }
}
```

For iOS, `app.id` should be the numeric App Store ID and `app.bundleId` should be the bundle ID when available. For Android, `app.id` and `app.packageName` should both use the package name. Platform-specific fields that are unavailable should be `null`; lists should be `[]`.

### 3. Fetch raw reviews

User fetches recent or helpful reviews for later LLM analysis.

Example:

```bash
appstore-scraper reviews com.spotify.music --platform android --country us --lang en --sort newest --limit 200
appstore-scraper reviews 324684580 --platform ios --country us --sort newest --limit 200
```

Expected output:

```json
{
  "appId": "com.spotify.music",
  "platform": "android",
  "country": "us",
  "lang": "en",
  "sort": "newest",
  "limit": 200,
  "count": 200,
  "fetchedAt": "2026-04-27T11:30:00Z",
  "reviews": [
    {
      "id": "review-id",
      "user": {
        "name": "Reviewer",
        "imageUrl": "https://...",
        "url": null
      },
      "rating": 4,
      "title": "Great app",
      "text": "Review text...",
      "reviewedAt": "2026-04-27T00:00:00Z",
      "version": "1.2.3",
      "url": "https://...",
      "helpfulCount": 12,
      "developerResponse": {
        "text": "Developer response...",
        "respondedAt": "2026-04-27T00:00:00Z"
      }
    }
  ]
}
```

## CLI Requirements

### Command Structure

The CLI binary should be named `appstore-scraper`.

Required commands:

```bash
appstore-scraper search <term>
appstore-scraper details <app-id>
appstore-scraper reviews <app-id>
```

Common flags:

- `--platform ios|android` required.
- `--country <code>` default `us`.
- `--lang <code>` default `en` where supported.
- `--output json|csv` default `json`.
- `--timeout <duration>` default `30s`.

Search-specific flags:

- `--limit <n>` default `10`, max `250`.

Review-specific flags:

- `--limit <n>` default `100`, max `1000`.
- `--sort newest|rating|helpfulness` default `newest`.
- `ios` supports `newest` and `helpfulness`.
- `android` supports `newest`, `rating`, and `helpfulness`.
- Unsupported platform/sort combinations must return a validation error before making network calls.

### Help Text

The CLI must provide concise help text for the root command and each required subcommand.

Root command:

```text
appstore-scraper retrieves public mobile app store data.

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
  -h, --help              Show help
```

Search command:

```text
Search for apps by keyword.

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
      --timeout duration  Request timeout (default 30s)
```

Details command:

```text
Fetch app metadata by store ID, bundle ID, or package name.

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
      --timeout duration  Request timeout (default 30s)
```

Reviews command:

```text
Fetch raw app reviews.

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
  android  newest, rating, helpfulness

Global Flags:
      --platform string   Store platform: ios or android
      --country string    Two-letter store country code (default "us")
      --lang string       Language code where supported (default "en")
      --output string     Output format: json or csv (default "json")
      --timeout duration  Request timeout (default 30s)
```

### Output Behavior

- JSON output must be deterministic and machine-readable.
- CSV output must be supported for `search`, `details`, and `reviews`.
- Errors must be written to stderr.
- Successful command output must be written to stdout.
- Non-zero exit codes must be used for invalid input, network failure, scraper failure, and unsupported platform behavior.

## Data Model Requirements

### Platform Normalization

The CLI must normalize iOS and Android responses into common shapes while preserving platform-specific fields when useful.

General rules:

- JSON fields must use `camelCase`.
- CSV columns must use `snake_case`.
- Do not wrap command output in generic `data`, `command`, or schema version objects.
- Each command should return a command-native top-level object: `search` returns `results`, `details` returns `app`, and `reviews` returns `reviews`.
- Include execution context at the top level: `platform`, `country`, `lang`, `fetchedAt`, and command-specific inputs such as `query`, `appId`, `sort`, and `limit`.
- All timestamps must be RFC3339 strings in UTC.
- Date-only source values must be normalized to midnight UTC.
- Use `null` for unavailable common scalar fields.
- Use `[]` for unavailable list fields.

Common app fields:

- `id`
- `bundleId`
- `packageName`
- `title`
- `summary`
- `description`
- `developer`
- `iconUrl`
- `screenshotUrls`
- `rating`
- `pricing`
- `categories`
- `contentRating`
- `releasedAt`
- `updatedAt`
- `version`
- `releaseNotes`
- `storeUrl`

### Review Normalization

Common review fields:

- `id`
- `user`
- `rating`
- `title`
- `text`
- `reviewedAt`
- `version`
- `url`
- `helpfulCount`
- `developerResponse`

`developerResponse` should be `null` when unavailable. For iOS, developer responses are expected to be unavailable through public review endpoints.

### Search CSV Export

When `search --output csv` is used, the CLI must write RFC 4180-compatible CSV to stdout.

Required CSV columns:

- `platform`
- `country`
- `lang`
- `query`
- `limit`
- `fetched_at`
- `id`
- `bundle_id`
- `package_name`
- `title`
- `developer_id`
- `developer_name`
- `icon_url`
- `rating_score`
- `rating_count`
- `price`
- `currency`
- `formatted_price`
- `free`
- `store_url`

### Details CSV Export

When `details --output csv` is used, the CLI must write RFC 4180-compatible CSV to stdout.

The details CSV output must contain exactly one data row.

Required CSV columns:

- `platform`
- `country`
- `lang`
- `app_id`
- `fetched_at`
- `id`
- `bundle_id`
- `package_name`
- `title`
- `summary`
- `description`
- `developer_id`
- `developer_name`
- `developer_website`
- `developer_email`
- `icon_url`
- `screenshot_urls`
- `rating_score`
- `rating_count`
- `review_count`
- `price`
- `currency`
- `formatted_price`
- `free`
- `categories`
- `content_rating`
- `released_at`
- `updated_at`
- `version`
- `release_notes`
- `store_url`

List fields such as `screenshot_urls` and `categories` must be joined with `|`.

### Review CSV Export

When `reviews --output csv` is used, the CLI must write RFC 4180-compatible CSV to stdout.

Required CSV columns:

- `platform`
- `country`
- `lang`
- `app_id`
- `sort`
- `limit`
- `fetched_at`
- `review_id`
- `user_name`
- `rating`
- `title`
- `text`
- `reviewed_at`
- `version`
- `url`
- `helpful_count`
- `developer_response_text`
- `developer_responded_at`

Rules:

- Always include the header row.
- Escape commas, quotes, and newlines correctly using Go's standard CSV writer.
- Use empty strings for unavailable platform-specific fields.
- Preserve review text newlines unless a future `--single-line` option is added.
- Do not include nested JSON blobs in CSV cells.

## Store Support

### Apple App Store

Required:

- Search by term.
- Details by numeric app ID.
- Details by bundle ID when feasible.
- Reviews by numeric app ID.

Preferred:

- Resolve bundle ID to numeric ID before fetching reviews.

Known limitations:

- Review pagination is limited by public App Store endpoints.
- Developer responses are not available through the public review endpoint.

### Google Play Store

Required:

- Search by term.
- Details by package name.
- Reviews by package name.
- Include developer responses when available.

Known limitations:

- Google Play data is scraped from public endpoints and may change without notice.
- Sorting and pagination behavior may differ from public store UI behavior.

## Technical Requirements

- Language: Go.
- Minimum Go version: Go 1.22.
- Use a maintainable CLI framework.
- Keep scraping logic separate from CLI command wiring.
- Use context-aware HTTP requests with timeouts.
- Use typed response structs for public output.
- Use table-driven tests for normalization logic.
- Add a `README.md` that documents installation, command usage, flags, output formats, examples, limitations, and only functionality that is actually implemented.

## Reliability Requirements

- Validate required flags before making network calls.
- Return useful error messages for app not found, unsupported platform, and network errors.
- Enforce maximum limits for search and review fetching.
- Add retry support only for transient failures, with conservative defaults.
- Avoid aggressive scraping behavior that could trigger store blocking.

## Testing Requirements

Unit tests:

- CLI argument validation.
- Platform normalization.
- JSON output schema.
- CSV output schemas and escaping.
- Error formatting.

Integration tests:

- Search one known iOS app.
- Search one known Android app.
- Fetch details for one known iOS app.
- Fetch details for one known Android app.
- Fetch a small number of reviews for one known app per platform.

Integration tests should be optional or tagged because public store data can change.

## Milestones

### Milestone 1: Project Scaffold

- Initialize Go module.
- Add CLI framework.
- Add root command and help text.
- Add JSON output helper.
- Add CI for lint/test/build.

### Milestone 2: Search

- Implement iOS search.
- Implement Android search.
- Normalize search output.
- Add tests for command validation and output formatting.

### Milestone 3: Details

- Implement iOS details.
- Implement Android details.
- Normalize metadata output.
- Add fixture-based tests.

### Milestone 4: Reviews

- Implement Android reviews with developer response fields.
- Implement iOS reviews with pagination limit.
- Normalize review output.
- Add CSV export for reviews.
- Add limit and sort handling.

### Milestone 5: Release

- Add `README.md` with installation instructions, examples, and documented limitations.
- Add GitHub Actions release workflow.
- Publish v0.1.0.

## Success Metrics

- A user can fetch app reviews from both stores with one command.
- JSON output can be piped directly into an LLM review analysis script.
- CSV output can be imported into spreadsheets, databases, and batch labeling workflows.
- Search, details, and reviews work for common public apps on both platforms.
- The README accurately describes only implemented functionality.
- No mock ASO scoring or unimplemented keyword research appears in the CLI.
