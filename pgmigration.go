package pgmigration

import (
	"database/sql"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

type postgres struct {
	db *sql.DB
}

// CreateMigrationsTable create the table to keep track of versions of migration
func (pg *postgres) createMigrationsTable() error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	err = func(tx *sql.Tx) error {
		_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id BIGSERIAL PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
		);`)
		return err
	}(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

// HasMigrated check for migration
func (pg *postgres) hasMigrated(name string) (bool, error) {
	var count int
	err := pg.db.QueryRow("SELECT count(1) FROM migrations WHERE name = $1", name).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Migrate perform exact migration
func (pg *postgres) migrateScript(filename string, script string) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	err = func(tx *sql.Tx) error {
		_, err := pg.db.Exec(script)
		if err != nil {
			return err
		}
		_, err = pg.db.Exec("INSERT INTO migrations(name) VALUES ($1)", filename)
		return err
	}(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

// Migrate run the migrations written in SQL scripts
func Migrate(db *sql.DB, assetNames func() []string, asset func(name string) ([]byte, error), lastScript *string) error {
	pg := &postgres{db: db}

	// Initialize migrations table, if it does not exist yet
	if err := pg.createMigrationsTable(); err != nil {
		return err
	}

	sqlFiles := assetNames()
	sort.Strings(sqlFiles)
	for _, filename := range sqlFiles {
		_, fn := filepath.Split(filename)
		if strings.HasPrefix(fn, "ignore") {
			log.Println("Script ignored:", filename)
			continue
		}

		ext := filepath.Ext(filename)
		if ".sql" != ext {
			log.Println("File ignored as it has no .sql extension:", filename)
			continue
		}

		// if exists in migrations table, leave it
		// else execute sql
		migrated, err := pg.hasMigrated(filename)
		if err != nil {
			return err
		}
		if migrated {
			log.Println("Already migrated:", filename)
			continue
		}

		b, err := asset(filename)
		if err != nil {
			return err
		}
		script := strings.TrimSpace(string(b))
		if len(script) == 0 {
			log.Println("Skipping empty file:", filename)
			continue // empty file
		}
		// Run migrations
		err = pg.migrateScript(filename, script)
		if err != nil {
			return err
		}
		log.Println("Migrated:", filename)

		if lastScript != nil {
			if *lastScript == filename {
				log.Println("Last script reached:", filename)
				break
			}
		}
	}
	return nil
}
