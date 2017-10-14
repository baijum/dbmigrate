package pgmigration

import (
	"database/sql"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

// database interface needs to be inmplemented to migrate a new type of database
type database interface {
	createMigrationsTable() error
	hasMigrated(name string) (bool, error)
	migrateScript(filename string, migration string) error
}

// postgres migrates PostgreSQL databases
type postgres struct {
	database *sql.DB
}

// CreateMigrationsTable create the table to keep track of versions of migration
func (postgres *postgres) createMigrationsTable() error {
	tx, err := postgres.database.Begin()
	if err != nil {
		return err
	}
	err = func(tx *sql.Tx) error {
		_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL,
			name TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL,
			CONSTRAINT "pk_migrations_id" PRIMARY KEY (id)
		);`)
		if err != nil {
			return err
		}
		_, err = tx.Exec("CREATE UNIQUE INDEX idx_migrations_name ON migrations(name)")
		if err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				return err
			}
		}
		return nil
	}(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

// HasMigrated check for migration
func (postgres *postgres) hasMigrated(name string) (bool, error) {
	var count int
	err := postgres.database.QueryRow("SELECT count(1) FROM migrations WHERE name = $1", name).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Migrate perform exact migration
func (postgres *postgres) migrateScript(filename string, migration string) error {
	tx, err := postgres.database.Begin()
	if err != nil {
		return err
	}
	err = func(tx *sql.Tx) error {
		_, err := postgres.database.Exec(migration)
		if err != nil {
			return err
		}
		_, err = postgres.database.Exec("INSERT INTO migrations(name, created_at) VALUES($1, current_timestamp)", filename)
		return err
	}(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

// Run migrations written in SQL scripts
func Run(db *sql.DB, migrationsSource MigrationsSource) error {
	database := &postgres{database: db}

	// Initialize migrations table, if it does not exist yet
	if err := database.createMigrationsTable(); err != nil {
		return err
	}

	sqlFiles := migrationsSource.AssetNames()
	sort.Strings(sqlFiles)
	for _, filename := range sqlFiles {
		ext := filepath.Ext(filename)
		if ".sql" != ext {
			continue
		}

		// if exists in migrations table, leave it
		// else execute sql
		migrated, err := database.hasMigrated(filename)
		if err != nil {
			return err
		}
		if migrated {
			log.Println("Already migrated", filename)
			continue
		}
		b, err := migrationsSource.Asset(filename)
		if err != nil {
			return err
		}
		migration := string(b)
		if len(migration) == 0 {
			log.Println("Skipping empty file", filename)
			continue // empty file
		}
		// Run migrations
		err = database.migrateScript(filename, migration)
		if err != nil {
			return err
		}
		log.Println("Migrated", filename)
	}

	return nil
}

// MigrationsSource includes migrations diectory and bindata AssetDir call back function
type MigrationsSource struct {
	AssetNames func() []string
	Asset      func(name string) ([]byte, error)
}

// NewMigrationsSource create a migratios source object
// which includes migrations diectory and bindata AssetDir call back function
func NewMigrationsSource(assetNames func() []string, asset func(name string) ([]byte, error)) MigrationsSource {
	ms := MigrationsSource{
		AssetNames: assetNames,
		Asset:      asset,
	}
	return ms
}
