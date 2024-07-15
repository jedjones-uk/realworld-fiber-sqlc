package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"realworld-fiber-sqlc/usecase/database/sqlc"
	"time"
)

var DB *pgxpool.Pool

func UpdateProfileTX(db *pgxpool.Pool, profileParams *sqlc.UpdateProfileParams, userParams *sqlc.UpdateUserParams) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			fmt.Printf("failed to rollback transaction: %v", err)
		}
	}(tx, ctx)

	queries := sqlc.New(tx)
	qtx := queries.WithTx(tx)

	if err := qtx.UpdateProfile(ctx, *profileParams); err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	if err := qtx.UpdateUser(ctx, *userParams); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
