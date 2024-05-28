package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/golang-rennes/mission-observability/internal/database"
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
	users, err := u.usersStore.GetAll(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, users)
}

func (u *Users) GetUser(c echo.Context) error {
	// User ID from path `users/:id`
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	user, err := u.usersStore.GetByID(c.Request().Context(), userId)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)
}

func (u *Users) CreateUser(c echo.Context) error {
	var user database.User
	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return err
	}
	user, err = u.usersStore.Create(c.Request().Context(), user)
	if err != nil {
		slog.Error("error creating user", "err", err)
		c.JSON(http.StatusInternalServerError, err)
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (u *Users) DeleteUser(c echo.Context) error {
	// User ID from path `users/:id`
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
