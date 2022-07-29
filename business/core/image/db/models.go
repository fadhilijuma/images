package db

import "time"

// Image represents an individual image.
type Image struct {
	ID           string    `db:"image_id"`      // Unique identifier.
	ImageURL     string    `db:"image_url"`     // Display image URL to the path on file system containing image.
	UserID       string    `db:"user_id"`       // ID of the user who created the image.
	DateUploaded time.Time `db:"date_uploaded"` // When the image was uploaded.
}
