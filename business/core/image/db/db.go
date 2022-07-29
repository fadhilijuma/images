// Package db contains image related CRUD functionality.
package db

import (
	"context"
	"fmt"

	"github.com/fadhilijuma/images/business/sys/database"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Store manages the set of APIs for user access.
type Store struct {
	log          *zap.SugaredLogger
	tr           database.Transactor
	db           sqlx.ExtContext
	isWithinTran bool
}

// NewStore constructs a data for api access.
func NewStore(log *zap.SugaredLogger, db *sqlx.DB) Store {
	return Store{
		log: log,
		tr:  db,
		db:  db,
	}
}

// WithinTran runs passed function and do commit/rollback at the end.
func (s Store) WithinTran(ctx context.Context, fn func(sqlx.ExtContext) error) error {
	if s.isWithinTran {
		fn(s.db)
	}
	return database.WithinTran(ctx, s.log, s.tr, fn)
}

// Tran return new Store with transaction in it.
func (s Store) Tran(tx sqlx.ExtContext) Store {
	return Store{
		log:          s.log,
		tr:           s.tr,
		db:           tx,
		isWithinTran: true,
	}
}

// Create adds an Image to the database. It returns the created Image with
// fields like ID and DateUploaded populated.
func (s Store) Create(ctx context.Context, image Image) error {
	const q = `
	INSERT INTO images
		(image_id, image_url, user_id, date_uploaded)
	VALUES
		(:image_id, :image_url, :user_id, :date_uploaded)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, image); err != nil {
		return fmt.Errorf("inserting image: %w", err)
	}

	return nil
}

// Update modifies metadata about an Image. It will error if the specified ID is
// invalid or does not reference an existing Image.
func (s Store) Update(ctx context.Context, image Image) error {
	const q = `
	UPDATE
		images
	SET
		"image_url" = :image_url,
		"user_id" = :user_id
	WHERE
		image_id = :image_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, image); err != nil {
		return fmt.Errorf("updating image imageID[%s]: %w", image.ID, err)
	}

	return nil
}

// Delete removes the Image identified by a given ID.
func (s Store) Delete(ctx context.Context, imageID string) error {
	data := struct {
		ImageID string `db:"image_id"`
	}{
		ImageID: imageID,
	}

	const q = `
	DELETE FROM
		images
	WHERE
		image_id = :image_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("deleting image imageID[%s]: %w", imageID, err)
	}

	return nil
}

// Query gets all Images from the database.
func (s Store) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]Image, error) {
	data := struct {
		Offset      int `db:"offset"`
		RowsPerPage int `db:"rows_per_page"`
	}{
		Offset:      (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
	}

	const q = `
	SELECT
		*
	FROM
		Images
	ORDER BY
		user_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var images []Image
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &images); err != nil {
		return nil, fmt.Errorf("selecting products: %w", err)
	}

	return images, nil
}

// QueryByID finds the Image identified by a given ID.
func (s Store) QueryByID(ctx context.Context, imageID string) (Image, error) {
	data := struct {
		ImageID string `db:"image_id"`
	}{
		ImageID: imageID,
	}

	const q = `
	SELECT
		*
	FROM
		images 
	GROUP BY
		image_id`

	var prd Image
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &prd); err != nil {
		return Image{}, fmt.Errorf("selecting image imageID[%q]: %w", imageID, err)
	}

	return prd, nil
}

// QueryByUserID finds the Image identified by a given User ID.
func (s Store) QueryByUserID(ctx context.Context, userID string) ([]Image, error) {
	data := struct {
		UserID string `db:"user_id"`
	}{
		UserID: userID,
	}

	const q = `
	SELECT
		*
	FROM
		images
	GROUP BY
		user_id`

	var prds []Image
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &prds); err != nil {
		return nil, fmt.Errorf("selecting products userID[%s]: %w", userID, err)
	}

	return prds, nil
}
