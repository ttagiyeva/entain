package main

import (
	"context"

	"github.com/golang-migrate/migrate"
	"go.uber.org/fx"

	"github.com/ttagiyeva/entain/internal/config"
	"github.com/ttagiyeva/entain/internal/database"
	"github.com/ttagiyeva/entain/internal/logger"
	"github.com/ttagiyeva/entain/internal/service"
	"github.com/ttagiyeva/entain/internal/transaction"
	"github.com/ttagiyeva/entain/internal/transaction/delivery/http"
	"github.com/ttagiyeva/entain/internal/transaction/repository"
	"github.com/ttagiyeva/entain/internal/transaction/usecase"
	"github.com/ttagiyeva/entain/internal/user"
	userRepo "github.com/ttagiyeva/entain/internal/user/repository"
)

// main is the entry point of the application.
func main() {
	fx.New(
		fx.Provide(
			config.New,
			logger.NewLogger,
			service.NewServer,
			http.NewHandler,
			database.NewPostgres,

			fx.Annotate(
				func(postgres *database.Postgres) transaction.Database {
					return postgres
				},

				fx.As(new(transaction.Database)),
			),

			fx.Annotate(
				usecase.New,
				fx.As(new(transaction.Usecase)),
			),

			fx.Annotate(
				func(postgres *database.Postgres) transaction.Repository {
					return repository.New(postgres.Connection)
				},

				fx.As(new(transaction.Repository)),
			),

			fx.Annotate(
				func(postgres *database.Postgres) user.Repository {
					return userRepo.New(postgres.Connection)
				},

				fx.As(new(user.Repository)),
			),
		),
		// Creating connection to database
		fx.Invoke(
			func(p *database.Postgres, c *config.Config) {
				err := p.Connect(context.Background(), c)
				if err != nil {
					panic(err)
				}
			},
		),
		// Executing database migrations
		fx.Invoke(
			func(p *database.Postgres) {
				err := p.MigrateUp()
				if err != nil && err == migrate.ErrNoChange {
					panic(err)
				}
			},
		),
		fx.Invoke(
			func(uc transaction.Usecase) {
				go uc.PostProcess(context.Background())
			},

			service.RegisterRouters,
		),
	).Run()
}
