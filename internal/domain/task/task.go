package task

import (
	"time"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// Status представляет статус задачи
type Status string

const (
	StatusNew        Status = "new"        // Новая задача
	StatusInProgress Status = "in_progress" // Задача в процессе выполнения
	StatusDone       Status = "done"        // Завершенная задача
)

// PeriodicityType представляет тип периодичности задачи
type PeriodicityType string

const (
	PeriodicityNone      PeriodicityType = "none"           // Без периодичности
	PeriodicityDaily     PeriodicityType = "daily"          // Ежедневная периодичность
	PeriodicityMonthly   PeriodicityType = "monthly"        // Ежемесячная периодичность
	PeriodicitySpecific  PeriodicityType = "specific_dates" // Конкретные даты
	PeriodicityEvenOdd   PeriodicityType = "even_odd"       // Четные/нечетные дни
)

// EvenOddType представляет тип четности дней
type EvenOddType string

const (
	EvenOddEven EvenOddType = "even" // Четные дни
	EvenOddOdd  EvenOddType = "odd"  // Нечетные дни
)

// DateArray представляет массив дат для JSON сериализации
type DateArray []time.Time

// IntArray представляет массив целых чисел для JSON сериализации
type IntArray []int

// Task представляет задачу с настройками периодичности
type Task struct {
	ID                   int64            `json:"id"`                          // Уникальный идентификатор
	Title                string           `json:"title"`                        // Заголовок задачи
	Description          string           `json:"description"`                   // Описание задачи
	Status               Status           `json:"status"`                        // Статус задачи
	PeriodicityType      PeriodicityType  `json:"periodicity_type,omitempty"`    // Тип периодичности
	PeriodicityInterval  int              `json:"periodicity_interval,omitempty"` // Интервал дней для ежедневной периодичности
	PeriodicityDays      IntArray         `json:"periodicity_days,omitempty"`    // Дни месяца для ежемесячной периодичности
	PeriodicityDates     DateArray        `json:"periodicity_dates,omitempty"`    // Конкретные даты для периодичности
	PeriodicityEvenOdd   EvenOddType      `json:"periodicity_even_odd,omitempty"` // Четные/нечетные дни
	NextOccurrence       *time.Time       `json:"next_occurrence,omitempty"`      // Следующее выполнение задачи
	CreatedAt            time.Time         `json:"created_at"`                    // Дата создания
	UpdatedAt            time.Time         `json:"updated_at"`                    // Дата обновления
}

// Valid проверяет валидность статуса задачи
func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

// Valid проверяет валидность типа периодичности
func (pt PeriodicityType) Valid() bool {
	switch pt {
	case PeriodicityNone, PeriodicityDaily, PeriodicityMonthly, PeriodicitySpecific, PeriodicityEvenOdd:
		return true
	default:
		return false
	}
}

// Valid проверяет валидность типа четности дней
func (eot EvenOddType) Valid() bool {
	switch eot {
	case EvenOddEven, EvenOddOdd:
		return true
	default:
		return false
	}
}

// Value реализует driver.Valuer для DateArray (сериализация в JSON)
func (da DateArray) Value() (driver.Value, error) {
	return json.Marshal(da)
}

// Scan реализует sql.Scanner для DateArray (десериализация из JSON)
func (da *DateArray) Scan(value interface{}) error {
	if value == nil {
		*da = DateArray{}
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprintf("failed to unmarshal DateArray value: %v", value))
	}
	
	return json.Unmarshal(bytes, da)
}

// Value реализует driver.Valuer для IntArray (сериализация в JSON)
func (ia IntArray) Value() (driver.Value, error) {
	return json.Marshal(ia)
}

// Scan реализует sql.Scanner для IntArray (десериализация из JSON)
func (ia *IntArray) Scan(value interface{}) error {
	if value == nil {
		*ia = IntArray{}
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprintf("failed to unmarshal IntArray value: %v", value))
	}
	
	return json.Unmarshal(bytes, ia)
}
