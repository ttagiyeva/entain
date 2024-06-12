package service

import (
	. "net/http"

	"github.com/labstack/echo/v4"

	"github.com/ttagiyeva/entain/internal/database"
	"github.com/ttagiyeva/entain/internal/transaction/delivery/http"
)

// RegisterRouters registers all routers for the service.
func RegisterRouters(e *echo.Echo, h *http.Handler, db *database.Postgres) error {
	e.GET("/health", healthCheck(db))

	grp := e.Group("api/v1")
	grp.POST("/users/:id/transactions", h.Process)

	return nil
}

// Healthcheck of the service.
func healthCheck(db *database.Postgres) echo.HandlerFunc {
	return func(c echo.Context) error {
		if db.Ping() != nil {
			return c.JSON(StatusInternalServerError, echo.Map{"status": "failed"})
		}

		return c.JSON(StatusOK, echo.Map{"status": "ok"})
	}
}
