package testaddress

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type queryRowScanner interface {
	Scan(dest ...any) error
}

type rowsScanner interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
	Close()
}

type dbTx interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type dbRunner interface {
	QueryRow(ctx context.Context, sql string, arguments ...any) queryRowScanner
	Query(ctx context.Context, sql string, arguments ...any) (rowsScanner, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (dbTx, error)
}

type pgxPoolRunner struct {
	pool *pgxpool.Pool
}

func (r pgxPoolRunner) QueryRow(ctx context.Context, sql string, arguments ...any) queryRowScanner {
	return r.pool.QueryRow(ctx, sql, arguments...)
}

func (r pgxPoolRunner) Query(ctx context.Context, sql string, arguments ...any) (rowsScanner, error) {
	rows, err := r.pool.Query(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}
	return pgxRowsRunner{rows: rows}, nil
}

func (r pgxPoolRunner) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return r.pool.Exec(ctx, sql, arguments...)
}

func (r pgxPoolRunner) Begin(ctx context.Context) (dbTx, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return pgxTxRunner{tx: tx}, nil
}

type pgxTxRunner struct {
	tx pgx.Tx
}

type pgxRowsRunner struct {
	rows pgx.Rows
}

func (r pgxRowsRunner) Next() bool {
	return r.rows.Next()
}

func (r pgxRowsRunner) Scan(dest ...any) error {
	return r.rows.Scan(dest...)
}

func (r pgxRowsRunner) Err() error {
	return r.rows.Err()
}

func (r pgxRowsRunner) Close() {
	r.rows.Close()
}

func (r pgxTxRunner) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return r.tx.Exec(ctx, sql, arguments...)
}

func (r pgxTxRunner) Commit(ctx context.Context) error {
	return r.tx.Commit(ctx)
}

func (r pgxTxRunner) Rollback(ctx context.Context) error {
	return r.tx.Rollback(ctx)
}

type PostgresRepository struct {
	db dbRunner
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: pgxPoolRunner{pool: pool}}
}

func newRepository(db dbRunner) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, address TestMailAddress) (*TestMailAddress, error) {
	const query = `
		INSERT INTO test_mail_addresses (owner_user_id, email)
		VALUES ($1, $2)
		RETURNING id, owner_user_id, email, created_at, updated_at, deleted_at
	`

	var created TestMailAddress
	err := r.db.QueryRow(ctx, query, address.OwnerUserID, address.Email).Scan(
		&created.ID,
		&created.OwnerUserID,
		&created.Email,
		&created.CreatedAt,
		&created.UpdatedAt,
		&created.DeletedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrDuplicateEmail
		}
		return nil, fmt.Errorf("create test mail address: %w", err)
	}

	return &created, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int64) (*TestMailAddress, error) {
	const query = `
		SELECT id, owner_user_id, email, created_at, updated_at, deleted_at
		FROM test_mail_addresses
		WHERE id = $1 AND deleted_at IS NULL
	`

	var address TestMailAddress
	err := r.db.QueryRow(ctx, query, id).Scan(
		&address.ID,
		&address.OwnerUserID,
		&address.Email,
		&address.CreatedAt,
		&address.UpdatedAt,
		&address.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTestMailAddressNotFound
		}
		return nil, fmt.Errorf("get test mail address by id: %w", err)
	}

	return &address, nil
}

func (r *PostgresRepository) GetByEmail(ctx context.Context, email string) (*TestMailAddress, error) {
	const query = `
		SELECT id, owner_user_id, email, created_at, updated_at, deleted_at
		FROM test_mail_addresses
		WHERE email = $1
	`

	var address TestMailAddress
	err := r.db.QueryRow(ctx, query, email).Scan(
		&address.ID,
		&address.OwnerUserID,
		&address.Email,
		&address.CreatedAt,
		&address.UpdatedAt,
		&address.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTestMailAddressNotFound
		}
		return nil, fmt.Errorf("get test mail address by email: %w", err)
	}

	return &address, nil
}

func (r *PostgresRepository) ListByOwner(ctx context.Context, ownerUserID int64) ([]TestMailAddress, error) {
	const query = `
		SELECT id, owner_user_id, email, created_at, updated_at, deleted_at
		FROM test_mail_addresses
		WHERE owner_user_id = $1 AND deleted_at IS NULL
		ORDER BY id ASC
	`

	rows, err := r.db.Query(ctx, query, ownerUserID)
	if err != nil {
		return nil, fmt.Errorf("query test mail addresses by owner: %w", err)
	}
	defer rows.Close()

	addresses := make([]TestMailAddress, 0)
	for rows.Next() {
		var address TestMailAddress
		err := rows.Scan(
			&address.ID,
			&address.OwnerUserID,
			&address.Email,
			&address.CreatedAt,
			&address.UpdatedAt,
			&address.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("list test mail addresses by owner: %w", err)
		}
		addresses = append(addresses, address)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate test mail addresses by owner: %w", err)
	}

	return addresses, nil
}

func (r *PostgresRepository) Update(ctx context.Context, address TestMailAddress) (*TestMailAddress, error) {
	const query = `
		UPDATE test_mail_addresses
		SET email = $2, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, owner_user_id, email, created_at, updated_at, deleted_at
	`

	var updated TestMailAddress
	err := r.db.QueryRow(ctx, query, address.ID, address.Email).Scan(
		&updated.ID,
		&updated.OwnerUserID,
		&updated.Email,
		&updated.CreatedAt,
		&updated.UpdatedAt,
		&updated.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTestMailAddressNotFound
		}
		if isUniqueViolation(err) {
			return nil, ErrDuplicateEmail
		}
		return nil, fmt.Errorf("update test mail address: %w", err)
	}

	return &updated, nil
}

func (r *PostgresRepository) SoftDelete(ctx context.Context, id int64, deletedAt time.Time) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin soft delete test mail address tx: %w", err)
	}
	defer tx.Rollback(ctx)

	const query = `
		UPDATE test_mail_addresses
		SET deleted_at = $2, updated_at = $2
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := tx.Exec(ctx, query, id, deletedAt)
	if err != nil {
		return fmt.Errorf("soft delete test mail address: %w", err)
	}
	if result.RowsAffected() != 1 {
		return ErrTestMailAddressNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit soft delete test mail address tx: %w", err)
	}

	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
