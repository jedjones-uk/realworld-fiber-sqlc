package dto

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

var DB *pgxpool.Pool

func NewPool() (*pgxpool.Pool, error) {
	// Connect to the dto
	log.Printf("Connecting to the database")
	connStr := "postgres://postgres:postgres@localhost:5432/realworld"
	dbpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		fmt.Printf("Unable to connect to the database: %v\n", err)
		log.Fatal(err)
	}
	//defer dbpool.Close()
	log.Printf("Connected to the database")
	// Check if the connection is successful
	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatal(err, "failed to ping the database")
	}

	return dbpool, nil

}
