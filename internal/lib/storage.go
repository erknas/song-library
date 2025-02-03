package lib

import (
	"fmt"

	"github.com/erknas/song-library/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func PoolConfig(cfg *config.Config) (*pgxpool.Config, error) {
	dns := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	return pgxpool.ParseConfig(dns)
}
