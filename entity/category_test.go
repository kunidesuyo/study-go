package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"awesomeProject/entity"
)

func TestCategory(t *testing.T) {
	category := entity.Category{
		ID:   1,
		Name: "sports",
	}

	assert.Equal(t, 1, category.ID)
	assert.Equal(t, "sports", string(category.Name))
}
