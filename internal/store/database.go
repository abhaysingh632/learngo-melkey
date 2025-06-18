package store

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx", "host=localhost port=5432 user=postgres password=postgres dbname=postgres")

	if err != nil {
		return nil, fmt.Errorf("DB Open: %w", err)
	}
	fmt.Fprint(os.Stdout, "Connected to Database Successfully!!!\n")
	return db, err
}

func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()

	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	
	if err != nil {
		return fmt.Errorf("migrate %w", err)
	}


	err = goose.Up(db, dir)

	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}


	return nil
}
