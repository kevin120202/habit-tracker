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
}

type HabitTags struct {
	ID      uuid.UUID
	HabitID uuid.UUID
	TagID   uuid.UUID
}

type PostgresHabitStore struct {
	db *sql.DB
}

func NewPostgresHabitStore(db *sql.DB) *PostgresHabitStore {
	return &PostgresHabitStore{db: db}
}

type HabitStore interface {
	CreateHabit(*Habit) (*Habit, error)
	GetHabitByID(id uuid.UUID) (*Habit, error)
	GetHabits() ([]*Habit, error)
	UpdateHabit(*Habit) error
	DeleteHabit(id uuid.UUID) error
	LogHabit(*HabitEntry) (*HabitEntry, error)
}

func (pg *PostgresHabitStore) CreateHabit(habit *Habit) (*Habit, error) {
	habit.ID = uuid.New()

	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO habits (id, name, description, frequency, target_count, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	err = tx.QueryRow(query, habit.ID, habit.Name, habit.Description, habit.Frequency, habit.TargetCount, habit.IsActive).Scan(&habit.ID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return habit, nil
}

func (pg *PostgresHabitStore) GetHabitByID(id uuid.UUID) (*Habit, error) {
	habit := &Habit{}

	query := `
		SELECT id, name, description, frequency, target_count, is_active, created_at, updated_at
		FROM habits
		WHERE id = $1`

	err := pg.db.QueryRow(query, id).Scan(&habit.ID, &habit.Name, &habit.Description, &habit.Frequency, &habit.TargetCount, &habit.IsActive, &habit.CreatedAt, &habit.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return habit, nil
}

func (pg *PostgresHabitStore) GetHabits() ([]*Habit, error) {
	query := `
		SELECT id, name, description, frequency, target_count, is_active, created_at, updated_at
		FROM habits
		ORDER BY name`

	rows, err := pg.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var habits []*Habit
	for rows.Next() {
		habit := &Habit{}
		err := rows.Scan(&habit.ID, &habit.Name, &habit.Description, &habit.Frequency, &habit.TargetCount, &habit.IsActive, &habit.CreatedAt, &habit.UpdatedAt)
		if err != nil {
			return nil, err
		}
		habits = append(habits, habit)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return habits, nil
}

func (pg *PostgresHabitStore) UpdateHabit(habit *Habit) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE habits
		SET name = $1, description = $2, frequency = $3, target_count = $4, is_active = $5, updated_at = $6
		WHERE id = $7
	`

	result, err := tx.Exec(query, habit.Name, habit.Description, habit.Frequency, habit.TargetCount, habit.IsActive, time.Now(), habit.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return tx.Commit()
}

func (pg *PostgresHabitStore) DeleteHabit(id uuid.UUID) error {
	query := `
		DELETE from habits
		WHERE id = $1`

	result, err := pg.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// func (pg *PostgresHabitStore) GetHabitOwner(habitID int64) (int, error) {
// 	return 0, nil
// }

func (pg *PostgresHabitStore) LogHabit(habitEntry *HabitEntry) (*HabitEntry, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	habitEntry.ID = uuid.New()
	query := `
		INSERT INTO habit_entries (id, habit_id, completion_date, note)
		VALUES ($1, $2, $3, $4)
		RETURNING id, habit_id, completion_date, note`

	err = tx.QueryRow(query, habitEntry.ID, habitEntry.HabitID, time.Now(), habitEntry.Note).Scan(&habitEntry.ID, &habitEntry.HabitID, &habitEntry.Completion, &habitEntry.Note)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return habitEntry, nil
}
