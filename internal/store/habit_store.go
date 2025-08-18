package store

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Habit struct {
	ID uuid.UUID
	// UserID      uuid.UUID
	Name        string
	Description string
	Frequency   string
	TargetCount int
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type HabitEntry struct {
	ID         uuid.UUID
	HabitID    uuid.UUID
	Completion time.Time
	Note       string
	CreatedAt  time.Time
}

type HabitTags struct {
	ID        uuid.UUID
	HabitID   uuid.UUID
	TagID     uuid.UUID
	CreatedAt time.Time
}

type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{db: db}
}

type HabitStore interface {
	CreateHabit(*Habit) (*Habit, error)
	GetHabitByID(id uuid.UUID) (*Habit, error)
	// GetHabits() ([]*Habit, error)
}

func (pg *PostgresWorkoutStore) CreateHabit(habit *Habit) (*Habit, error) {
	habit.ID = uuid.New()

	// Start a database transaction.
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // Ensure that the transaction is rolled back if anything fails.

	query := `
		INSERT INTO habits (id, name, description, frequency, target_count, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	// Execute the query and scan the returned workout ID into the workout struct.
	err = tx.QueryRow(query, habit.ID, habit.Name, habit.Description, habit.Frequency, habit.TargetCount, habit.IsActive).Scan(&habit.ID)
	if err != nil {
		return nil, err
	}

	// Commit the transaction to save all changes to the database.
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return habit, nil
}

func (pg *PostgresWorkoutStore) GetHabitByID(id uuid.UUID) (*Habit, error) {
	habit := &Habit{}

	query := `
		SELECT id, name, description, frequency, target_count, is_active
		FROM habits
		WHERE id = $1`

	err := pg.db.QueryRow(query, id).Scan(&habit.ID, &habit.Name, &habit.Description, &habit.Frequency, &habit.TargetCount, &habit.IsActive)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return habit, nil
}

// func (pg *PostgresWorkoutStore) GetHabits() []*Habit, error {

// }

func (pg *PostgresWorkoutStore) UpdateHabit(habit *Habit) error {
	return nil
}

func (pg *PostgresWorkoutStore) DeleteHabit(id int64) error {
	return nil
}

func (pg *PostgresWorkoutStore) GetHabitOwner(habitID int64) (int, error) {
	return 0, nil
}
