package image_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/fadhilijuma/images/business/core/image"
	"github.com/fadhilijuma/images/business/data/dbtest"
	"github.com/fadhilijuma/images/foundation/docker"
	"github.com/google/go-cmp/cmp"
)

var c *docker.Container

func TestMain(m *testing.M) {
	var err error
	c, err = dbtest.StartDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbtest.StopDB(c)

	m.Run()
}

func Test_Image(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, c, "testprod")
	t.Cleanup(teardown)

	core := image.NewCore(log, db)

	t.Log("Given the need to work with Image records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single Image.", testID)
		{
			ctx := context.Background()
			now := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)

			ni := image.NewImage{
				ImageURL: "images/image.jpg",
				UserID:   "123",
			}

			img, err := core.Create(ctx, ni, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create a image : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create a image.", dbtest.Success, testID)

			saved, err := core.QueryByID(ctx, img.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve image by ID: %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve image by ID.", dbtest.Success, testID)

			if diff := cmp.Diff(img, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same image. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same image.", dbtest.Success, testID)

			upd := image.UpdateImage{
				ImageURL: dbtest.StringPointer("images/pic.jpg"),
				UserID:   dbtest.StringPointer("12345"),
			}
			updatedTime := time.Date(2019, time.January, 1, 1, 1, 1, 0, time.UTC)

			if err := core.Update(ctx, img.ID, upd, updatedTime); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update image : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update image.", dbtest.Success, testID)

			products, err := core.Query(ctx, 1, 3)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve updated image : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve updated image.", dbtest.Success, testID)

			// Check specified fields were updated. Make a copy of the original image
			// and change just the fields we expect then diff it with what was saved.
			want := img
			want.ImageURL = *upd.ImageURL
			want.UserID = *upd.UserID

			var idx int
			for i, p := range products {
				if p.ID == want.ID {
					idx = i
				}
			}
			if diff := cmp.Diff(want, products[idx]); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same image. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same image.", dbtest.Success, testID)

			upd = image.UpdateImage{
				ImageURL: dbtest.StringPointer("images/image.jpg"),
				UserID:   dbtest.StringPointer("1234"),
			}

			if err := core.Update(ctx, img.ID, upd, updatedTime); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update just some fields of image : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update just some fields of image.", dbtest.Success, testID)

			saved, err = core.QueryByID(ctx, img.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve updated image : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve updated image.", dbtest.Success, testID)

			if saved.ImageURL != *upd.ImageURL {
				t.Fatalf("\t%s\tTest %d:\tShould be able to see updated ImageURL field : got %q want %q.", dbtest.Failed, testID, saved.ImageURL, *upd.ImageURL)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see updated Name field.", dbtest.Success, testID)
			}

			if err := core.Delete(ctx, img.ID); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete image : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete image.", dbtest.Success, testID)

			_, err = core.QueryByID(ctx, img.ID)
			if !errors.Is(err, image.ErrNotFound) {
				t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve deleted image : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve deleted image.", dbtest.Success, testID)
		}
	}
}
