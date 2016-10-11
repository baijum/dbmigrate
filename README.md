# PostgreSQL Database Migrations

[![GoDoc](https://godoc.org/github.com/baijum/pgmigration?status.svg)](https://godoc.org/github.com/baijum/pgmigration)

Package pgmigration provides support for PostgreSQL database
migrations.

## Install

This package provide support for migrations written in plain SQL
scripts and also in the Go code.  The migrations written in Go code is
supported through the `Migrate` function (see below).

In your project, place your migrations in a separate folder, for
example, **db/migrations**.

**Migrations are sorted using their file name and then applied in the
sorted order.** Since sorting is important, name your migrations
accordingly.  For example, add a timestamp before migration name.  Or
use any other ordering scheme you'll like.

Note that migration file names are saved into a table, and the table
is used later on to detect which migrations have already been applied.
In other words, **don't rename your migration files once they've been
applied to your DB**.

Install [go-bindata](https://github.com/jteeuwen/go-bindata) and run this
command:

    go-bindata -pkg myapp -o bindata.go db/migrations/

The `bindata.go` file will contain your migrations. Regenerate your
`bindata.go` file whenever you add migrations.

## Use

In your app code, import `pgmigration` package:
```golang
import (
    "log"

    "github.com/baijum/pgmigration"
)
```

Then, run the migrations.

**Make sure the migrations have an .sql ending.**

After app startup and after a `sql.DB` instance is initialized in your
app, run the migrations.  Assuming you have a variable called **DB**
that points to `sql.DB` and the migrations are located in
**db/migrations**, execute the following code:

```golang
// DB is the database connection wrapper
var DB *sql.DB

// SchemaMigrate migrate database schema
func SchemaMigrate() error {
	ms := pgmigration.NewMigrationsSource(AssetNames, Asset)
	var err error
	pg, err := pgmigration.Run(DB, ms)
	if err != nil {
		return err
	}
	err = pg.Migrate("unique-code-migrations-name-00001", func(tx *sql.Tx) error { return nil })
	if err != nil {
		return err
	}
	return err
}
```

The `Migrate` method can be called to run any migrations written inside
your code.

## Credits

This project is a fork of https://github.com/tanel/dbmigrate written by Tanel Lebedev
