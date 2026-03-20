package testaddress

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type fakeRow struct {
	values []any
	err    error
}

func (r fakeRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	if len(dest) != len(r.values) {
		return errors.New("scan length mismatch")
	}
	for i := range dest {
		switch d := dest[i].(type) {
		case *int:
			*d = r.values[i].(int)
		case *int64:
			*d = r.values[i].(int64)
		case *string:
			*d = r.values[i].(string)
		case *time.Time:
			*d = r.values[i].(time.Time)
		case **time.Time:
			*d = r.values[i].(*time.Time)
		default:
			return errors.New("unsupported scan destination")
		}
	}
	return nil
}

type fakeRows struct {
	rows []fakeRow
	idx  int
	err  error
}

func (r *fakeRows) Next() bool {
	return r.idx < len(r.rows)
}

func (r *fakeRows) Scan(dest ...any) error {
	if r.idx >= len(r.rows) {
		return errors.New("unexpected rows scan")
	}
	row := r.rows[r.idx]
	r.idx++
	return row.Scan(dest...)
}

func (r *fakeRows) Err() error {
	return r.err
}

func (r *fakeRows) Close() {}

type fakeDB struct {
	queryRows []fakeRow
	queryIdx  int
	rows      rowsScanner
	queryErr  error
	beginTx   dbTx
	beginErr  error
	execTag   pgconn.CommandTag
	execErr   error
}

func (f *fakeDB) QueryRow(_ context.Context, _ string, _ ...any) queryRowScanner {
	if f.queryIdx >= len(f.queryRows) {
		return fakeRow{err: errors.New("unexpected query row call")}
	}
	row := f.queryRows[f.queryIdx]
	f.queryIdx++
	return row
}

func (f *fakeDB) Query(_ context.Context, _ string, _ ...any) (rowsScanner, error) {
	if f.queryErr != nil {
		return nil, f.queryErr
	}
	if f.rows == nil {
		return &fakeRows{}, nil
	}
	return f.rows, nil
}

func (f *fakeDB) Exec(_ context.Context, _ string, _ ...any) (pgconn.CommandTag, error) {
	return f.execTag, f.execErr
}

func (f *fakeDB) Begin(_ context.Context) (dbTx, error) {
	if f.beginErr != nil {
		return nil, f.beginErr
	}
	return f.beginTx, nil
}

type fakeTx struct {
	execTag     pgconn.CommandTag
	execErr     error
	commitErr   error
	rollbackErr error
	committed   bool
	rolledBack  bool
}

func (f *fakeTx) Exec(_ context.Context, _ string, _ ...any) (pgconn.CommandTag, error) {
	return f.execTag, f.execErr
}

func (f *fakeTx) Commit(_ context.Context) error {
	f.committed = true
	return f.commitErr
}

func (f *fakeTx) Rollback(_ context.Context) error {
	f.rolledBack = true
	return f.rollbackErr
}

func TestCreateReturnsStoredAddress(t *testing.T) {
	now := time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)
	repo := newRepository(&fakeDB{queryRows: []fakeRow{{values: []any{int64(1), int64(7), "alpha@test.local", now, now, (*time.Time)(nil)}}}})

	created, err := repo.Create(context.Background(), TestMailAddress{OwnerUserID: 7, Email: "alpha@test.local"})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if created.ID != 1 || created.OwnerUserID != 7 || created.Email != "alpha@test.local" {
		t.Fatalf("unexpected created address: %+v", created)
	}
}

func TestCreateReturnsDuplicateEmailOnUniqueViolation(t *testing.T) {
	repo := newRepository(&fakeDB{queryRows: []fakeRow{{err: &pgconn.PgError{Code: "23505"}}}})

	_, err := repo.Create(context.Background(), TestMailAddress{OwnerUserID: 7, Email: "dup@test.local"})
	if !errors.Is(err, ErrDuplicateEmail) {
		t.Fatalf("expected ErrDuplicateEmail, got %v", err)
	}
}

func TestGetByIDReturnsNotFound(t *testing.T) {
	repo := newRepository(&fakeDB{queryRows: []fakeRow{{err: pgx.ErrNoRows}}})

	_, err := repo.GetByID(context.Background(), 1)
	if !errors.Is(err, ErrTestMailAddressNotFound) {
		t.Fatalf("expected ErrTestMailAddressNotFound, got %v", err)
	}
}

func TestGetByEmailReturnsStoredAddress(t *testing.T) {
	now := time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)
	repo := newRepository(&fakeDB{queryRows: []fakeRow{{values: []any{int64(4), int64(2), "beta@test.local", now, now, (*time.Time)(nil)}}}})

	address, err := repo.GetByEmail(context.Background(), "beta@test.local")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if address.Email != "beta@test.local" || address.OwnerUserID != 2 {
		t.Fatalf("unexpected address: %+v", address)
	}
}

func TestListByOwnerReturnsEmptyWhenNoActiveAddresses(t *testing.T) {
	repo := newRepository(&fakeDB{rows: &fakeRows{rows: []fakeRow{}}})

	addresses, err := repo.ListByOwner(context.Background(), 9)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(addresses) != 0 {
		t.Fatalf("expected empty list, got %+v", addresses)
	}
}

func TestListByOwnerReturnsActiveAddresses(t *testing.T) {
	now := time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)
	repo := newRepository(&fakeDB{rows: &fakeRows{rows: []fakeRow{
		{values: []any{int64(1), int64(3), "one@test.local", now, now, (*time.Time)(nil)}},
		{values: []any{int64(2), int64(3), "two@test.local", now, now, (*time.Time)(nil)}},
	}}})

	addresses, err := repo.ListByOwner(context.Background(), 3)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(addresses) != 2 {
		t.Fatalf("expected two addresses, got %+v", addresses)
	}
	if addresses[0].Email != "one@test.local" || addresses[1].Email != "two@test.local" {
		t.Fatalf("unexpected addresses: %+v", addresses)
	}
}

func TestUpdateReturnsDuplicateEmailOnUniqueViolation(t *testing.T) {
	repo := newRepository(&fakeDB{queryRows: []fakeRow{{err: &pgconn.PgError{Code: "23505"}}}})

	_, err := repo.Update(context.Background(), TestMailAddress{ID: 11, Email: "dup@test.local"})
	if !errors.Is(err, ErrDuplicateEmail) {
		t.Fatalf("expected ErrDuplicateEmail, got %v", err)
	}
}

func TestSoftDeleteMarksDeletedAndCommits(t *testing.T) {
	tx := &fakeTx{execTag: pgconn.NewCommandTag("UPDATE 1")}
	repo := newRepository(&fakeDB{beginTx: tx})

	err := repo.SoftDelete(context.Background(), 12, time.Date(2026, 3, 20, 1, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !tx.committed {
		t.Fatalf("expected transaction commit")
	}
}

func TestSoftDeleteReturnsNotFoundWhenNoRowUpdated(t *testing.T) {
	tx := &fakeTx{execTag: pgconn.NewCommandTag("UPDATE 0")}
	repo := newRepository(&fakeDB{beginTx: tx})

	err := repo.SoftDelete(context.Background(), 77, time.Now().UTC())
	if !errors.Is(err, ErrTestMailAddressNotFound) {
		t.Fatalf("expected ErrTestMailAddressNotFound, got %v", err)
	}
}
