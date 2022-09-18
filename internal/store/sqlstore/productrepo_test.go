package sqlstore_test

import (
	"testing"

	"github.com/VladimirBlinov/MarketPlace/internal/model"
	"github.com/VladimirBlinov/MarketPlace/internal/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestProductRepo_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("product", "users", "category", "material")

	s := sqlstore.New(db)
	u := model.TestAdminUser(t)
	s.User().Create(u)

	c := model.TestCategory(t)
	s.Product().CreateCategory(c)

	m := model.TestMaterial(t)
	s.Product().CreateMaterial(m)

	p := model.TestProduct(t)
	p.UserID = u.ID
	p.CategoryID = c.CategoryID
	p.MaterialID = m.MaterialID
	err := s.Product().Create(p)

	assert.NoError(t, err)
	assert.NotNil(t, p)
}

func TestProductRepo_FindByUserId(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("product", "users", "category")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	s.User().Create(u)

	c := model.TestCategory(t)
	s.Product().CreateCategory(c)

	m := model.TestMaterial(t)
	s.Product().CreateMaterial(m)

	p1 := model.TestProduct(t)
	p2 := model.TestProduct(t)

	p1.UserID = u.ID
	p1.CategoryID = c.CategoryID
	p1.MaterialID = m.MaterialID
	s.Product().Create(p1)
	p2.UserID = u.ID
	p2.CategoryID = c.CategoryID
	p2.MaterialID = m.MaterialID
	s.Product().Create(p2)

	productsList, err := s.Product().FindByUserId(u.ID)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(productsList))
}

func TestProductRepo_GetCategories(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("category")

	s := sqlstore.New(db)
	c1 := model.TestCategory(t)
	c2 := model.TestCategory(t)

	c1.ParentCategoryID = 0

	s.Product().CreateCategory(c1)
	s.Product().CreateCategory(c2)

	categories, err := s.Product().GetCategories()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(categories))
}

func TestProductRepo_CreateCategory(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("category")

	s := sqlstore.New(db)

	testCases := []struct {
		name     string
		category func() *model.Category
		isValid  bool
	}{
		{
			name: "valid_parent_nil",
			category: func() *model.Category {
				c := model.TestCategory(t)
				c.ParentCategoryID = 0
				return c
			},
			isValid: true,
		},
		{
			name: "valid",
			category: func() *model.Category {
				c := model.TestCategory(t)
				c.ParentCategoryID = 1
				return c
			},
			isValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, s.Product().CreateCategory(tc.category()))
			} else {
				assert.Error(t, s.Product().CreateCategory(tc.category()))
			}
		})
	}
}

func TestProductRepo_CreateMaterial(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("material")

	s := sqlstore.New(db)

	testCases := []struct {
		name     string
		material func() *model.Material
		isValid  bool
	}{
		{
			name: "valid",
			material: func() *model.Material {
				m := *model.TestMaterial(t)
				return &m
			},
			isValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, s.Product().CreateMaterial(tc.material()))
			} else {
				assert.Error(t, s.Product().CreateMaterial(tc.material()))
			}
		})
	}
}

func TestProductRepo_GetMaterials(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("material")

	s := sqlstore.New(db)

	m := model.TestMaterial(t)
	m1 := model.TestMaterial(t)

	s.Product().CreateMaterial(m)
	s.Product().CreateMaterial(m1)

	materials, err := s.Product().GetMaterials()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(materials))
}
