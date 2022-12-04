package entity

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type BaseEntity struct {
	ID        int64            `db:"id" json:"id"`
	CreatedAt pgtype.Timestamp `db:"created_at" json:"created_at"`
	UpdatedAt pgtype.Timestamp `db:"updated_at" json:"updated_at"`
}

// GetID returns the user ID.
func (b BaseEntity) GetID() int64 {
	return b.ID
}
