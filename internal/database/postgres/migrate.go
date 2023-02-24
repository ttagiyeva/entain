package postgres

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// startMigration runs the migrations.
func startMigration(conn *sql.DB, log *zap.SugaredLogger) error {
	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		log.Errorf("error while creating postgres driver: %v", err)

		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:////home/turan/hd/Files/personal/projects/Entain/migrations",
		"postgres",
		driver,
	)

	if err != nil {
		log.Errorf("error while creating migrate instance: %v", err)

		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Errorf("error while running migrations: %v", err)

		return err
	}

	return nil
}
