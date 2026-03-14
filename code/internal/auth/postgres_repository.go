package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) FindByLoginID(ctx context.Context, loginID string) (*User, error) {
	const query = `
		SELECT id, login_id, password_hash, failed_attempts, locked_until
		FROM users
		WHERE login_id = $1
	`

	var user User
	err := r.pool.QueryRow(ctx, query, loginID).Scan(
		&user.ID,
		&user.LoginID,
		&user.PasswordHash,
		&user.FailedAttempts,
		&user.LockedUntil,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by login id: %w", err)
	}

	return &user, nil
}

func (r *PostgresRepository) IncrementFailure(ctx context.Context, userID int64, now time.Time) (int, *time.Time, error) {
	const query = `
		UPDATE users
		SET
			failed_attempts = failed_attempts + 1,
			locked_until = CASE
				WHEN failed_attempts + 1 >= $2 THEN $3
				ELSE locked_until
			END,
			updated_at = $1
		WHERE id = $4
		RETURNING failed_attempts, locked_until
	`

	lockedUntilValue := now.Add(LockDuration)
	var failedAttempts int
	var lockedUntil *time.Time
	err := r.pool.QueryRow(ctx, query, now, MaxFailedAttempts, lockedUntilValue, userID).Scan(&failedAttempts, &lockedUntil)
	if err != nil {
		return 0, nil, fmt.Errorf("increment failure attempts: %w", err)
	}

	return failedAttempts, lockedUntil, nil
}

func (r *PostgresRepository) ResetFailures(ctx context.Context, userID int64) error {
	const query = `
		UPDATE users
		SET failed_attempts = 0, locked_until = NULL, updated_at = NOW()
		WHERE id = $1
	`

	if _, err := r.pool.Exec(ctx, query, userID); err != nil {
		return fmt.Errorf("reset failures: %w", err)
	}
	return nil
}
