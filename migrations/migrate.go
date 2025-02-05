package migrations

import (
	"errors"
	"fmt"
	"time"

	"github.com/erknas/song-library/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func New(cfg *config.Config) error {
	dns := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	time.Sleep(time.Second * 3)

	m, err := migrate.New(cfg.MigrationPath, dns)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no new migragtion")
			return nil
		}
		return err
	}

	fmt.Println("successful migration")

	return nil
}
