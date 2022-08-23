package sqlstore_test

import (
	"testing"

	"github.com/VladimirBlinov/MarketPlace/internal/app/model"
	"github.com/VladimirBlinov/MarketPlace/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestProductRepo_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("product", "users")

	s := sqlstore.New(db)
	u := model.TestAdminUser(t)
	s.User().Create(u)
	p := model.TestProduct(t)
	p.UserID = u.ID
	err := s.Product().Create(p)

	assert.NoError(t, err)
	assert.NotNil(t, p)
}
