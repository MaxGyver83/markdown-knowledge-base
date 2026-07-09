package documents

import "time"

type Document struct {
	ID        int64
	Title     string
	Filename  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
