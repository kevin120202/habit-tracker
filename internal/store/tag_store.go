package store

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	ID        uuid.UUID
	Name      string
	Color     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PostgresTagStore struct {
	db *sql.DB
}

func NewPostgresTagStore(db *sql.DB) *PostgresTagStore {
	return &PostgresTagStore{db: db}
}

type TagStore interface {
	CreateTag(*Tag) (*Tag, error)
}
