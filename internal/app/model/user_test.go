package model_test

import (
	"testing"

	"github.com/VladimirBlinov/MarketPlace/internal/app/model"
	"github.com/stretchr/testify/assert"
)

func Test_EncryptPasswordBeforeCreate(t *testing.T) {
	u := model.TestUser(t)
	assert.NoError(t, u.EncryptPasswordBeforeCreate())
	assert.NotEmpty(t, u.EncryptedPassword)
}
