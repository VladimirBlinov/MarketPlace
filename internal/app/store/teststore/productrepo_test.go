package teststore_test

import (
	"testing"

	"github.com/VladimirBlinov/MarketPlace/internal/app/model"
	"github.com/VladimirBlinov/MarketPlace/internal/app/store/teststore"
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
