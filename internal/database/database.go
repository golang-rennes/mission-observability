package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/dlmiddlecote/sqlstats"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
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
	dbx.DB.SetMaxOpenConns(50)
	dbx.DB.SetMaxIdleConns(50)
	dbx.DB.SetConnMaxLifetime(time.Minute)
	dbx.DB.SetConnMaxIdleTime(time.Second * 30)

	collector := sqlstats.NewStatsCollector("mission_observability_db", dbx)
	prometheus.MustRegister(collector)

	db := &DBClient{
		dbx,
	}

	slog.Info(fmt.Sprintf("Successfully connected to database mission-observability"), "database", "mission-observability")

	return db, nil
}
