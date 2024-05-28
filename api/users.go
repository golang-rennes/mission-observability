package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	missionErrors "github.com/golang-rennes/mission-observability/errors"
	"github.com/golang-rennes/mission-observability/internal/database"
	"github.com/golang-rennes/mission-observability/logutils"
	"github.com/labstack/echo/v4"
)

type Users struct {
	usersStore database.UsersStore
}

func NewUsers(usersStore database.UsersStore) *Users {
	return &Users{
		usersStore: usersStore,
	}
}

func (u *Users) ListUsers(c echo.Context) error {
	logger := logutils.LoggerFromContext(c.Request().Context())

	users, err := u.usersStore.GetAll(c.Request().Context())
	if err != nil {
		logger.Error(fmt.Sprintf("Unable to get all users : %v", err))
		return err
	}
	return c.JSON(http.StatusOK, users)
}

func (u *Users) GetUser(c echo.Context) error {
	logger := logutils.LoggerFromContext(c.Request().Context())

	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logger.Warn(fmt.Sprintf("Unable to parse id %s as int64", c.Param("id")))
		return err
	}

	user, err := u.usersStore.GetByID(c.Request().Context(), userId)
	if err != nil {
		logger.Warn(fmt.Sprintf("Unable to get user: %v", err))
		return err
	}
	return c.JSON(http.StatusOK, user)
}

func (u *Users) CreateUser(c echo.Context) error {
	logger := logutils.LoggerFromContext(c.Request().Context())

	var user database.User
	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if err != nil {
		err := missionErrors.NewInvalidUserBody(err)
		return err
	}

	logger.Info(fmt.Sprintf("Create user %q", user.Name))
	user, err = u.usersStore.Create(c.Request().Context(), user)
	if err != nil {
		slog.Error(fmt.Sprintf("error creating user : %v", err))
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (u *Users) DeleteUser(c echo.Context) error {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err = u.usersStore.Delete(c.Request().Context(), userId)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, nil)
}
