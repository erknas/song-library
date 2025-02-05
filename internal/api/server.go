package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/erknas/song-library/docs"
	"github.com/erknas/song-library/internal/config"
	"github.com/erknas/song-library/internal/lib"
	"github.com/erknas/song-library/internal/service"
	httpSwagger "github.com/swaggo/http-swagger"
)

const ctxTimeout time.Duration = time.Second * 10

type Server struct {
	log *slog.Logger
	srv service.Servicer
}

func NewServer(log *slog.Logger, srv service.Servicer) *Server {
	return &Server{
		log: log,
		srv: srv,
	}
}

//	@title			song-library API
//	@version		0.0.1
//	@description	API for managing songs
//	@host			localhost:3000
func (s *Server) Start(ctx context.Context, cfg *config.Config) {
	router := http.NewServeMux()

	s.registerRoutes(router)

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	quitch := make(chan os.Signal, 1)
	signal.Notify(quitch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	s.log.Info("starting server", "port", srv.Addr, "addr", fmt.Sprintf("http://localhost%s", srv.Addr))

	<-quitch

	shutdownCtx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		s.log.Error("failed to shutdown server", "error", err)
		return
	}

	s.log.Info("server shutdown")
}

func (s *Server) registerRoutes(router *http.ServeMux) {
	router.HandleFunc("/songs", lib.MakeHTTPFunc(s.handleSongs))
	router.HandleFunc("/song", lib.MakeHTTPFunc(s.handleSong))

	router.Handle("/swagger/", httpSwagger.WrapHandler)
}
