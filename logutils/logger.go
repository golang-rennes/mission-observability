package logutils

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/Graylog2/go-gelf/gelf"
	sloggraylog "github.com/samber/slog-graylog/v2"
	slogmulti "github.com/samber/slog-multi"
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
	gelfWriter, err := gelf.NewWriter("localhost:12201")
	if err != nil {
		log.Fatalf("gelf.NewWriter: %s", err)
	}
	fmt.Println("graylog is ok")
	handler := sloggraylog.Option{Level: slog.LevelDebug, Writer: gelfWriter}.NewGraylogHandler()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()

			logger := slog.New(slogmulti.Fanout(handler, slog.Default().Handler()))

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

			logger = logger.With("status", status)
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
