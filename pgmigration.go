package pgmigration

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
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

// Postgres migrates PostgreSQL databases
type Postgres struct {
	database *sql.DB
}

// CreateMigrationsTable create the table to keep track of versions of migration
func (postgres *Postgres) createMigrationsTable() error {
	_, err := postgres.database.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL,
			name TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL,
			CONSTRAINT "PK_migrations_id" PRIMARY KEY (id)
	);`)
	if err != nil {
		return err
	}
	_, err = postgres.database.Exec("CREATE UNIQUE INDEX idx_migrations_name ON migrations(name)")
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return err
		}
	}
	return nil
}

// HasMigrated check for migration
func (postgres *Postgres) hasMigrated(name string) (bool, error) {
	var count int
	err := postgres.database.QueryRow("SELECT count(1) FROM migrations WHERE name = $1", name).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Migrate perform exact migration
func (postgres *Postgres) migrateScript(filename string, migration string) error {
	_, err := postgres.database.Exec(migration)
	if err != nil {
		return err
	}
	_, err = postgres.database.Exec("INSERT INTO migrations(name, created_at) VALUES($1, current_timestamp)", filename)
	return err
}

// Migrate run migrations written in your code using ORM or any other SQL toolkit
func (postgres *Postgres) Migrate(name string, migrations func() error) error {
	err := migrations()
	if err != nil {
		return err
	}
	_, err = postgres.database.Exec("INSERT INTO migrations(name, created_at) VALUES($1, current_timestamp)", name)
	return err
}

// Run migrations written in SQL scripts
func Run(db *sql.DB, migrationsFolder string) (*Postgres, error) {
	postgres := &Postgres{database: db}
	return postgres, applyMigrations(postgres, migrationsFolder)
}

// applyMigrations applies migrations from migrationsFolder to database.
func applyMigrations(database database, migrationsFolder string) error {
	// Initialize migrations table, if it does not exist yet
	if err := database.createMigrationsTable(); err != nil {
		return err
	}

	// Scan migration file names in migrations folder
	d, err := os.Open(migrationsFolder)
	if err != nil {
		return err
	}
	dir, err := d.Readdir(-1)
	if err != nil {
		return err
	}

	// Run migrations
	var sqlFiles []string
	for _, f := range dir {
		ext := filepath.Ext(f.Name())
		if ".sql" == ext {
			sqlFiles = append(sqlFiles, f.Name())
		}
	}
	sort.Strings(sqlFiles)
	for _, filename := range sqlFiles {
		// if exists in migrations table, leave it
		// else execute sql
		migrated, err := database.hasMigrated(filename)
		if err != nil {
			return err
		}
		fullpath := filepath.Join(migrationsFolder, filename)
		if migrated {
			log.Println("Already migrated", fullpath)
			continue
		}
		b, err := ioutil.ReadFile(fullpath)
		if err != nil {
			return err
		}
		migration := string(b)
		if len(migration) == 0 {
			log.Println("Skipping empty file", fullpath)
			continue // empty file
		}
		err = database.migrateScript(filename, migration)
		if err != nil {
			return err
		}
		log.Println("Migrated", fullpath)
	}

	return nil
}
