/*
Package pgmigration provides support for PostgreSQL database migrations.

This package provide support for migrations written in plain SQL
scripts and also in the Go code.  The migrations written in Go code is
supported through the "Migrate" function (see below).

In your project, place your migrations in a separate folder, for
example, "db/migrate".

Migrations are sorted using their file name and then applied in the
sorted order.  Since sorting is important, name your migrations
accordingly.  For example, add a timestamp before migration name.  Or
use any other ordering scheme you'll like.

Note that migration file names are saved into a table, and the table
is used later on to detect which migrations have already been applied.
In other words, don't rename your migration files once they've been
applied to your DB.

Install go-bindata ( https://github.com/jteeuwen/go-bindata ) and run
this command:

    go-bindata -pkg myapp -o bindata.go db/migrate/

The "bindata.go" file will contain your migrations. Regenerate your
"bindata.go" file whenever you add migrations.

In your app code, import pgmigration package:

    import (
        "log"
        "path/filepath"

        "github.com/baijum/pgmigration"
    )

Then, run the migrations.

Make sure the migrations have an ".sql" ending.

After app startup and after a "sql.DB" instance is initialized in your
app, run the migrations.  Assuming you have a variable called "db"
that points to "sql.DB" and the migrations are located in
"db/migrate", execute the following code:

    ms := pgmigration.MigrationsSource{
        AssetDir: AssetDir,
        Dir: "db/migrate",
    }
    if pg, err := pgmigration.Run(db, ms); err != nil {
        log.Fatal(err)
    }
    pg.Migrate("unique-migrations-name-00001", func(tx *sql.TX) error {...})
    pg.Migrate("unique-migrations-name-00002", func(tx *sql.TX) error {...})

The "Migrate" method can be called to run any migrations written inside
your code.
*/
package pgmigration
