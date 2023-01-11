package mysql

import (
	"reflect"
	"testing"
	"time"

	"github.com/narinderv/snipText/pkg/models"
)

func TestUserGet(t *testing.T) {

	// Skip test if 'Short' flag is on
	if testing.Short() {
		t.Skip("Skipping Integration Test")
	}

	// Test Cases
	tests := []struct {
		name      string
		userID    int
		wantUser  *models.User
		wantError error
	}{
		{
			name:   "Valid User",
			userID: 1,
			wantUser: &models.User{
				ID:       1,
				UserName: "Narinder Verma",
				Email:    "narinderv@example.com",
				Created:  time.Date(2023, 1, 10, 17, 25, 22, 0, time.UTC),
			},
			wantError: nil,
		},
		{
			name:      "Invalid User ID",
			userID:    2,
			wantUser:  nil,
			wantError: models.ErrNoRecord,
		},
		{
			name:      "Zero User ID",
			userID:    0,
			wantUser:  nil,
			wantError: models.ErrNoRecord,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			// Setup database connection and test data
			db, teardownFunc := newTestDB(t)

			defer teardownFunc()
			userModel := UserModel{
				DB: db,
			}

			user, err := userModel.Get(test.userID)
			if err != test.wantError {
				t.Errorf("Want %q got %q", test.wantError, err)
			}

			if !reflect.DeepEqual(user, test.wantUser) {
				t.Errorf("Want %q got %q", test.wantUser, user)
			}
		})
	}
}
