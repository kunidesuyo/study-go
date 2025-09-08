package gateway

import (
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"go-api-arch-clean-template/entity"
)

type AlbumRepository interface {
	Create(album *entity.Album) (*entity.Album, error)
	Get(ID int) (*entity.Album, error)
	Save(*entity.Album) (*entity.Album, error)
	Delete(ID int) error
}

type albumRepository struct {
	db *gorm.DB
}

func NewAlbumRepository(db *gorm.DB) AlbumRepository {
	return &albumRepository{db: db}
}

func (a *albumRepository) GetOrCreateCategory(album *entity.Album) error {
	var category entity.Category
	tx := a.db.FirstOrCreate(&category, entity.Category{Name: album.Category.Name})
	if tx.Error != nil {
		return tx.Error
	}
	album.CategoryID = category.ID
	album.Category = category
	return nil
}

func (a *albumRepository) Create(album *entity.Album) (*entity.Album, error) {
	if err := a.GetOrCreateCategory(album); err != nil {
		return nil, err
	}
	if err := a.db.Create(album).Error; err != nil {
		return nil, err
	}
	return album, nil
}

func (a *albumRepository) Get(ID int) (*entity.Album, error) {
	var album = entity.Album{}
	if err := a.db.Preload("Category").First(&album, ID).Error; err != nil {
		return nil, err
	}
	return &album, nil
}

func (a *albumRepository) Save(album *entity.Album) (*entity.Album, error) {
	selectedAlbum, err := a.Get(album.ID)
	if err != nil {
		return nil, err
	}

	if err := a.GetOrCreateCategory(album); err != nil {
		return nil, err
	}

	// Copy files other than files with different time.Time and types.Date
	if err := copier.CopyWithOption(selectedAlbum, album, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}
	if err := a.db.Save(&selectedAlbum).Error; err != nil {
		return nil, err
	}

	return selectedAlbum, nil
}

func (a *albumRepository) Delete(ID int) error {
	album := entity.Album{ID: ID}
	if err := a.db.Where("id = ?", &album.ID).Delete(&album).Error; err != nil {
		return err
	}
	return nil
}
