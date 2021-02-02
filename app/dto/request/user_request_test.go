package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRequestToUser(t *testing.T) {
	userRequest := UserRequest{
		Email:    "jupiter.stein@gmail.com",
		FullName: "Júpiter Stein",
	}

	user := userRequest.ToUser()

	assert.Equal(t, user.FullName, userRequest.FullName)
	assert.Equal(t, user.Email, userRequest.Email)
}
