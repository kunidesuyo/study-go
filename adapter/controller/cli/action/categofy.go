package action

import (
	"go-api-arch-clean-template/entity"
	"go-api-arch-clean-template/pkg/logger"
	"go-api-arch-clean-template/usecase"
)

var CategoryName string

type CategoryAction struct {
	categoryUseCase usecase.CategoryUseCase
}

func NewCategoryAction(categoryUseCase usecase.CategoryUseCase) *CategoryAction {
	return &CategoryAction{
		categoryUseCase: categoryUseCase,
	}
}

func (a *CategoryAction) CreateCategory(name string) (*entity.Category, error) {
	category, err := entity.NewCategory(name)
	if err != nil {
		return nil, err
	}
	createdCategory, err := a.categoryUseCase.GetOrCreate(category)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return createdCategory, nil
}
