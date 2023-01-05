package mysql

import (
	"database/sql" // New import
	"strings"

	"github.com/go-sql-driver/mysql" // New import
	"github.com/narinderv/snipText/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(userName string, email string, password string) error {

	// Create the hash of the password before inserting
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := "INSERT INTO users(user_name, email, hashed_password, created) VALUES(?, ?, ?, NOW())"

	_, err = m.DB.Exec(stmt, userName, email, string(passwordHash))
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	stmt := "SELECT id, user_name, email, hashed_password, created FROM users WHERE id= ?"

	row := m.DB.QueryRow(stmt, id)

	userInfo := &models.User{}

	err := row.Scan(&userInfo.ID, &userInfo.UserName, &userInfo.Email, &userInfo.Password, &userInfo.Created)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {

	// First get hashed password corresponding to the provided email
	stmt := "SELECT id, hashed_password from users WHERE email = ?"

	row := m.DB.QueryRow(stmt, email)

	var id int
	var hashedPassword []byte

	err := row.Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	// Validate the password
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	return id, nil
}
