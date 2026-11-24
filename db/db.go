package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var Db *pgx.Conn

func ConnectDB() {
	var err error

	connectStr := os.Getenv("DB_URL")

	Db, err = pgx.Connect(context.Background(), connectStr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Database connected successfully")
}
