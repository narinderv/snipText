package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("no matching record found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDuplicateEmail     = errors.New("duplicate Email")
)

type SnipText struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type User struct {
	ID       int
	UserName string
	Email    string
	Password []byte
	Created  time.Time
}
