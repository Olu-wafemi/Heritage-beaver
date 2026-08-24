package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/pressly/goose/v3"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/config"
)

func main() {
	dir := flag.String("dir", "migrations", "directory containing .sql migrations")
	flag.Parse()

	cfg := config.Load()
	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("set dialect: %v", err)
	}

	if err := goose.Up(db, *dir); err != nil {
		log.Fatalf("apply migrations: %v", err)
	}
}
