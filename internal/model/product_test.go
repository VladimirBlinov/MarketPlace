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
