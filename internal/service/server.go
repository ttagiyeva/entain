package service

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
)

// NewServer creates a new echo server.
func NewServer(lc fx.Lifecycle) (*echo.Echo, error) {
	engine := echo.New()

	engine.Use(middleware.Recover())
	engine.Use(middleware.CORS())

	errCh := make(chan error)
	succCh := make(chan int)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				addr := ":8080"
				listener, err := net.Listen("tcp", addr)
				if err != nil {
					engine.Logger.Fatal("shutting down the server is started due to listener")
					errCh <- err
				}
				engine.Listener = listener

				succCh <- 0

				if err := engine.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
					engine.Logger.Fatal("shutting down the server is started due to server")
					errCh <- err
				}
			}()

			select {
			case <-succCh:
				return nil
			case e := <-errCh:
				return e
			case <-ctx.Done():
				return ctx.Err()
			}
		},
		OnStop: func(ctx context.Context) error {
			engine.Logger.Info("shutting down the server gracefully")
			return engine.Shutdown(ctx)
		},
	})

	return engine, nil
}
