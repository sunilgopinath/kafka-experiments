package consumer

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

const postgresConn = "postgresql://postgres@localhost:5432/user_tracking"

// ConnectDB returns a new PostgreSQL connection
func ConnectDB() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), postgresConn)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	return conn, err
}
