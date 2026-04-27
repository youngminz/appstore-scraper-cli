package model

import "time"

type Developer struct {
	ID      *string `json:"id"`
	Name    *string `json:"name"`
	Website *string `json:"website"`
	Email   *string `json:"email"`
}

type Rating struct {
	Score       *float64      `json:"score"`
	Count       *int64        `json:"count"`
	ReviewCount *int64        `json:"reviewCount"`
	Histogram   map[int]int64 `json:"histogram"`
}

type Pricing struct {
	Price          *float64 `json:"price"`
	Currency       *string  `json:"currency"`
	FormattedPrice *string  `json:"formattedPrice"`
	Free           *bool    `json:"free"`
}

type Category struct {
	ID   *string `json:"id"`
	Name *string `json:"name"`
}

type App struct {
	ID             string     `json:"id"`
	BundleID       *string    `json:"bundleId"`
	PackageName    *string    `json:"packageName"`
	Title          *string    `json:"title"`
	Summary        *string    `json:"summary"`
	Description    *string    `json:"description"`
	Developer      Developer  `json:"developer"`
	IconURL        *string    `json:"iconUrl"`
	ScreenshotURLs []string   `json:"screenshotUrls"`
	Rating         Rating     `json:"rating"`
	Pricing        Pricing    `json:"pricing"`
	Categories     []Category `json:"categories"`
	ContentRating  *string    `json:"contentRating"`
	ReleasedAt     *time.Time `json:"releasedAt"`
	UpdatedAt      *time.Time `json:"updatedAt"`
	Version        *string    `json:"version"`
	ReleaseNotes   *string    `json:"releaseNotes"`
	StoreURL       *string    `json:"storeUrl"`
}

type SearchResponse struct {
	Query     string    `json:"query"`
	Platform  string    `json:"platform"`
	Country   string    `json:"country"`
	Lang      string    `json:"lang"`
	Limit     int       `json:"limit"`
	Count     int       `json:"count"`
	FetchedAt time.Time `json:"fetchedAt"`
	Results   []App     `json:"results"`
}

type DetailsResponse struct {
	AppID     string    `json:"appId"`
	Platform  string    `json:"platform"`
	Country   string    `json:"country"`
	Lang      string    `json:"lang"`
	FetchedAt time.Time `json:"fetchedAt"`
	App       App       `json:"app"`
}

type User struct {
	Name     *string `json:"name"`
	ImageURL *string `json:"imageUrl"`
	URL      *string `json:"url"`
}

type DeveloperResponse struct {
	Text        *string    `json:"text"`
	RespondedAt *time.Time `json:"respondedAt"`
}

type Review struct {
	ID                *string            `json:"id"`
	User              User               `json:"user"`
	Rating            *int               `json:"rating"`
	Title             *string            `json:"title"`
	Text              *string            `json:"text"`
	ReviewedAt        *time.Time         `json:"reviewedAt"`
	Version           *string            `json:"version"`
	URL               *string            `json:"url"`
	HelpfulCount      *int64             `json:"helpfulCount"`
	DeveloperResponse *DeveloperResponse `json:"developerResponse"`
}

type ReviewsResponse struct {
	AppID     string    `json:"appId"`
	Platform  string    `json:"platform"`
	Country   string    `json:"country"`
	Lang      string    `json:"lang"`
	Sort      string    `json:"sort"`
	Limit     int       `json:"limit"`
	Count     int       `json:"count"`
	FetchedAt time.Time `json:"fetchedAt"`
	Reviews   []Review  `json:"reviews"`
}

func StringPtr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func FloatPtr(v float64) *float64 {
	return &v
}

func Int64Ptr(v int64) *int64 {
	return &v
}

func IntPtr(v int) *int {
	return &v
}

func BoolPtr(v bool) *bool {
	return &v
}
