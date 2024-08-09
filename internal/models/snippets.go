package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
type SnippetModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := "INSERT INTO snippets (title, content, created, expires) VALUES (? , ?, datetime('now'), date(date('now'), ?));"
	rs, err := m.DB.Exec(stmt, title, content, fmt.Sprintf("+%d days", expires))
	if err != nil {
		return 0, err
	}
	id, err := rs.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := "SELECT id, title, content, created, expires FROM snippets" +
		" WHERE expires > datetime('now') AND id = ?;"

	rs := m.DB.QueryRow(stmt, id)

	snippet := &Snippet{}

	err := rs.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorNotFound
		} else {
			return nil, err
		}
	}
	return snippet, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := "SELECT id, title, content, created, expires FROM snippets" +
		" WHERE expires > datetime('now') order by created desc limit 10;"

	rows, err := m.DB.Query(stmt)
	defer rows.Close()

	if err != nil {
		return nil, err
	}
	snippets := []*Snippet{}
	for rows.Next() {
		s := &Snippet{}
		rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
