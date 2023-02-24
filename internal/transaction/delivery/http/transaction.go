package http

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/ttagiyeva/entain/internal/config"
	"github.com/ttagiyeva/entain/internal/model"
	"github.com/ttagiyeva/entain/internal/transaction"
	"go.uber.org/zap"
)

// Handler is a structure which manages http handlers.
type Handler struct {
	log     *zap.SugaredLogger
	usecase transaction.Usecase
}

// NewHandler creates a new http handler.
func NewHandler(log *zap.SugaredLogger, u transaction.Usecase) *Handler {
	return &Handler{
		log:     log,
		usecase: u,
	}
}

func (h *Handler) Process(ctx echo.Context) error {
	transaction := &model.Transaction{}

	err := ctx.Bind(&transaction)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, model.Error{
			Code:    http.StatusBadRequest,
			Message: model.ErrorBadRequest,
		})
	}

	sourceType := ctx.Request().Header.Get(config.SourceType)
	_, ok := model.SourceType[sourceType]
	if !ok {

		return ctx.JSON(http.StatusBadRequest, model.Error{
			Code:    http.StatusBadRequest,
			Message: model.ErrorInvalidSourceType,
		})
	}

	transaction.SourceType = sourceType
	transaction.UserID = ctx.Param("id")

	if strings.Trim(transaction.TransactionID, " ") == "" {

		return ctx.JSON(http.StatusBadRequest, model.Error{
			Code:    http.StatusBadRequest,
			Message: model.ErrorInvalidTransactionId,
		})
	}

	if transaction.Amount <= 0 {

		return ctx.JSON(http.StatusBadRequest, model.Error{
			Code:    http.StatusBadRequest,
			Message: model.ErrorInvalidAmont,
		})
	}

	_, ok = model.State[transaction.State]
	if !ok {
		return ctx.JSON(http.StatusBadRequest, model.Error{
			Code:    http.StatusBadRequest,
			Message: model.ErrorInvalidState,
		})
	}

	c := ctx.Request().Context()

	err = h.usecase.Process(c, transaction)
	if err != nil {
		return ctx.JSON(getStatusCode(err), model.Error{
			Code:    getStatusCode(err),
			Message: err.Error(),
		})
	}

	return ctx.NoContent(http.StatusOK)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case model.ErrorInternalServer:
		return http.StatusInternalServerError
	case model.ErrorNotFound:
		return http.StatusNotFound
	case model.ErrorInsufficientBalance:
		return http.StatusForbidden
	case model.ErrorTransactionAlreadyExists:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
