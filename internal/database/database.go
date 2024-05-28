package database

import (
	"context"
	"fmt"
	"log/slog"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type DBClient struct {
	*sqlx.DB
}

// const connString = "postgresql://postgres:Password123@localhost:5432/moviesdb?sslmode=disable"
func NewDBClient(ctx context.Context, connString string) (*DBClient, error) {
	dbx, err := sqlx.ConnectContext(ctx, "pgx", connString)
	if err != nil {
		return nil, err
	}

	db := &DBClient{
		dbx,
	}

	slog.Info(fmt.Sprintf("Successfully connected to database mission-observability"), "database", "mission-observability")

	return db, nil
}
