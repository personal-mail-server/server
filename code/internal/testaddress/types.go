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

type CreateRequest struct {
	Email string `json:"email"`
}

type Response struct {
	ID          int64      `json:"id"`
	OwnerUserID int64      `json:"ownerUserId"`
	Email       string     `json:"email"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
}

type ListResponse struct {
	Addresses []*Response `json:"addresses"`
}

type Repository interface {
	Create(ctx context.Context, address TestMailAddress) (*TestMailAddress, error)
	GetByID(ctx context.Context, id int64) (*TestMailAddress, error)
	GetByEmail(ctx context.Context, email string) (*TestMailAddress, error)
	ListByOwner(ctx context.Context, ownerUserID int64) ([]TestMailAddress, error)
	Update(ctx context.Context, address TestMailAddress) (*TestMailAddress, error)
	SoftDelete(ctx context.Context, id int64, deletedAt time.Time) error
}

func NewResponse(address *TestMailAddress) *Response {
	if address == nil {
		return nil
	}
	return &Response{
		ID:          address.ID,
		OwnerUserID: address.OwnerUserID,
		Email:       address.Email,
		CreatedAt:   address.CreatedAt,
		UpdatedAt:   address.UpdatedAt,
		DeletedAt:   address.DeletedAt,
	}
}

func NewListResponse(addresses []TestMailAddress) *ListResponse {
	items := make([]*Response, 0, len(addresses))
	for i := range addresses {
		address := addresses[i]
		items = append(items, NewResponse(&address))
	}
	return &ListResponse{Addresses: items}
}
