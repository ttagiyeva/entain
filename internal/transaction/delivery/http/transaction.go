package http

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/ttagiyeva/entain/internal/constants"
	"github.com/ttagiyeva/entain/internal/model"
	"github.com/ttagiyeva/entain/internal/transaction"
)

// Handler is a structure which manages http handlers.
type Handler struct {
	log     *slog.Logger
	usecase transaction.Usecase
}

// NewHandler creates a new http handler.
func NewHandler(log *slog.Logger, u transaction.Usecase) *Handler {
	return &Handler{
		log:     log,
		usecase: u,
	}
}

func (h *Handler) Process(ctx echo.Context) error {
	transaction := &model.Transaction{}

	err := ctx.Bind(transaction)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, model.Error{
			Code:    http.StatusBadRequest,
			Message: model.ErrorBadRequest,
		})
	}

	sourceType := ctx.Request().Header.Get(constants.SourceType)
	_, ok := constants.SourceTypes[sourceType]
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
			Message: model.ErrorInvalidAmount,
		})
	}

	_, ok = constants.States[transaction.State]
	if !ok {
		return ctx.JSON(http.StatusBadRequest, model.Error{
			Code:    http.StatusBadRequest,
			Message: model.ErrorInvalidState,
		})
	}

	c := ctx.Request().Context()

	err = h.usecase.Process(c, transaction)
	if err != nil {
		h.log.With("body", transaction).Error("processing transaction failed", "error", err)

		resp := getStatusCode(err)

		return ctx.JSON(resp.Code, resp)
	}

	return ctx.NoContent(http.StatusOK)
}

func getStatusCode(err error) model.Error {
	if errors.Is(err, model.ErrorNotFound) {
		return model.Error{Code: http.StatusNotFound, Message: model.ErrorNotFound.Error()}
	}
	if errors.Is(err, model.ErrorInsufficientBalance) {
		return model.Error{Code: http.StatusForbidden, Message: model.ErrorInsufficientBalance.Error()}
	}
	if errors.Is(err, model.ErrorTransactionAlreadyExists) {
		return model.Error{Code: http.StatusConflict, Message: model.ErrorTransactionAlreadyExists.Error()}
	}

	return model.Error{Code: http.StatusInternalServerError, Message: model.ErrorInternalServer.Error()}
}
