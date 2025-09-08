package gateway_test

import (
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"

	"go-api-arch-clean-template/adapter/gateway"
	"go-api-arch-clean-template/entity"
	"go-api-arch-clean-template/pkg"
	"go-api-arch-clean-template/pkg/tester"
)

type AlbumRepositorySuite struct {
	tester.DBSQLiteSuite
	repository gateway.AlbumRepository
}

func TestAlbumRepositorySuite(t *testing.T) {
	suite.Run(t, new(AlbumRepositorySuite))
}

func (suite *AlbumRepositorySuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	suite.repository = gateway.NewAlbumRepository(suite.DB)
}

func (suite *AlbumRepositorySuite) MockDB() sqlmock.Sqlmock {
	mock, mockGormDB := tester.MockDB()
	suite.repository = gateway.NewAlbumRepository(mockGormDB)
	return mock
}

func (suite *AlbumRepositorySuite) AfterTest(suiteName, testName string) {
	suite.repository = gateway.NewAlbumRepository(suite.DB)
}

func (suite *AlbumRepositorySuite) TestAlbumRepositoryCRUD() {
	now := pkg.Str2time("2023-01-01")
	album := &entity.Album{
		Title:       "test",
		ReleaseDate: now,
		Category:    entity.Category{Name: entity.CategoryName("sports")},
	}
	album, err := suite.repository.Create(album)
	suite.Assert().Nil(err)
	suite.Assert().NotZero(album.ID)
	suite.Assert().Equal("test", album.Title)
	suite.Assert().Equal(now, album.ReleaseDate)
	suite.Assert().NotZero(album.Category.ID)
	suite.Assert().Equal("sports", string(album.Category.Name))

	getAlbum, err := suite.repository.Get(album.ID)
	suite.Assert().Nil(err)
	suite.Assert().Equal("test", getAlbum.Title)
	suite.Assert().Equal(now, album.ReleaseDate)
	suite.Assert().Equal(album.Category.ID, getAlbum.Category.ID)
	suite.Assert().Equal("sports", string(getAlbum.Category.Name))

	getAlbum.Title = "updated"
	updatedAlbum, err := suite.repository.Save(getAlbum)
	suite.Assert().Nil(err)
	suite.Assert().Equal("updated", updatedAlbum.Title)
	suite.Assert().NotNil(updatedAlbum.ReleaseDate)
	suite.Assert().NotNil(updatedAlbum.Category.ID)
	suite.Assert().Equal("sports", string(updatedAlbum.Category.Name))

	err = suite.repository.Delete(updatedAlbum.ID)
	suite.Assert().Nil(err)
	deletedAlbum, err := suite.repository.Get(updatedAlbum.ID)
	suite.Assert().Nil(deletedAlbum)
	suite.Assert().True(strings.Contains("record not found", err.Error()))
}

func (suite *AlbumRepositorySuite) TestAlbumCreateFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE `categories`.`name` = ? ORDER BY `categories`.`id` LIMIT ?")).WithArgs("sports", 1).WillReturnError(errors.New("create error"))

	album := &entity.Album{
		Title:       "test",
		ReleaseDate: time.Now(),
		Category:    entity.Category{Name: entity.CategoryName("sports")},
	}

	createdAlbum, err := suite.repository.Create(album)
	suite.Assert().Nil(createdAlbum)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("create error", err.Error())
}

func (suite *AlbumRepositorySuite) TestAlbumGetFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `albums` WHERE `albums`.`id` = ? ORDER BY `albums`.`id` LIMIT ?")).WithArgs(1, 1).WillReturnError(errors.New("get error"))

	album, err := suite.repository.Get(1)
	suite.Assert().Nil(album)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("get error", err.Error())
}

func (suite *AlbumRepositorySuite) TestAlbumDeleteFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta("DELETE FROM `albums` WHERE id = ? AND `albums`.`id` = ?")).WithArgs(1, 1).WillReturnError(errors.New("delete error"))
	mockDB.ExpectRollback()
	mockDB.ExpectCommit()

	err := suite.repository.Delete(1)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("delete error", err.Error())
}

func (suite *AlbumRepositorySuite) TestAlbumSaveFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `albums` WHERE `albums`.`id` = ? ORDER BY `albums`.`id` LIMIT ?")).WithArgs(1, 1).WillReturnError(errors.New("save error"))

	album := &entity.Album{
		ID:       1,
		Title:    "test",
		Category: entity.Category{Name: entity.CategoryName("sports")},
	}

	album, err := suite.repository.Save(album)
	suite.Assert().Nil(album)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("save error", err.Error())
}
