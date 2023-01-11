package stubs

import (
	"time"

	"github.com/narinderv/snipText/pkg/models"
)

var testSnip = &models.SnipText{
	ID:      1,
	Title:   "Snip For Testing",
	Content: "This is a snip for testing purpose",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnipModel struct{}

func (m *SnipModel) Insert(title string, content string, expiry string) (int, error) {

	return 2, nil
}

func (m *SnipModel) Get(id int) (*models.SnipText, error) {

	switch id {
	case 1:
		return testSnip, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *SnipModel) GetLatest() ([]*models.SnipText, error) {

	return []*models.SnipText{testSnip}, nil
}
