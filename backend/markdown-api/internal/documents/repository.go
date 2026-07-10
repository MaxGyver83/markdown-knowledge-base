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
		WHERE id = $1
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

	err := r.db.QueryRow(`
		INSERT INTO documents (
			title,
			filename,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`,
		doc.Title,
		fmt.Sprintf("tmp_%d", now.UnixNano()),
		now,
		now,
	).Scan(&doc.ID)

	if err != nil {
		return err
	}

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
		SET title = $1,
		    updated_at = $2
		WHERE id = $3
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
		WHERE id = $1
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
		SET filename = $1,
		    updated_at = $2
		WHERE id = $3
	`,
		filename,
		time.Now(),
		id,
	)

	return err
}
