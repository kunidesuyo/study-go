package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"go-api-arch-clean-template/adapter/controller/gin/presenter"
	"go-api-arch-clean-template/entity"
	"go-api-arch-clean-template/pkg"
)

type MockAlbumUseCase struct {
	mock.Mock
}

func NewMockAlbumUseCase() *MockAlbumUseCase {
	return &MockAlbumUseCase{}
}

func (m *MockAlbumUseCase) Create(album *entity.Album) (*entity.Album, error) {
	args := m.Called(album)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Album), args.Error(1)
}

func (m *MockAlbumUseCase) Get(ID int) (*entity.Album, error) {
	args := m.Called(ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Album), args.Error(1)
}
func (m *MockAlbumUseCase) Save(album *entity.Album) (*entity.Album, error) {
	args := m.Called(album)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Album), args.Error(1)
}

func (m *MockAlbumUseCase) Delete(ID int) error {
	args := m.Called(ID)
	if args.Get(0) == nil {
		return args.Error(1)
	}
	return nil
}

type AlbumHandlersSuite struct {
	suite.Suite
	albumHandler *AlbumHandler
}

func TestAlbumHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(AlbumHandlersSuite))
}

func (suite *AlbumHandlersSuite) TestCreate() {
	now := pkg.Str2time("2023-01-01")
	mockUseCase := NewMockAlbumUseCase()
	album := &entity.Album{
		Title:       "album",
		ReleaseDate: now,
		Category:    entity.Category{Name: "sports"},
	}

	mockUseCase.On("Create", album).Return(&entity.Album{
		ID:          1,
		Title:       "album",
		ReleaseDate: now,
		CategoryID:  1,
		Category: entity.Category{
			ID:   1,
			Name: "sports",
		},
	}, nil)
	suite.albumHandler = NewAlbumHandler(mockUseCase)

	request, _ := presenter.NewCreateAlbumRequest("/api/v1", presenter.CreateAlbumJSONRequestBody{
		Title:       "album",
		Category:    presenter.Category{Name: "sports"},
		ReleaseDate: presenter.ReleaseDate{Time: now},
	})
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.albumHandler.CreateAlbum(ginContext)

	suite.Assert().Equal(http.StatusCreated, w.Code)
	bodyBytes, _ := io.ReadAll(w.Body)
	var albumGetResponse presenter.AlbumResponse
	err := json.Unmarshal(bodyBytes, &albumGetResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusCreated, w.Code)
	suite.Assert().Equal("album", albumGetResponse.Title)
	suite.Assert().Equal("sports", string(albumGetResponse.Category.Name))
	suite.Assert().NotNil(albumGetResponse.ReleaseDate)
}

func (suite *AlbumHandlersSuite) TestCreateRequestBodyFailure() {
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("POST", "/api/v1/album", nil)
	req.Header.Add("Content-Type", "application/json")
	ginContext.Request = req

	suite.albumHandler.CreateAlbum(ginContext)
	suite.Assert().Equal(http.StatusBadRequest, w.Code)
	suite.Assert().JSONEq(`{"message": "invalid request"}`, w.Body.String())
}

func (suite *AlbumHandlersSuite) TestCreateFailure() {
	now := pkg.Str2time("2023-01-01")
	mockUseCase := NewMockAlbumUseCase()
	album := &entity.Album{
		Title:       "album",
		ReleaseDate: now,
		Category:    entity.Category{Name: "sports"},
	}

	mockUseCase.On("Create", album).Return(nil,
		errors.New("invalid"))
	suite.albumHandler = NewAlbumHandler(mockUseCase)

	request, _ := presenter.NewCreateAlbumRequest("/api/v1", presenter.CreateAlbumJSONRequestBody{
		Title:       "album",
		Category:    presenter.Category{Name: "sports"},
		ReleaseDate: presenter.ReleaseDate{Time: now},
	})
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.albumHandler.CreateAlbum(ginContext)

	suite.Assert().Equal(http.StatusInternalServerError, w.Code)
	suite.Assert().JSONEq(`{"message":"invalid"}`, w.Body.String())
}

