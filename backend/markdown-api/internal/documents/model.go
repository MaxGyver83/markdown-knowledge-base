package documents

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("document not found")

type Document struct {
	ID        int64
	Title     string
	Filename  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
