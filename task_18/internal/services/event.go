package services

import (
	"context"
	"time"

	"github.com/M-kos/wb_level2/task_18/internal/domains"
)

type EventRepository interface {
	Event(ctx context.Context, id int) (*domains.Event, error)
	List(ctx context.Context, userId int, from, to time.Time) ([]*domains.Event, error)
	Create(ctx context.Context, newEvent *domains.Event) (*domains.Event, error)
	Update(ctx context.Context, event *domains.Event) (*domains.Event, error)
	Delete(ctx context.Context, id int) error
}

type EventService struct {
	repo EventRepository
}

func NewEventService(repo EventRepository) *EventService {
	return &EventService{
		repo: repo,
	}
}

func (es *EventService) EventsForDay(ctx context.Context, userId int, date time.Time) ([]*domains.Event, error) {
	from := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	to := from.Add(24 * time.Hour)

	return es.repo.List(ctx, userId, from, to)
}

func (es *EventService) EventsForWeek(ctx context.Context, userId int, date time.Time) ([]*domains.Event, error) {
	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	from := time.Date(date.Year(), date.Month(), date.Day()-weekday+1, 0, 0, 0, 0, date.Location())
	to := from.Add(7 * 24 * time.Hour)

	return es.repo.List(ctx, userId, from, to)
}

func (es *EventService) EventsForMonth(ctx context.Context, userId int, date time.Time) ([]*domains.Event, error) {
	from := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	to := from.AddDate(0, 1, 0)

	return es.repo.List(ctx, userId, from, to)
}

func (es *EventService) Create(ctx context.Context, newEvent *domains.Event) (*domains.Event, error) {
	return es.repo.Create(ctx, newEvent)
}

func (es *EventService) Update(ctx context.Context, event *domains.Event) (*domains.Event, error) {
	return es.repo.Update(ctx, event)
}

func (es *EventService) Delete(ctx context.Context, eventId int) error {
	return es.repo.Delete(ctx, eventId)
}
