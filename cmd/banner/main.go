package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/LightAir/bas/docs"
	"github.com/LightAir/bas/internal/config"
	"github.com/LightAir/bas/internal/core"
	"github.com/LightAir/bas/internal/logger"
	rmqqueue "github.com/LightAir/bas/internal/queue/rmq"
	"github.com/LightAir/bas/internal/server/http"
	initstorage "github.com/LightAir/bas/internal/storage/init"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/banner/config.yaml", "Path to configuration file")
}

// @title       Banner Advertising Service
// @version     1.0
// @description This is a Banner Advertising Service
//
// @host        localhost:8000
// @BasePath    /
// .
func main() {
	flag.Parse()

	cfg, err := config.Parse(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logg := logger.New()

	storage, err := initstorage.NewStorage(cfg, logg)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	err = storage.Connect(ctx)
	if err != nil {
		log.Fatalf("failed to load driver: %v", err)
	}

	rmq := rmqqueue.NewRmq(cfg)
	err = rmq.Connect(ctx)
	if err != nil {
		log.Fatalf("failed to connect rabbitmq: %v", err)
	}

	app := core.NewApp(storage, cfg, rmq, logg)
	server := http.NewServer(logg, app, cfg)

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("banner is running...")

	go func() {
		if err := server.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
		}
	}()

	<-ctx.Done()
}
