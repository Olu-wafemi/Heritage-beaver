package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/config"
	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/store/postgres"
)

func main() {
	dir := flag.String("dir", "migrations", "directory containing .sql migrations")
	flag.Parse()

	cfg := config.Load()
	db, err := postgres.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	entries, err := os.ReadDir(*dir)
	if err != nil {
		log.Fatalf("read migrations dir: %v", err)
	}

	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		files = append(files, filepath.Join(*dir, entry.Name()))
	}
	sort.Strings(files)

	for _, file := range files {
		contents, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("read migration %s: %v", file, err)
		}

		if _, err := db.ExecContext(context.Background(), string(contents)); err != nil {
			log.Fatalf("apply migration %s: %v", file, err)
		}

		fmt.Printf("applied %s\n", file)
	}
}
