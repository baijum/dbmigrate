# PostgreSQL Database Migrations

[![GoDoc](https://godoc.org/github.com/baijum/pgmigration?status.svg)](https://godoc.org/github.com/baijum/pgmigration)

Package pgmigration provides support for PostgreSQL database
migrations.

## Install

This package provide support for migrations written in plain SQL
scripts and also in the Go code.  The migrations written in Go code is
supported through the `Migrate` function (see below).

In your project, place your migrations in a separate folder, for
example, **db/migrate**.

**Migrations are sorted using their file name and then applied in the
sorted order.** Since sorting is important, name your migrations
accordingly.  For example, add a timestamp before migration name.  Or
use any other ordering scheme you'll like.

Note that migration file names are saved into a table, and the table
is used later on to detect which migrations have already been applied.
In other words, **don't rename your migration files once they've been
applied to your DB**.

## Use

In your app code, import `pgmigration` package:
```golang
import (
    "log"
    "path/filepath"

    "github.com/baijum/pgmigration"
)
```

Then, run the migrations.

**Make sure the migrations have an .sql ending.**

After app startup and after a `sql.DB` instance is initialized in your
app, run the migrations.  Assuming you have a variable called **db**
that points to `sql.DB` and the migrations are located in
**db/migrate**, execute the following code:

```golang
if pg, err := pgmigration.Run(db, filepath.Join("db", "migrate")); err != nil {
    log.Fatal(err)
}
pg.Migrate("unique-migrations-name-00001", func() error {...})
pg.Migrate("unique-migrations-name-00002", func() error {...})
```

The `Migrate` method can be called to run any migrations written inside
your code.

## Credits

This project is a fork of https://github.com/tanel/dbmigrate written by Tanel Lebedev
