# PostgreSQL Database Migrations

[![GoDoc](https://godoc.org/github.com/baijum/pgmigration?status.svg)](https://godoc.org/github.com/baijum/pgmigration)

Package pgmigration provides support for PostgreSQL database
migrations.

## Install

This package provide support for migrations written in plain SQL
scripts.

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
import "github.com/baijum/pgmigration"
```

Then, run the migrations.

**Make sure the migrations have an .sql ending.**  Any script with file
name prefix as `ignore` will be ignored.

After app startup and after a `sql.DB` instance is initialized in your
app, run the migrations.  Assuming you have a variable called **DB**
that points to `sql.DB` and the migrations are located in
**db/migrations**, execute the following code:

```golang
// DB is the database connection wrapper
var DB *sql.DB

// SchemaMigrate migrate database schema
func SchemaMigrate() error {
	return pgmigration.Migrate(DB, AssetNames, Asset, nil)
}
```

The last argument is string pointer which can refer to the last script,
if there is any.  This will be useful to write tests.

## Credits

This project started as a fork of https://github.com/tanel/dbmigrate written by Tanel Lebedev
