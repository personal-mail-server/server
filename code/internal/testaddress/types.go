package testaddress

import (
	"context"
	"time"
)

type TestMailAddress struct {
	ID          int64
	OwnerUserID int64
	Email       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type Repository interface {
	Create(ctx context.Context, address TestMailAddress) (*TestMailAddress, error)
	GetByID(ctx context.Context, id int64) (*TestMailAddress, error)
	GetByEmail(ctx context.Context, email string) (*TestMailAddress, error)
	ListByOwner(ctx context.Context, ownerUserID int64) ([]TestMailAddress, error)
	Update(ctx context.Context, address TestMailAddress) (*TestMailAddress, error)
	SoftDelete(ctx context.Context, id int64, deletedAt time.Time) error
}
