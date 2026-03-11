package dto

import (
	"time"

	"github.com/M-kos/wb_level2/task_18/internal/domains"
)

type EventDto struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

func EventDtoFromDomain(event *domains.Event) *EventDto {
	return &EventDto{
		ID:          event.ID,
		UserID:      event.UserID,
		Title:       event.Title,
		Description: event.Description,
		Date:        event.Date.Format(time.DateOnly),
	}
}
