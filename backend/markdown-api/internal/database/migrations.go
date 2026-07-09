package database

import (
	"embed"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func LoadMigrations() ([]Migration, error) {
	entries, err := migrationFiles.ReadDir("migrations")

	if err != nil {
		return nil, err
	}

	var migrations []Migration

	for _, entry := range entries {
		data, err := migrationFiles.ReadFile(
			"migrations/" + entry.Name(),
		)

		if err != nil {
			return nil, err
		}

		migrations = append(migrations, Migration{
			Version: entry.Name(),
			SQL:     string(data),
		})
	}

	return migrations, nil
}
