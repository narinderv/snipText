package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("No matching record found")

type SnipText struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
