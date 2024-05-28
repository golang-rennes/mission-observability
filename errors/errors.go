package error

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Type string

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserInvalidBody = errors.New("invalid user body")
)

const NotFound Type = "not_found"

func NewNotFound(id int64) error {
	return fmt.Errorf("%w : %d ", ErrUserNotFound, id)
}

func NewInvalidUserBody(err error) error {
	return fmt.Errorf("%w: %v", ErrUserInvalidBody, err)
}

func ToEchoHTTPError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, ErrUserNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}
	if errors.Is(err, ErrUserInvalidBody) {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return echo.NewHTTPError(http.StatusInternalServerError, err)
}
