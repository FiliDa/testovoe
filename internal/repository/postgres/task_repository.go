package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	const query = `
		INSERT INTO tasks (
			title, description, status, 
			periodicity_type, periodicity_interval, periodicity_days_of_month, 
			periodicity_specific_dates, periodicity_even_odd, next_occurrence,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, title, description, status, 
			periodicity_type, periodicity_interval, periodicity_days_of_month, 
			periodicity_specific_dates, periodicity_even_odd, next_occurrence,
			created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query, 
		task.Title, task.Description, task.Status,
		task.PeriodicityType, task.PeriodicityInterval, task.PeriodicityDays,
		task.PeriodicityDates, task.PeriodicityEvenOdd, task.NextOccurrence,
		task.CreatedAt, task.UpdatedAt,
	)
	created, err := scanTask(row)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	const query = `
		SELECT 
			id, title, description, status, 
			periodicity_type, periodicity_interval, periodicity_days_of_month, 
			periodicity_specific_dates, periodicity_even_odd, next_occurrence,
			created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	found, err := scanTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrNotFound
		}

		return nil, err
	}

	return found, nil
}

func (r *Repository) Update(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	const query = `
		UPDATE tasks
		SET title = $1,
			description = $2,
			status = $3,
			periodicity_type = $4,
			periodicity_interval = $5,
			periodicity_days_of_month = $6,
			periodicity_specific_dates = $7,
			periodicity_even_odd = $8,
			next_occurrence = $9,
			updated_at = $10
		WHERE id = $11
		RETURNING id, title, description, status, 
			periodicity_type, periodicity_interval, periodicity_days_of_month, 
			periodicity_specific_dates, periodicity_even_odd, next_occurrence,
			created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query, 
		task.Title, task.Description, task.Status,
		task.PeriodicityType, task.PeriodicityInterval, task.PeriodicityDays,
		task.PeriodicityDates, task.PeriodicityEvenOdd, task.NextOccurrence,
		task.UpdatedAt, task.ID,
	)
	updated, err := scanTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrNotFound
		}

		return nil, err
	}

	return updated, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM tasks WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return taskdomain.ErrNotFound
	}

	return nil
}

func (r *Repository) List(ctx context.Context) ([]taskdomain.Task, error) {
	const query = `
		SELECT 
			id, title, description, status, 
			periodicity_type, periodicity_interval, periodicity_days_of_month, 
			periodicity_specific_dates, periodicity_even_odd, next_occurrence,
			created_at, updated_at
		FROM tasks
		ORDER BY id DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]taskdomain.Task, 0)
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, *task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

type taskScanner interface {
	Scan(dest ...any) error
}

func scanTask(scanner taskScanner) (*taskdomain.Task, error) {
	var (
		task                     taskdomain.Task
		status                   string
		periodicityType          *string
		periodicityInterval      *int
		periodicityDays          taskdomain.IntArray
		periodicityDates         taskdomain.DateArray
		periodicityEvenOdd       *string
		nextOccurrence           *time.Time
	)

	if err := scanner.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&status,
		&periodicityType,
		&periodicityInterval,
		&periodicityDays,
		&periodicityDates,
		&periodicityEvenOdd,
		&nextOccurrence,
		&task.CreatedAt,
		&task.UpdatedAt,
	); err != nil {
		return nil, err
	}

	task.Status = taskdomain.Status(status)

	// Set periodicity fields
	if periodicityType != nil {
		task.PeriodicityType = taskdomain.PeriodicityType(*periodicityType)
	}
	if periodicityInterval != nil {
		task.PeriodicityInterval = *periodicityInterval
	}
	task.PeriodicityDays = periodicityDays
	task.PeriodicityDates = periodicityDates
	if periodicityEvenOdd != nil {
		task.PeriodicityEvenOdd = taskdomain.EvenOddType(*periodicityEvenOdd)
	}
	task.NextOccurrence = nextOccurrence

	return &task, nil
}
