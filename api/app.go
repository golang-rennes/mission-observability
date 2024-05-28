package api

import (
	"context"
	"net/http"

	"github.com/golang-rennes/mission-observability/config"
	"github.com/golang-rennes/mission-observability/internal/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run(ctx context.Context, config *config.Config) error {
	router := echo.New()
	router.Logger.SetHeader("${time_rfc3339} ${status}")
	router.Use(middleware.Recover())
	router.Use(middleware.Logger())

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

	return router.Start(":8080")
}
