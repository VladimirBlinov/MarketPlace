package teststore_test

import (
	"testing"

	"github.com/VladimirBlinov/MarketPlace/internal/model"
	"github.com/VladimirBlinov/MarketPlace/internal/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestProductRepo_Create(t *testing.T) {
	s := teststore.New()
	u := model.TestUser(t)
	p := model.TestProduct(t)

	p.UserID = u.ID

	assert.NoError(t, s.Product().Create(p))
	assert.NotNil(t, p)
}

func TestProductRepo_FindByUserId(t *testing.T) {
	s := teststore.New()
	u := model.TestUser(t)
	s.User().Create(u)

	p1 := model.TestProduct(t)
	p2 := model.TestProduct(t)

	p1.UserID = u.ID
	s.Product().Create(p1)
	p2.UserID = u.ID
	s.Product().Create(p2)

	productsList, err := s.Product().FindByUserId(u.ID)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(productsList))
}
