package entity

import (
	"errors"
)

const (
	Food   CategoryName = "food"
	Music  CategoryName = "music"
	Sports CategoryName = "sports"
)

type CategoryName string

func NewCategoryName(value string) (*CategoryName, error) {
	var categoryName CategoryName
	if err := categoryName.Set(value); err != nil {
		return nil, err
	}
	return &categoryName, nil
}

func (c *CategoryName) IsValid() bool {
	return *c == Food || *c == Music || *c == Sports
}

func (c *CategoryName) Set(value string) error {
	newCategoryName := CategoryName(value)
	if !newCategoryName.IsValid() {
		return errors.New("Invalid value for CategoryName")
	}
	*c = newCategoryName
	return nil
}

type Category struct {
	ID   int
	Name CategoryName
}

func NewCategory(name string) (*Category, error) {
	categoryName, err := NewCategoryName(name)
	if err != nil {
		return nil, err
	}
	return &Category{
		Name: *categoryName,
	}, nil
}
