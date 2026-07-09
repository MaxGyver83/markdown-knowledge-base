package storage

import (
	"os"
	"path/filepath"
)

type MarkdownStorage struct {
	basePath string
}

func NewMarkdownStorage(basePath string) *MarkdownStorage {
	return &MarkdownStorage{
		basePath: basePath,
	}
}

func (s *MarkdownStorage) Init() error {
	return os.MkdirAll(
		s.basePath,
		0755,
	)
}

func (s *MarkdownStorage) Save(
	filename string,
	content string,
) error {
	path := filepath.Join(
		s.basePath,
		filename,
	)

	return os.WriteFile(
		path,
		[]byte(content),
		0644,
	)
}

func (s *MarkdownStorage) Load(
	filename string,
) (string, error) {

	path := filepath.Join(
		s.basePath,
		filename,
	)

	data, err := os.ReadFile(path)

	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (s *MarkdownStorage) Delete(
	filename string,
) error {

	path := filepath.Join(
		s.basePath,
		filename,
	)

	return os.Remove(path)
}
