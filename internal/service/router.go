package service

import (
	"github.com/labstack/echo/v4"
	"github.com/ttagiyeva/entain/internal/transaction/delivery/http"
)

// RegisterRouters registers all routers for the service.
func RegisterRouters(e *echo.Echo, h *http.Handler) error {
	grp := e.Group("api/v1")
	grp.POST("/users/:id/transactions", h.Process)

	return nil
}
