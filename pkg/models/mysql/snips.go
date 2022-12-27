package mysql

import (
	"database/sql"

	"github.com/narinderv/snipText/pkg/models"
)

type SnipModel struct {
	DB *sql.DB
}

func (m *SnipModel) Insert(title string, content string, expiry string) (int, error) {

	stmt := "INSERT INTO snips(title, content, created, expires) VALUES(?, ?, NOW(), DATE_ADD(NOW(), INTERVAL ? DAY))"

	res, err := m.DB.Exec(stmt, title, content, expiry)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnipModel) Get(id int) (*models.SnipText, error) {
	stmt := "SELECT id, title, content, created, expires FROM snips WHERE expires > NOW() AND id= ?"

	row := m.DB.QueryRow(stmt, id)

	snipInfo := &models.SnipText{}

	err := row.Scan(&snipInfo.ID, &snipInfo.Title, &snipInfo.Content, &snipInfo.Created, &snipInfo.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return snipInfo, nil
}

func (m *SnipModel) GetLatest() ([]*models.SnipText, error) {

	stmt := "SELECT id, title, content, created, expires FROM snips WHERE expires > NOW() ORDER BY created DESC LIMIT 10"

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// Free the result set after use, only if it is a valid one
	// important to free this to ensure connections are released
	defer rows.Close()

	snips := []*models.SnipText{}

	for rows.Next() {
		snip := &models.SnipText{}

		err := rows.Scan(&snip.ID, &snip.Title, &snip.Content, &snip.Created, &snip.Expires)
		if err != nil {
			return nil, err
		}

		snips = append(snips, snip)
	}

	// Check if any error had occured during the iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snips, nil
}
