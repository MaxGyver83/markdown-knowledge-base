package documents

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Get(id int64) (Document, error) {
	var doc Document

	err := r.db.QueryRow(`
		SELECT
			id,
			title,
			filename,
			created_at,
			updated_at
		FROM documents
		WHERE id = ?
	`,
		id,
	).Scan(
		&doc.ID,
		&doc.Title,
		&doc.Filename,
		&doc.CreatedAt,
		&doc.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return Document{}, ErrNotFound
	}

	if err != nil {
		return Document{}, err
	}

	return doc, nil
}

func (r *Repository) List() ([]Document, error) {
	rows, err := r.db.Query(`
		SELECT
			id,
			title,
			filename,
			created_at,
			updated_at
		FROM documents
		ORDER BY created_at DESC
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var documents []Document

	for rows.Next() {
		var doc Document

		err := rows.Scan(
			&doc.ID,
			&doc.Title,
			&doc.Filename,
			&doc.CreatedAt,
			&doc.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		documents = append(documents, doc)
	}

	return documents, rows.Err()
}

func (r *Repository) Create(doc *Document) error {
	now := time.Now()

	result, err := r.db.Exec(`
		INSERT INTO documents (
			title,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?)
	`,
		doc.Title,
		now,
		now,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	doc.ID = id
	doc.Filename = fmt.Sprintf("%d.md", doc.ID)
	doc.CreatedAt = now
	doc.UpdatedAt = now

	return nil
}

func (r *Repository) Update(
	id int64,
	title string,
) error {
	_, err := r.db.Exec(`
		UPDATE documents
		SET title = ?,
		    updated_at = ?
		WHERE id = ?
	`,
		title,
		time.Now(),
		id,
	)

	return err
}

func (r *Repository) Delete(id int64) error {
	_, err := r.db.Exec(`
		DELETE FROM documents
		WHERE id = ?
	`,
		id,
	)

	return err
}

func (r *Repository) UpdateFilename(
	id int64,
	filename string,
) error {
	_, err := r.db.Exec(`
		UPDATE documents
		SET filename = ?,
		    updated_at = ?
		WHERE id = ?
	`,
		filename,
		time.Now(),
		id,
	)

	return err
}
