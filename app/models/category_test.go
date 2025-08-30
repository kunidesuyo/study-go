package models_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"awesomeProject/app/models"
	"awesomeProject/pkg/tester"
)

type CategoryTestSuite struct {
	tester.DBMySQLSuite
}

func TestCategoryTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryTestSuite))
}

func (suite *CategoryTestSuite) TestCategory() {
	category, err := models.GetOrCreateCategory("test")
	suite.Assert().Nil(err)
	suite.Assert().NotNil(category.ID)
	suite.Assert().Equal("test", category.Name)

	category2, err := models.GetOrCreateCategory("test")
	suite.Assert().Nil(err)
	suite.Assert().Equal("test", category2.Name)
	suite.Assert().Equal(category.ID, category.ID)
}
