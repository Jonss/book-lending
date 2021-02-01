package request

import (
	"testing"

	"github.com/google/uuid"
)

func TestUserRequestToUser(t *testing.T) {
	userRequest := UserRequest{
		Email:    "jupiter.stein@gmail.com",
		FullName: "Júpiter Stein",
	}

	user := userRequest.ToUser()

	if user.Email != userRequest.Email {
		t.Errorf("Email expected is %s. Got %s", userRequest.Email, user.Email)
	}

	if user.FullName != userRequest.FullName {
		t.Errorf("FullName expected is %s. Got %s", userRequest.FullName, user.FullName)
	}

	if user.LoggedUserId == uuid.Nil {
		t.Errorf("LoggedUserId must exist. Received %s", user.LoggedUserId)
	}
}
