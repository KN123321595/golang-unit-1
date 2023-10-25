package model

type ApodMetadata struct {
	ID             int    `db:"id" json:"id"`
	Explanation    string `db:"explanation" json:"explanation"`
	MediaType      string `db:"media_type" json:"media_type"`
	ServiceVersion string `db:"service_version" json:"service_version"`
	Title          string `db:"title" json:"title"`
	Hdurl          string `db:"hdurl" json:"hdurl"`
	ThumbnailUrl   string `db:"thumbnail_url" json:"thumbnail_url"`
	Url            string `db:"url" json:"url"`
	ImagePath      string `db:"image_path" json:"image_path"`
	Date           string `db:"date" json:"date"`
}
