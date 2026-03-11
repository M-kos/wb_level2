package dto

import (
	"time"

	"github.com/M-kos/wb_level2/task_18/internal/domains"
	"github.com/go-playground/validator/v10"
)

type CreateEvent struct {
	UserId      int    `json:"user_id" validate:"required"`
	Title       string `json:"title" validate:"required,min=2"`
	Description string `json:"description" validate:"omitempty,min=1"`
	Date        string `json:"date" validate:"required,datetime=2006-01-02"`
}

func (ce *CreateEvent) Validate() error {
	validate := validator.New()

	return validate.Struct(ce)
}

func (ce *CreateEvent) ToDomain() *domains.Event {
	date, err := time.Parse(time.DateOnly, ce.Date)

	if err != nil {
		date = time.Time{}
	}

	return &domains.Event{
		UserID:      ce.UserId,
		Title:       ce.Title,
		Description: ce.Description,
		Date:        date,
	}
}
