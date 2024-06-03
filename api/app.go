package api

import (
	"context"
	"net/http"
	"runtime/pprof"

	"github.com/golang-rennes/mission-observability/config"
	missionErrors "github.com/golang-rennes/mission-observability/errors"
	"github.com/golang-rennes/mission-observability/logutils"

	"github.com/golang-rennes/mission-observability/internal/database"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func Run(ctx context.Context, config *config.Config) error {
	router := echo.New()
	router.Logger.SetHeader("${time_rfc3339} ${status}")
	router.Use(middleware.Recover())

	router.HTTPErrorHandler = func(err error, c echo.Context) {
		err = missionErrors.ToEchoHTTPError(err)
		router.DefaultHTTPErrorHandler(err, c)
	}
	router.Use(otelecho.Middleware("mission-observability", otelecho.WithTracerProvider(otel.GetTracerProvider()), otelecho.WithPropagators(otel.GetTextMapPropagator())))
	router.Use(pprofMiddleware())
	router.Use(logutils.LoggerContextMiddleware())

	router.Use(echoprometheus.NewMiddleware("mission_observability"))
	router.GET("/metrics", echoprometheus.NewHandler())

	db, err := database.NewDBClient(ctx, config.ConnString)
	if err != nil {
		return err
	}

	router.GET("/freedom", func(c echo.Context) error {
		return c.String(http.StatusOK, "For the Managed Democraty!")
	})
	users := NewUsers(db)

	router.GET("/users", users.ListUsers)
	router.POST("/users", users.CreateUser)
	router.GET("/users/:id", users.GetUser)
	router.DELETE("/users/:id", users.DeleteUser)

	router.GET("/connexion_boum", users.ListUsersBoum)
	router.GET("/factorial", users.FactorialUsers)

	return router.Start(":8080")
}

func pprofMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()
			traceId := GetTraceID(request.Context())
			pprof.Do(c.Request().Context(),
				pprof.Labels(
					"span", traceId,
				),
				func(ctx context.Context) {
					c.SetRequest(request.Clone(ctx))
					if err := next(c); err != nil {
						c.Error(err)
					}
				},
			)
			return nil
		}
	}
}

func GetTraceID(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		traceID := spanCtx.TraceID()
		return traceID.String()
	}
	return ""
}
