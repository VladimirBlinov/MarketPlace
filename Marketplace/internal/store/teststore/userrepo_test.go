package teststore_test

import (
	"testing"

	"github.com/VladimirBlinov/MarketPlace/MarketPlace/internal/model"
	"github.com/VladimirBlinov/MarketPlace/MarketPlace/internal/store"
	"github.com/VladimirBlinov/MarketPlace/MarketPlace/internal/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestUserRepo_Create(t *testing.T) {
	s := teststore.New()
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}

func TestFindByEmail(t *testing.T) {
	s := teststore.New()

	email := "user@example.org"

	_, err := s.User().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u := model.TestUser(t)
	u.Email = email

	s.User().Create(u)

	u, err = s.User().FindByEmail(email)

	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestFindById(t *testing.T) {
	s := teststore.New()
	u1 := model.TestUser(t)
	s.User().Create(u1)

	u2, err := s.User().FindById(u1.ID)

	assert.NoError(t, err)
	assert.NotNil(t, u2)
}
