package main

import (
	"context"
	"log"

	"github.com/erknas/song-library/internal/api"
	"github.com/erknas/song-library/internal/config"
	"github.com/erknas/song-library/internal/logger"
	"github.com/erknas/song-library/internal/service"
	"github.com/erknas/song-library/internal/storage"
	"github.com/erknas/song-library/migrations"
)

func main() {
	var (
		ctx    = context.Background()
		cfg    = config.Load()
		logger = logger.New(cfg.Env)
	)

	if err := migrations.New(cfg); err != nil {
		log.Fatal(err)
	}

	store, err := storage.NewPostgresPool(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	srv := service.New(cfg.ThirdPartyAPIURL, logger, store)

	server := api.NewServer(logger, srv)
	server.Start(ctx, cfg)
}
