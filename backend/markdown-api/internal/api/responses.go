package api

import "time"

type DocumentMetadata struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DocumentResponse struct {
	DocumentMetadata

	Content string `json:"content"`
}
