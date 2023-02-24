package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/ttagiyeva/entain/internal/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func newClient(conf *config.Config, log *zap.SugaredLogger) (*sql.DB, error) {
	dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.DB.Host,
		conf.DB.Port,
		conf.DB.User,
		conf.DB.Password,
		conf.DB.Name,
	)
	conn, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Errorf("error couldn't open db connection %v", err)
		return nil, err
	}

	return conn, nil
}

// NNewProvider creates a new postgres client.
func NewProvider(lc fx.Lifecycle, conf *config.Config, log *zap.SugaredLogger) (*sql.DB, error) {
	conn, err := newClient(conf, log)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := startMigration(conn, log); err != nil {

				return err
			}

			return nil
		},

		OnStop: func(ctx context.Context) error {
			err := conn.Close()
			if err != nil {
				log.Errorf("error couldn't close db connection %v", err)

				return err
			}

			return nil
		},
	})

	return conn, nil
}
