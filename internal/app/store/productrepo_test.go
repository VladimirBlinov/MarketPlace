package store_test

import (
	"testing"

	"github.com/VladimirBlinov/MarketPlace/internal/app/model"
	"github.com/VladimirBlinov/MarketPlace/internal/app/store"
	"github.com/stretchr/testify/assert"
)

// TODO: fulfill FK tables
func TestProductRepo_Create(t *testing.T) {
	s, teardown := store.TestStore(t, databaseURL)

	defer teardown("product")

	p, err := s.Product().Create(model.TestProduct(t))

	assert.NoError(t, err)
	assert.NotNil(t, p)
}
