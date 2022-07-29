package image

import (
	"github.com/fadhilijuma/images/business/core/image/db"
	"time"
)

// Image represents an individual image.
type Image struct {
	ID           string    `json:"id"`            // Unique identifier.
	ImageURL     string    `json:"image_url"`     // Display image url of the image.
	UserID       string    `json:"user_id"`       // User who uploaded the image.
	DateUploaded time.Time `json:"date_uploaded"` // When the image was added.
}

// NewImage is what we require from clients when adding an image.
type NewImage struct {
	ImageURL string `json:"image_url" validate:"required"`
	UserID   string `json:"user_id" validate:"required"`
}

// UpdateImage defines what information may be provided to modify an
// existing Image. All fields are optional so clients can send just the
// fields they want changed. It uses pointer fields so we can differentiate
// between a field that was not provided and a field that was provided as
// explicitly blank. Normally we do not want to use pointers to basic types but
// we make exceptions around marshalling/unmarshalling.
type UpdateImage struct {
	ImageURL *string `json:"image_url"`
	UserID   *string `json:"user_id"`
}

// =============================================================================

func toImage(dbPrd db.Image) Image {
	pu := (*Image)(&dbPrd)
	return *pu
}

func toImageSlice(dbImages []db.Image) []Image {
	images := make([]Image, len(dbImages))
	for i, dbImage := range dbImages {
		images[i] = toImage(dbImage)
	}
	return images
}
