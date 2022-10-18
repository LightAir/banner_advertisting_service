package main

import (
	"embed"
	"flag"
	"log"

	"github.com/LightAir/bas/internal/config"
	_ "github.com/lib/pq"
	goose "github.com/pressly/goose/v3"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/banner/config.yaml", "Path to configuration file")
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	flag.Parse()

	cfg, err := config.Parse(configFile)
	if err != nil {
		log.Fatal(err)
	}

	goose.SetBaseFS(embedMigrations)

	db, err := goose.OpenDBWithDriver("postgres", config.GetDsn(cfg.DB.SQL))
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatalf("goose up: %v", err)
	}
}
