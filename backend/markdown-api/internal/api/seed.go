package api

import (
	"embed"
	"fmt"
	"log/slog"
	"markdown-api/internal/documents"
	"path"
	"sort"
	"time"
)

//go:embed seeddocs/*.md
var seedDocs embed.FS

type seedDoc struct {
	Title   string
	Content string
}

func loadSeedDocs() ([]seedDoc, error) {
	entries, err := seedDocs.ReadDir("seeddocs")
	if err != nil {
		return nil, fmt.Errorf("read seeddocs dir: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	var docs []seedDoc
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		data, err := seedDocs.ReadFile(path.Join("seeddocs", entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", entry.Name(), err)
		}

		title := deriveTitle(string(data), entry.Name())
		docs = append(docs, seedDoc{
			Title:   title,
			Content: string(data),
		})
	}

	return docs, nil
}

func deriveTitle(content, filename string) string {
	for _, line := range []string{
		"# Project Architecture",
		"# Docker",
		"# Kubernetes",
		"# Terraform",
	} {
		if len(content) >= len(line) && content[:len(line)] == line {
			return line[2:]
		}
	}
	return filename
}

func SeedDatabase(logger *slog.Logger, repository *documents.Repository, storage Storage) error {
	count, err := repository.Count()
	if err != nil {
		return fmt.Errorf("count documents: %w", err)
	}

	if count > 0 {
		logger.Info("database already seeded", "documents", count)
		return nil
	}

	demoDocs, err := loadSeedDocs()
	if err != nil {
		return fmt.Errorf("load seed docs: %w", err)
	}

	logger.Info("seeding database with demo documents", "count", len(demoDocs))

	for _, demo := range demoDocs {
		if err := insertDoc(repository, storage, demo.Title, demo.Content); err != nil {
			return fmt.Errorf("seed %q: %w", demo.Title, err)
		}
	}

	logger.Info("database seeded successfully")
	return nil
}

func ResetDatabase(logger *slog.Logger, repository *documents.Repository, storage Storage) error {
	docs, err := repository.List()
	if err != nil {
		return fmt.Errorf("list documents: %w", err)
	}

	for _, doc := range docs {
		if err := storage.Delete(doc.Filename); err != nil {
			logger.Warn("failed to delete file", "filename", doc.Filename, "error", err)
		}
	}

	if err := repository.DeleteAll(); err != nil {
		return fmt.Errorf("delete all documents: %w", err)
	}

	demoDocs, err := loadSeedDocs()
	if err != nil {
		return fmt.Errorf("load seed docs: %w", err)
	}

	for _, demo := range demoDocs {
		if err := insertDoc(repository, storage, demo.Title, demo.Content); err != nil {
			return fmt.Errorf("reset insert %q: %w", demo.Title, err)
		}
	}

	return nil
}

func insertDoc(repository *documents.Repository, storage Storage, title, content string) error {
	doc := documents.Document{
		Title:     title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := repository.Create(&doc); err != nil {
		return fmt.Errorf("create: %w", err)
	}

	filename := fmt.Sprintf("%d.md", doc.ID)
	if err := storage.Save(filename, content); err != nil {
		return fmt.Errorf("save file: %w", err)
	}

	if err := repository.UpdateFilename(doc.ID, filename); err != nil {
		return fmt.Errorf("update filename: %w", err)
	}

	return nil
}
