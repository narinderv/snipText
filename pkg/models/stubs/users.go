package stubs

import (
	"time"

	"github.com/narinderv/snipText/pkg/models"
)

var testUser = &models.User{
	ID:       1,
	UserName: "Narinder",
	Email:    "narinderv@gmail.com",
	Created:  time.Now(),
}

type UserModel struct{}

func (m *UserModel) Insert(userName string, email string, password string) error {

	switch email {
	case "narinderv@gmail.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Get(id int) (*models.User, error) {

	switch id {
	case 1:
		return testUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {

	switch email {
	case "narinderv@gmail.com":
		return 1, nil
	default:
		return 0, models.ErrInvalidCredentials
	}
}