func (suite *AlbumHandlersSuite) TestGet() {
	now := pkg.Str2time("2023-01-01")
	mockUseCase := NewMockAlbumUseCase()
	mockUseCase.On("Get", 1).Return(&entity.Album{
		ID:          1,
		Title:       "album",
		ReleaseDate: now,
		CategoryID:  1,
		Category: entity.Category{
			ID:   1,
			Name: "sports",
		},
	}, nil)
	suite.albumHandler = NewAlbumHandler(mockUseCase)

	request, _ := presenter.NewGetAlbumByIdRequest("/api/v1", 1)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.albumHandler.GetAlbumById(ginContext, 1)

	bodyBytes, _ := io.ReadAll(w.Body)
	var albumGetResponse presenter.AlbumResponse
	err := json.Unmarshal(bodyBytes, &albumGetResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal("album", albumGetResponse.Title)
	suite.Assert().Equal("sports", string(albumGetResponse.Category.Name))
	suite.Assert().NotNil(albumGetResponse.ReleaseDate)
}

func (suite *AlbumHandlersSuite) TestGetNoAlbumFailure() {
	mockUseCase := NewMockAlbumUseCase()
	mockUseCase.On("Get", 1111).Return(nil,
		errors.New("invalid"))
	suite.albumHandler = NewAlbumHandler(mockUseCase)

	request, _ := presenter.NewGetAlbumByIdRequest("/api/v1", 1111)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.albumHandler.GetAlbumById(ginContext, 1111)

	suite.Assert().Equal(http.StatusInternalServerError, w.Code)
	suite.Assert().JSONEq(`{"message": "invalid"}`, w.Body.String())
}

func (suite *AlbumHandlersSuite) TestUpdate() {

	now := pkg.Str2time("2023-01-01")
	mockUseCase := NewMockAlbumUseCase()
	title := "updated"
	categoryName := "food"

	album := &entity.Album{
		ID:       1,
		Title:    title,
		Category: entity.Category{Name: entity.CategoryName(categoryName)},
	}

	mockUseCase.On("Save", album).Return(&entity.Album{
		ID:          1,
		Title:       title,
		ReleaseDate: now,
		CategoryID:  1,
		Category: entity.Category{
			ID:   1,
			Name: entity.CategoryName(categoryName),
		},
	}, nil)

	suite.albumHandler = NewAlbumHandler(mockUseCase)

	category := presenter.Category{
		Name: presenter.CategoryName(categoryName),
	}
	request, _ := presenter.NewUpdateAlbumByIdRequest("/api/v1", 1,
		presenter.UpdateAlbumByIdJSONRequestBody{
			Title:    &title,
			Category: &category,
		},
	)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.albumHandler.UpdateAlbumById(ginContext, 1)

	bodyBytes, _ := io.ReadAll(w.Body)
	var albumGetResponse presenter.AlbumResponse
	err := json.Unmarshal(bodyBytes, &albumGetResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal("updated", albumGetResponse.Title)
	suite.Assert().Equal("food", string(albumGetResponse.Category.Name))
	suite.Assert().NotNil(albumGetResponse.ReleaseDate)
}

func (suite *AlbumHandlersSuite) TestUpdateRequestBodyFailure() {
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("PATCH", "/api/v1/album", nil)
	req.Header.Add("Content-Type", "application/json")
	ginContext.Request = req

	suite.albumHandler.CreateAlbum(ginContext)

	suite.Assert().Equal(http.StatusBadRequest, w.Code)
	suite.Assert().JSONEq(`{"message": "invalid request"}`, w.Body.String())
}

func (suite *AlbumHandlersSuite) TestUpdateAlbumFailure() {
	mockUseCase := NewMockAlbumUseCase()
	title := "updated"
	categoryName := "food"
	album := &entity.Album{
		ID:       1111,
		Title:    title,
		Category: entity.Category{Name: entity.CategoryName(categoryName)},
	}

	mockUseCase.On("Save", album).Return(nil,
		errors.New("invalid"))
	suite.albumHandler = NewAlbumHandler(mockUseCase)

	category := presenter.Category{
		Name: presenter.CategoryName(categoryName),
	}
	request, _ := presenter.NewUpdateAlbumByIdRequest("/api/v1", 1111,
		presenter.UpdateAlbumByIdJSONRequestBody{
			Title:    &title,
			Category: &category,
		},
	)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.albumHandler.UpdateAlbumById(ginContext, 1111)

	suite.Assert().Equal(http.StatusInternalServerError, w.Code)
	suite.Assert().JSONEq(`{"message": "invalid"}`, w.Body.String())
}

func (suite *AlbumHandlersSuite) TestDelete() {
	mockUseCase := NewMockAlbumUseCase()
	mockUseCase.On("Delete", 1).Return(nil, nil)
	suite.albumHandler = NewAlbumHandler(mockUseCase)

	request, _ := presenter.NewDeleteAlbumByIdRequest("/api/v1", 1)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.albumHandler.DeleteAlbumById(ginContext, 1)

	suite.Assert().Equal(http.StatusNoContent, w.Code)
}

func (suite *AlbumHandlersSuite) TestDeleteAlbumFailure() {
	mockUseCase := NewMockAlbumUseCase()
	mockUseCase.On("Delete", 1111).Return(nil, errors.New("invalid"))
	suite.albumHandler = NewAlbumHandler(mockUseCase)

	request, _ := presenter.NewDeleteAlbumByIdRequest("/api/v1", 1111)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.albumHandler.DeleteAlbumById(ginContext, 1111)

	suite.Assert().Equal(http.StatusInternalServerError, w.Code)
	suite.Assert().JSONEq(`{"message": "invalid"}`, w.Body.String())
}
