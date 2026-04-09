package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		Title:                normalized.Title,
		Description:          normalized.Description,
		Status:               normalized.Status,
		PeriodicityType:      normalized.PeriodicityType,
		PeriodicityInterval:  normalized.PeriodicityInterval,
		PeriodicityDays:      normalized.PeriodicityDays,
		PeriodicityDates:     normalized.PeriodicityDates,
		PeriodicityEvenOdd:   normalized.PeriodicityEvenOdd,
	}
	now := s.now()
	model.CreatedAt = now
	model.UpdatedAt = now

	// Calculate next occurrence for recurring tasks
	if model.PeriodicityType != "" && model.PeriodicityType != taskdomain.PeriodicityNone {
		nextOccurrence, err := calculateNextOccurrence(model, now)
		if err != nil {
			return nil, err
		}
		model.NextOccurrence = &nextOccurrence
	}

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		ID:                   id,
		Title:                normalized.Title,
		Description:          normalized.Description,
		Status:               normalized.Status,
		PeriodicityType:      normalized.PeriodicityType,
		PeriodicityInterval:  normalized.PeriodicityInterval,
		PeriodicityDays:      normalized.PeriodicityDays,
		PeriodicityDates:     normalized.PeriodicityDates,
		PeriodicityEvenOdd:   normalized.PeriodicityEvenOdd,
		UpdatedAt:            s.now(),
	}

	// Calculate next occurrence for recurring tasks
	if model.PeriodicityType != "" && model.PeriodicityType != taskdomain.PeriodicityNone {
		nextOccurrence, err := calculateNextOccurrence(model, s.now())
		if err != nil {
			return nil, err
		}
		model.NextOccurrence = &nextOccurrence
	}

	updated, err := s.repo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	return s.repo.List(ctx)
}

func validateCreateInput(input CreateInput) (CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return CreateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}

	if !input.Status.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	// Validate periodicity settings
	if err := validatePeriodicity(input.PeriodicityType, input.PeriodicityInterval, 
		input.PeriodicityDays, input.PeriodicityDates, input.PeriodicityEvenOdd); err != nil {
		return CreateInput{}, err
	}

	return input, nil
}

// validatePeriodicity validates periodicity settings
func validatePeriodicity(
	periodicityType taskdomain.PeriodicityType,
	periodicityInterval int,
	periodicityDays taskdomain.IntArray,
	periodicityDates taskdomain.DateArray,
	periodicityEvenOdd taskdomain.EvenOddType,
) error {
	// Validate periodicity type
	if !periodicityType.Valid() && periodicityType != "" {
		return fmt.Errorf("%w: invalid periodicity type", ErrInvalidInput)
	}

	// If no periodicity type is specified, no further validation needed
	if periodicityType == "" || periodicityType == taskdomain.PeriodicityNone {
		return nil
	}

	// Validate based on periodicity type
	switch periodicityType {
	case taskdomain.PeriodicityDaily:
		if periodicityInterval <= 0 {
			return fmt.Errorf("%w: daily periodicity requires positive interval", ErrInvalidInput)
		}

	case taskdomain.PeriodicityMonthly:
		if len(periodicityDays) == 0 {
			return fmt.Errorf("%w: monthly periodicity requires at least one day", ErrInvalidInput)
		}
		for _, day := range periodicityDays {
			if day < 1 || day > 31 {
				return fmt.Errorf("%w: monthly day must be between 1 and 31", ErrInvalidInput)
			}
		}

	case taskdomain.PeriodicitySpecific:
		if len(periodicityDates) == 0 {
			return fmt.Errorf("%w: specific dates periodicity requires at least one date", ErrInvalidInput)
		}
		now := time.Now()
		for _, date := range periodicityDates {
			if date.Before(now) {
				return fmt.Errorf("%w: specific dates must be in the future", ErrInvalidInput)
			}
		}

	case taskdomain.PeriodicityEvenOdd:
		if !periodicityEvenOdd.Valid() {
			return fmt.Errorf("%w: even/odd periodicity requires valid type (even or odd)", ErrInvalidInput)
		}

	default:
		return fmt.Errorf("%w: unsupported periodicity type", ErrInvalidInput)
	}

	return nil
}

// calculateNextOccurrence calculates the next occurrence date for a recurring task
func calculateNextOccurrence(task *taskdomain.Task, currentTime time.Time) (time.Time, error) {
	switch task.PeriodicityType {
	case taskdomain.PeriodicityDaily:
		// Add n days to current date
		return currentTime.AddDate(0, 0, task.PeriodicityInterval), nil

	case taskdomain.PeriodicityMonthly:
		// Find the next occurrence in the current or next month
		now := currentTime
		currentDay := now.Day()
		
		// Try to find the next day in the current month
		for _, day := range task.PeriodicityDays {
			if day > currentDay {
				// This day exists in the current month
				return time.Date(now.Year(), now.Month(), day, 0, 0, 0, 0, time.UTC), nil
			}
		}
		
		// If no day found in current month, use first day of next month
		firstDay := task.PeriodicityDays[0]
		nextMonth := now.AddDate(0, 1, 0)
		return time.Date(nextMonth.Year(), nextMonth.Month(), firstDay, 0, 0, 0, 0, time.UTC), nil

	case taskdomain.PeriodicitySpecific:
		// Return the first specific date that is in the future
		for _, date := range task.PeriodicityDates {
			if date.After(currentTime) {
				return date, nil
			}
		}
		return time.Time{}, fmt.Errorf("no future dates available for specific dates periodicity")

	case taskdomain.PeriodicityEvenOdd:
		// Find the next even or odd day
		now := currentTime
		for i := 0; i < 31; i++ { // Limit to 31 days to avoid infinite loop
			nextDay := now.AddDate(0, 0, i)
			day := nextDay.Day()
			
			// Check if day matches even/odd requirement
			if (task.PeriodicityEvenOdd == taskdomain.EvenOddEven && day%2 == 0) ||
			   (task.PeriodicityEvenOdd == taskdomain.EvenOddOdd && day%2 == 1) {
				return nextDay, nil
			}
		}
		return time.Time{}, fmt.Errorf("could not find suitable even/odd day")

	default:
		return time.Time{}, fmt.Errorf("unsupported periodicity type: %s", task.PeriodicityType)
	}
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return UpdateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if !input.Status.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	// Validate periodicity settings
	if err := validatePeriodicity(input.PeriodicityType, input.PeriodicityInterval, 
		input.PeriodicityDays, input.PeriodicityDates, input.PeriodicityEvenOdd); err != nil {
		return UpdateInput{}, err
	}

	return input, nil
}
