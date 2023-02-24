package main

import (
	"context"

	"github.com/ttagiyeva/entain/internal/config"
	"github.com/ttagiyeva/entain/internal/database/postgres"
	"github.com/ttagiyeva/entain/internal/logger"
	"github.com/ttagiyeva/entain/internal/service"
	"github.com/ttagiyeva/entain/internal/transaction"
	"github.com/ttagiyeva/entain/internal/transaction/delivery/http"
	repository "github.com/ttagiyeva/entain/internal/transaction/repository/postgres"
	"github.com/ttagiyeva/entain/internal/transaction/usecase"
	userRepo "github.com/ttagiyeva/entain/internal/user/repository/postgres"
	"go.uber.org/fx"
)

// main is the entry point of the application.
func main() {
	fx.New(
		fx.Provide(
			config.New,
			postgres.NewProvider,
			logger.NewZapLogger,
			service.NewServer,

			userRepo.New,
			http.NewHandler,
			repository.New,
			usecase.New,
		),
		fx.Invoke(func(uc transaction.Usecase) {
			go uc.PostProcess(context.Background())
		},
			service.RegisterRouters),
	).Run()
}
