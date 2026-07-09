package database

import (
	"database/sql"
	"fmt"
)

type Migration struct {
	Version string
	SQL     string
}

func Migrate(
	db *sql.DB,
	migrations []Migration,
) error {

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY
		)
	`)

	if err != nil {
		return err
	}

	for _, migration := range migrations {
		applied, err := isApplied(
			db,
			migration.Version,
		)

		if err != nil {
			return err
		}

		if applied {
			continue
		}

		fmt.Println(
			"Applying migration:",
			migration.Version,
		)

		tx, err := db.Begin()

		if err != nil {
			return err
		}

		_, err = tx.Exec(
			migration.SQL,
		)

		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO schema_migrations(version)
			VALUES (?)
		`,
			migration.Version,
		)

		if err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func isApplied(
	db *sql.DB,
	version string,
) (bool, error) {

	var count int

	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM schema_migrations
		WHERE version = ?
	`,
		version,
	).Scan(&count)

	return count > 0, err
}
