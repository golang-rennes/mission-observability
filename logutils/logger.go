package logutils

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type contextKey string

const (
	loggerKey contextKey = "logger-key"
)

// LoggerFromContext gets the logger from the given Golang context
func LoggerFromContext(ctx context.Context) *slog.Logger {
	if ctx.Value(loggerKey) != nil {
		newLoggers, ok := ctx.Value(loggerKey).(*slog.Logger)
		if ok {
			return newLoggers
		}
	}

	return slog.Default()
}

// LoggerToContext set the logger the Golang context
// It the returns the new context, containing new logger.
func LoggerToContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func LoggerContextMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()
			logger := LoggerFromContext(request.Context())
			now := time.Now()
			logger = logger.With(
				"request-id", uuid.New(),
				"request", request.RequestURI,
				"method", request.Method,
				"params", request.URL.Query(),
				"path", request.URL.EscapedPath(),
				"date", now,
			)
			logger.Info(fmt.Sprintf("%s %s", request.Method, request.RequestURI))

			ctx := LoggerToContext(request.Context(), logger)
			c.SetRequest(request.Clone(ctx))
			if err := next(c); err != nil {
				c.Error(err)
			}

			status := c.Response().Status
			logMessage := fmt.Sprintf("%d %s %s", status, request.Method, request.RequestURI)
			switch {
			case status >= http.StatusInternalServerError:
				logger.Error(logMessage)
			case status < http.StatusOK || status >= http.StatusBadRequest:
				logger.Warn(logMessage)
			default:
				logger.Info(logMessage)
			}

			return nil
		}
	}
}
