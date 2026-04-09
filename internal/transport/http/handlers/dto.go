package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type taskMutationDTO struct {
	Title                string                    `json:"title"`
	Description          string                    `json:"description"`
	Status               taskdomain.Status         `json:"status"`
	PeriodicityType      taskdomain.PeriodicityType  `json:"periodicity_type,omitempty"`
	PeriodicityInterval  int                       `json:"periodicity_interval,omitempty"`
	PeriodicityDays      taskdomain.IntArray       `json:"periodicity_days,omitempty"`
	PeriodicityDates     taskdomain.DateArray      `json:"periodicity_dates,omitempty"`
	PeriodicityEvenOdd   taskdomain.EvenOddType    `json:"periodicity_even_odd,omitempty"`
}

type taskDTO struct {
	ID                   int64                     `json:"id"`
	Title                string                    `json:"title"`
	Description          string                    `json:"description"`
	Status               taskdomain.Status         `json:"status"`
	PeriodicityType      taskdomain.PeriodicityType  `json:"periodicity_type,omitempty"`
	PeriodicityInterval  int                       `json:"periodicity_interval,omitempty"`
	PeriodicityDays      taskdomain.IntArray       `json:"periodicity_days,omitempty"`
	PeriodicityDates     taskdomain.DateArray      `json:"periodicity_dates,omitempty"`
	PeriodicityEvenOdd   taskdomain.EvenOddType    `json:"periodicity_even_odd,omitempty"`
	NextOccurrence       *time.Time                `json:"next_occurrence,omitempty"`
	CreatedAt            time.Time                 `json:"created_at"`
	UpdatedAt            time.Time                 `json:"updated_at"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	return taskDTO{
		ID:                   task.ID,
		Title:                task.Title,
		Description:          task.Description,
		Status:               task.Status,
		PeriodicityType:      task.PeriodicityType,
		PeriodicityInterval:  task.PeriodicityInterval,
		PeriodicityDays:      task.PeriodicityDays,
		PeriodicityDates:     task.PeriodicityDates,
		PeriodicityEvenOdd:   task.PeriodicityEvenOdd,
		NextOccurrence:       task.NextOccurrence,
		CreatedAt:            task.CreatedAt,
		UpdatedAt:            task.UpdatedAt,
	}
}
