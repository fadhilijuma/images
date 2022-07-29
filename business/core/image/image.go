// Package image provides an example of a core business API. Right now these
// calls are just wrapping the data/store layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package image

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fadhilijuma/images/business/core/image/db"
	"github.com/fadhilijuma/images/business/sys/database"
	"github.com/fadhilijuma/images/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound  = errors.New("image not found")
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// Core manages the set of APIs for image access.
type Core struct {
	store db.Store
}

// NewCore constructs a core for image api access.
func NewCore(log *zap.SugaredLogger, sqlxDB *sqlx.DB) Core {
	return Core{
		store: db.NewStore(log, sqlxDB),
	}
}

// Create adds an Image to the database. It returns the created Image with
// fields like ID and DateCreated populated.
func (c Core) Create(ctx context.Context, np NewImage, now time.Time) (Image, error) {
	if err := validate.Check(np); err != nil {
		return Image{}, fmt.Errorf("validating data: %w", err)
	}

	dbImg := db.Image{
		ID:           validate.GenerateID(),
		ImageURL:     np.ImageURL,
		UserID:       np.UserID,
		DateUploaded: now,
	}

	if err := c.store.Create(ctx, dbImg); err != nil {
		return Image{}, fmt.Errorf("create: %w", err)
	}

	return toImage(dbImg), nil
}

// Update modifies data about a Product. It will error if the specified ID is
// invalid or does not reference an existing Product.
func (c Core) Update(ctx context.Context, productID string, up UpdateImage, now time.Time) error {
	if err := validate.CheckID(productID); err != nil {
		return ErrInvalidID
	}

	if err := validate.Check(up); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	dbImg, err := c.store.QueryByID(ctx, productID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating image productID[%s]: %w", productID, err)
	}

	if up.ImageURL != nil {
		dbImg.ImageURL = *up.ImageURL
	}
	if up.UserID != nil {
		dbImg.UserID = *up.UserID
	}
	dbImg.DateUploaded = now

	if err := c.store.Update(ctx, dbImg); err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil
}

// Delete removes the image identified by a given ID.
func (c Core) Delete(ctx context.Context, productID string) error {
	if err := validate.CheckID(productID); err != nil {
		return ErrInvalidID
	}

	if err := c.store.Delete(ctx, productID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query gets all Products from the database.
func (c Core) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]Image, error) {
	dbImg, err := c.store.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return toImageSlice(dbImg), nil
}

// QueryByID finds the image identified by a given ID.
func (c Core) QueryByID(ctx context.Context, productID string) (Image, error) {
	if err := validate.CheckID(productID); err != nil {
		return Image{}, ErrInvalidID
	}

	dbImg, err := c.store.QueryByID(ctx, productID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Image{}, ErrNotFound
		}
		return Image{}, fmt.Errorf("query: %w", err)
	}

	return toImage(dbImg), nil
}

// QueryByUserID finds the products identified by a given User ID.
func (c Core) QueryByUserID(ctx context.Context, userID string) ([]Image, error) {
	if err := validate.CheckID(userID); err != nil {
		return nil, ErrInvalidID
	}

	dbImg, err := c.store.QueryByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return toImageSlice(dbImg), nil
}
