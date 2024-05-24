package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	DB *pgxpool.Pool
}

func NewPostgres(dbname, user, password string) (*Postgres, error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@localhost/%s?sslmode=disable", user, password, dbname)
	dbpool, err := pgxpool.New(context.Background(), connStr)

	if err != nil {
		return nil, err
	}

	// try connecting to the database
	ticker := time.NewTicker(time.Second)
	timeoutExceeded := time.NewTimer(10 * time.Second)

	for {
		select {
		case <-timeoutExceeded.C:
			ticker.Stop()
			timeoutExceeded.Stop()
			return nil, fmt.Errorf("connection timeout")

		case <-ticker.C:
			if err := dbpool.Ping(context.Background()); err == nil {
				timeoutExceeded.Stop()
				return &Postgres{DB: dbpool}, nil
			}
		}
	}
}
