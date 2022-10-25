package teststore_test

import (
	store2 "github.com/VladimirBlinov/MarketPlace/MarketPlace/internal/store"
	"testing"

	"github.com/VladimirBlinov/MarketPlace/MarketPlace/internal/model"
	"github.com/VladimirBlinov/MarketPlace/MarketPlace/internal/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestProductRepo_Create(t *testing.T) {
	s := teststore.New()
	u := model.TestUser(t)
	p := model.TestProduct(t)

	s.User().Create(u)
	p.UserID = u.ID

	mpiList := &model.MarketPlaceItemsList{}
	mpiList.GetMPIList(p)

	assert.NoError(t, s.Product().Create(p, mpiList))
	assert.NotNil(t, p)
	assert.NotNil(t, mpiList)
}

func TestProductRepo_Update(t *testing.T) {
	s := teststore.New()
	u := model.TestUser(t)
	p := model.TestProduct(t)

	s.User().Create(u)
	p.UserID = u.ID

	mpiList := &model.MarketPlaceItemsList{}
	mpiList.GetMPIList(p)

	s.Product().Create(p, mpiList)

	p.Description = "new description"
	p.OzonSKU = 1111111

	mpiList.GetMPIList(p)

	assert.NoError(t, s.Product().Update(p, mpiList))
	assert.Equal(t, p.OzonSKU, s.ProductRepo.Products[p.ProductID].OzonSKU)
	assert.Equal(t, p.Description, s.ProductRepo.Products[p.ProductID].Description)
}

func TestProductRepo_Delete(t *testing.T) {
	s := teststore.New()
	u := model.TestUser(t)
	s.User().Create(u)

	p := model.TestProduct(t)
	p.UserID = u.ID
	mpiList := &model.MarketPlaceItemsList{}
	mpiList.GetMPIList(p)
	s.Product().Create(p, mpiList)

	err := s.Product().Delete(p.ProductID, p.UserID)
	assert.Nil(t, err)

	product, err := s.Product().GetProductById(p.ProductID)
	assert.Error(t, store2.ErrRecordNotFound, err)
	assert.Nil(t, product)
}

func TestProduct_GetProductById(t *testing.T) {
	s := teststore.New()
	u := model.TestUser(t)
	s.User().Create(u)

	p := model.TestProduct(t)
	mpiList1 := &model.MarketPlaceItemsList{}
	mpiList1.GetMPIList(p)
	p.UserID = u.ID
	s.Product().Create(p, mpiList1)

	product, err := s.Product().GetProductById(p.ProductID)

	assert.Nil(t, err)
	assert.NotNil(t, product.OzonSKU)
	assert.NotNil(t, product.WildberriesSKU)
}

func TestProductRepo_FindByUserId(t *testing.T) {
	s := teststore.New()
	u := model.TestUser(t)
	s.User().Create(u)

	p1 := model.TestProduct(t)
	p2 := model.TestProduct(t)

	mpiList1 := &model.MarketPlaceItemsList{}
	mpiList1.GetMPIList(p1)
	p1.UserID = u.ID
	s.Product().Create(p1, mpiList1)

	mpiList2 := &model.MarketPlaceItemsList{}
	p2.UserID = u.ID
	mpiList2.GetMPIList(p2)
	s.Product().Create(p2, mpiList2)

	productsList, err := s.Product().FindByUserId(u.ID)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(productsList))
	assert.NotNil(t, productsList[0].OzonSKU)
	assert.NotNil(t, productsList[0].WildberriesSKU)
}

func TestProductRepo_GetCategories(t *testing.T) {
	s := teststore.New()
	c1 := model.TestCategory(t)
	c2 := model.TestCategory(t)

	s.Product().CreateCategory(c1)
	s.Product().CreateCategory(c2)

	categories, err := s.Product().GetCategories()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(categories))
}

func TestProductRepo_CreateMaterial(t *testing.T) {
	s := teststore.New()
	m := model.TestMaterial(t)

	assert.NoError(t, s.Product().CreateMaterial(m))
	assert.NotNil(t, m)
}

func TestProductRepo_GetMaterials(t *testing.T) {
	s := teststore.New()

	m := model.TestMaterial(t)
	m1 := model.TestMaterial(t)

	s.Product().CreateMaterial(m)
	s.Product().CreateMaterial(m1)

	materials, err := s.Product().GetMaterials()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(materials))
}
