package database

import (
	"context"
	"database/sql"
	"errors"

	missionErrors "github.com/golang-rennes/mission-observability/errors"
	"github.com/golang-rennes/mission-observability/logutils"
)

type UsersStore interface {
	GetAll(ctx context.Context) ([]User, error)
	GetByID(ctx context.Context, id int64) (User, error)
	Create(ctx context.Context, user User) (User, error)
	Delete(ctx context.Context, id int) error
}

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (db *DBClient) GetAll(ctx context.Context) ([]User, error) {
	logger := logutils.LoggerFromContext(ctx)
	logger.Info("Getting all users")

	users := []User{}
	rows, err := db.QueryxContext(ctx, "SELECT id, name FROM users")
	if err != nil {
		return users, err
	}
	for rows.Next() {
		var user User
		err := rows.StructScan(&user)
		if err != nil {
			return []User{}, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (db *DBClient) GetByID(ctx context.Context, id int64) (User, error) {
	var user User
	err := db.GetContext(ctx, &user, "SELECT id, name FROM users WHERE id = $1", id)
	if errors.Is(err, sql.ErrNoRows) {
		return user, missionErrors.NewNotFound(id)
	}
	return user, err
}

func (db *DBClient) Create(ctx context.Context, user User) (User, error) {
	lastInsertId := 0

	err := db.QueryRowContext(ctx, `INSERT INTO users (name) VALUES($1) returning id`, user.Name).Scan(&lastInsertId)
	if err != nil {
		return user, err
	}
	user.ID = int64(lastInsertId)
	return user, err
}

func (db *DBClient) Delete(ctx context.Context, id int) error {
	_, err := db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	return err
}
