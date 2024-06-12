package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"

	"github.com/ttagiyeva/entain/internal/config"
)

//go:embed migrations/*.sql
var fs embed.FS

type Postgres struct {
	Connection *sqlx.DB
	m          *migrate.Migrate
}

// NewPostgres creates a new Postgres instance.
func NewPostgres() *Postgres {
	return &Postgres{}
}

func createConnectionString(host string, port uint16, user, password, dbName string) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName,
	)
}

func (p *Postgres) Connect(ctx context.Context, config *config.Config) error {
	conn, err := sqlx.ConnectContext(ctx, "postgres", createConnectionString(
		config.DB.Host,
		config.DB.Port,
		config.DB.User,
		config.DB.Password,
		config.DB.Name,
	))
	if err != nil {
		return err
	}

	p.Connection = conn

	d, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}

	migratePostgres, err := postgres.WithInstance(p.Connection.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance(
		"iofs", d,
		"postgres", migratePostgres,
	)
	if err != nil {
		return err
	}

	p.m = m

	return nil
}

// MigrateUp runs up database migrations.
func (p *Postgres) MigrateUp() error {
	return p.m.Up()
}

// MigrateDown runs down database migrations.
func (p *Postgres) MigrateDown() error {
	return p.m.Down()
}

// Ping verifies a connection to the database is still alive, establishing a connection if necessary.
func (p *Postgres) Ping() error {
	return p.Connection.Ping()
}

// Begin starts a transaction and returns it.
func (p *Postgres) BeginTx(ctx context.Context) (*sql.Tx, error) {
	tx, err := p.Connection.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// Rollback aborts the given transaction.
func (p *Postgres) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}

// Commit commits the given transaction.
func (p *Postgres) Commit(tx *sql.Tx) error {
	return tx.Commit()
}
