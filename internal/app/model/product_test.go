package model_test

import (
	"testing"

	"github.com/VladimirBlinov/MarketPlace/internal/app/model"
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
