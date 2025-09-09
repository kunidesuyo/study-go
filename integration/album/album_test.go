package integration

import (
	"context"
	"net/http"
	"testing"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/suite"

	"go-api-arch-clean-template/adapter/controller/gin/presenter"
	"go-api-arch-clean-template/pkg"
)

type AlbumTestSuite struct {
	suite.Suite
}

func TestAlbumSuite(t *testing.T) {
	suite.Run(t, new(AlbumTestSuite))
}

func (suite *AlbumTestSuite) TestAlbumCreateGetDelete() {
	// Create
	baseEndpoint := pkg.GetEndpoint("/api/v1")
	apiClient, _ := presenter.NewClientWithResponses(baseEndpoint)
	createResponse, err := apiClient.CreateAlbumWithResponse(context.Background(), presenter.CreateAlbumJSONRequestBody{
		Title:       "test",
		Category:    presenter.Category{Name: presenter.Sports},
		ReleaseDate: openapi_types.Date{Time: time.Now()},
	})
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusCreated, createResponse.StatusCode())
	suite.Assert().Nil(err)
	suite.Assert().NotNil(createResponse.JSON201.Id)
	suite.Assert().Equal("test", createResponse.JSON201.Title)
	suite.Assert().Equal("sports", string(createResponse.JSON201.Category.Name))
	suite.Assert().NotNil(createResponse.JSON201.ReleaseDate)

	// Get
	getResponse, err := apiClient.GetAlbumByIdWithResponse(context.Background(), createResponse.JSON201.Id)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusOK, getResponse.StatusCode())
	suite.Assert().Nil(err)
	suite.Assert().Equal(createResponse.JSON201.Id, getResponse.JSON200.Id)
	suite.Assert().Equal("test", getResponse.JSON200.Title)
	suite.Assert().Equal("sports", string(getResponse.JSON200.Category.Name))
	suite.Assert().NotNil(getResponse.JSON200.ReleaseDate)

	// Update
	title := "updated"
	category := presenter.Category{
		Name: presenter.Food,
	}
	updateResponse, err := apiClient.UpdateAlbumByIdWithResponse(context.Background(), getResponse.JSON200.Id, presenter.UpdateAlbumByIdJSONRequestBody{
		Title:    &title,
		Category: &category,
	})
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusOK, updateResponse.StatusCode())
	suite.Assert().Nil(err)
	suite.Assert().Equal("updated", updateResponse.JSON200.Title)
	suite.Assert().Equal("food", string(updateResponse.JSON200.Category.Name))
	suite.Assert().NotNil(updateResponse.JSON200.ReleaseDate)

	// Delete
	deleteResponse, err := apiClient.DeleteAlbumByIdWithResponse(context.Background(), updateResponse.JSON200.Id)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusNoContent, deleteResponse.StatusCode())
}
