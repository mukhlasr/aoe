package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func initDB(ctx context.Context, dsn string) error {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return fmt.Errorf("failed to ping the db: %w", err)
	}

	DB = conn
	return nil
}
