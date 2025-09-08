package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-api-arch-clean-template/entity"
	"go-api-arch-clean-template/pkg"
	"go-api-arch-clean-template/pkg/tester"
)

func TestAlbum(t *testing.T) {
	category := entity.Category{
		ID:   1,
		Name: "sports",
	}

	now := pkg.Str2time("2023-01-01")
	mockClock := tester.NewMockClock(now)
	album := entity.Album{
		ID:          1,
		Title:       "album",
		ReleaseDate: now,
		CategoryID:  1,
		Category:    category,
	}
	assert.Equal(t, 1, album.ID)
	assert.Equal(t, 0, album.Anniversary(mockClock))
	assert.Equal(t, "album", album.Title)
	assert.Equal(t, now, album.ReleaseDate)
	assert.Equal(t, 1, album.CategoryID)
	assert.Equal(t, 1, album.Category.ID)
	assert.Equal(t, "sports", string(album.Category.Name))
}

func TestAlbumAnniversary(t *testing.T) {
	mockedClock := tester.NewMockClock(pkg.Str2time("2022-04-01"))

	// non-leap
	album := entity.Album{ReleaseDate: pkg.Str2time("2022-04-01")}
	assert.Equal(t, 0, album.Anniversary(mockedClock))
	album = entity.Album{ReleaseDate: pkg.Str2time("2021-04-02")}
	assert.Equal(t, 0, album.Anniversary(mockedClock))
	album = entity.Album{ReleaseDate: pkg.Str2time("2021-04-01")}
	assert.Equal(t, 1, album.Anniversary(mockedClock))
	// leap
	album = entity.Album{ReleaseDate: pkg.Str2time("2020-04-02")}
	assert.Equal(t, 1, album.Anniversary(mockedClock))
	album = entity.Album{ReleaseDate: pkg.Str2time("2020-04-01")}
	assert.Equal(t, 2, album.Anniversary(mockedClock))
}
