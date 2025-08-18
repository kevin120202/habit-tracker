package store

import (
	"database/sql"
	"errors"
	"strings"
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
	GetTags() ([]*Tag, error)
}

func (pg *PostgresTagStore) GetTags() ([]*Tag, error) {
	query := `
		SELECT id, name, color, created_at, updated_at
		FROM tags
		ORDER BY name`

	rows, err := pg.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*Tag
	for rows.Next() {
		tag := &Tag{}
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Color, &tag.CreatedAt, &tag.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (pg *PostgresTagStore) CreateTag(tag *Tag) (*Tag, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	tag.ID = uuid.New()
	query := `
		INSERT INTO tags (id, name, color)
		VALUES ($1, $2, $3)
		RETURNING id`

	err = tx.QueryRow(query, tag.ID, tag.Name, tag.Color).Scan(&tag.ID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, errors.New("tag with this name already exists")
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return tag, nil
}
