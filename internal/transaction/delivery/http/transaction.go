package http

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"

	"github.com/ttagiyeva/entain/internal/model"
	"github.com/ttagiyeva/entain/internal/transaction"
)

const (
	SourceType = "Source-Type"
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

	transaction.UserID = ctx.Param("id")
	transaction.SourceType = ctx.Request().Header.Get(SourceType)

	sv := validator.New()

	err = sv.Struct(transaction)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, h.validatorError(err))
	}

	c := ctx.Request().Context()

	err = h.usecase.Process(c, transaction)
	if err != nil {
		h.log.With("body", transaction).Error("failed to process transaction", "error", err)

		resp := getError(err)

		return ctx.JSON(resp.Code, resp)
	}

	return ctx.NoContent(http.StatusOK)
}

func (h *Handler) validatorError(err error) model.Error {
	if _, ok := err.(*validator.InvalidValidationError); ok {
		h.log.Error("failed to assert validation error", "error", err)

		return model.Error{Code: http.StatusInternalServerError, Message: "Internal Server Error"}
	}

	var sb strings.Builder

	for i, err := range err.(validator.ValidationErrors) {
		if i > 0 {
			sb.WriteString(", ")
		}

		switch err.Tag() {
		case "required":
			sb.WriteString(fmt.Sprintf("%s field is required", err.Field()))
		case "oneof":
			sb.WriteString(fmt.Sprintf("Value of the %s field must be one of '%s'", err.Field(), err.Param()))
		case "gt":
			sb.WriteString(fmt.Sprintf("Value of the %s field must be greater than %s", err.Field(), err.Param()))
		}
	}

	return model.Error{Code: http.StatusBadRequest, Message: sb.String()}
}

func getError(err error) model.Error {
	switch {
	case errors.Is(err, model.ErrorUserNotFound):
		return model.Error{Code: http.StatusNotFound, Message: model.ErrorUserNotFound.Error()}
	case errors.Is(err, model.ErrorInsufficientBalance):
		return model.Error{Code: http.StatusForbidden, Message: model.ErrorInsufficientBalance.Error()}
	case errors.Is(err, model.ErrorTransactionAlreadyExists):
		return model.Error{Code: http.StatusConflict, Message: model.ErrorTransactionAlreadyExists.Error()}
	default:
		return model.Error{Code: http.StatusInternalServerError, Message: model.ErrorInternalServerError.Error()}
	}
}
