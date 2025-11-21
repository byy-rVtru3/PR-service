package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

type Postgres struct {
	conn *pgx.Conn
}

func (db *Postgres) createConnectPath() (string, error) {
	var dbParam [5]string
	for i, param := range []string{"DB_USER", "DB_PASSWORD", "DB_HOST", "DB_PORT", "DB_NAME"} {
		value := os.Getenv(param)
		if value == "" {
			return "", fmt.Errorf("environment variable %s is not set", param)
		}
		dbParam[i] = value
	}

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbParam[0], // DB_USER
		dbParam[1], // DB_PASSWORD
		dbParam[2], // DB_HOST
		dbParam[3], // DB_PORT
		dbParam[4], // DB_NAME
	)

	return dbURL, nil
}

func NewDB() (*Postgres, error) {
	db := &Postgres{}
	path, err := db.createConnectPath()
	if err != nil {
		return db, err
	}

	time.Sleep(10 * time.Second)

	conn, err := pgx.Connect(context.Background(), path)
	if err != nil {
		return db, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}
	db.conn = conn

	return db, nil
}

func (db *Postgres) CloseDB() error {
	return db.conn.Close(context.Background())
}
