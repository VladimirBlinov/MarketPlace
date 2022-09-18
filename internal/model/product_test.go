package model_test

import (
	"strings"
	"testing"

	"github.com/VladimirBlinov/MarketPlace/internal/model"
	"github.com/stretchr/testify/assert"
)

func Test_ProductValidate(t *testing.T) {
	testCases := []struct {
		name    string
		p       func() *model.Product
		isValid bool
	}{
		{
			name: "valid",
			p: func() *model.Product {
				return model.TestProduct(t)
			},
			isValid: true,
		},
		{
			name: "wrong category ID",
			p: func() *model.Product {
				p := model.TestProduct(t)
				p.CategoryID = 0
				return p
			},
			isValid: false,
		},
		{
			name: "short name",
			p: func() *model.Product {
				p := model.TestProduct(t)
				p.ProductName = ""
				return p
			},
			isValid: false,
		},
		{
			name: "long name",
			p: func() *model.Product {
				p := model.TestProduct(t)
				p.ProductName = strings.Repeat("a", 210)
				return p
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.p().Validate())
			} else {
				assert.Error(t, tc.p().Validate())
			}
		})
	}

}

func Test_CategoryValidate(t *testing.T) {
	testCases := []struct {
		name    string
		c       func() *model.Category
		isValid bool
	}{
		{
			name: "valid",
			c: func() *model.Category {
				return model.TestCategory(t)
			},
			isValid: true,
		},
		{
			name: "wrong parent category ID",
			c: func() *model.Category {
				c := model.TestCategory(t)
				c.ParentCategoryID = 0
				return c
			},
			isValid: true,
		},
		{
			name: "short name",
			c: func() *model.Category {
				c := model.TestCategory(t)
				c.CategoryName = ""
				return c
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.c().ValidateCategory())
			} else {
				assert.Error(t, tc.c().ValidateCategory())
			}
		})
	}
}

func Test_MaterialValidate(t *testing.T) {
	testCases := []struct {
		name    string
		m       func() *model.Material
		isValid bool
	}{
		{
			name: "valid",
			m: func() *model.Material {
				m := model.TestMaterial(t)
				return m
			},
			isValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.m().ValidateMaterial())
			} else {
				assert.Error(t, tc.m().ValidateMaterial())
			}
		})
	}
}
