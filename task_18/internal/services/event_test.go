package services

import (
	"context"
	"testing"
	"time"

	"github.com/M-kos/wb_level2/task_18/internal/domains"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRepo struct {
	events []*domains.Event
}

func (m *mockRepo) Event(_ context.Context, id int) (*domains.Event, error) {
	for _, e := range m.events {
		if e.ID == id {
			return e, nil
		}
	}
	return nil, domains.ErrEventNotFound
}

func (m *mockRepo) List(_ context.Context, userId int, from, to time.Time) ([]*domains.Event, error) {
	var result []*domains.Event
	for _, e := range m.events {
		if e.UserID == userId && !e.Date.Before(from) && e.Date.Before(to) {
			result = append(result, e)
		}
	}
	return result, nil
}

func (m *mockRepo) Create(_ context.Context, e *domains.Event) (*domains.Event, error) {
	e.ID = len(m.events) + 1
	m.events = append(m.events, e)
	return e, nil
}

func (m *mockRepo) Update(_ context.Context, e *domains.Event) (*domains.Event, error) {
	for i, ev := range m.events {
		if ev.ID == e.ID {
			m.events[i] = e
			return e, nil
		}
	}
	return nil, domains.ErrEventNotFound
}

func (m *mockRepo) Delete(_ context.Context, id int) error {
	for i, e := range m.events {
		if e.ID == id {
			m.events = append(m.events[:i], m.events[i+1:]...)
			return nil
		}
	}
	return domains.ErrEventNotFound
}

func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func newEvent(id, userId int, d time.Time, title string) *domains.Event {
	return &domains.Event{ID: id, UserID: userId, Date: d, Title: title}
}

func setupService() (*EventService, *mockRepo) {
	repo := &mockRepo{
		events: []*domains.Event{
			newEvent(1, 1, date(2026, 3, 11), "monday"),
			newEvent(2, 1, date(2026, 3, 12), "tuesday"),
			newEvent(3, 1, date(2026, 3, 17), "sunday"),
			newEvent(4, 1, date(2026, 3, 25), "late march"),
			newEvent(5, 1, date(2026, 4, 1), "april"),
			newEvent(6, 2, date(2026, 3, 11), "other user"),
		},
	}
	return NewEventService(repo), repo
}

func TestEventsForDay(t *testing.T) {
	svc, _ := setupService()
	ctx := context.Background()

	events, err := svc.EventsForDay(ctx, 1, date(2026, 3, 11))
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, "monday", events[0].Title)
}

func TestEventsForDay_NoEvents(t *testing.T) {
	svc, _ := setupService()
	ctx := context.Background()

	events, err := svc.EventsForDay(ctx, 1, date(2026, 3, 13))
	require.NoError(t, err)
	assert.Empty(t, events)
}

func TestEventsForDay_FiltersByUserId(t *testing.T) {
	svc, _ := setupService()
	ctx := context.Background()

	events, err := svc.EventsForDay(ctx, 2, date(2026, 3, 11))
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, "other user", events[0].Title)
}

func TestEventsForWeek(t *testing.T) {
	svc, _ := setupService()
	ctx := context.Background()

	events, err := svc.EventsForWeek(ctx, 1, date(2026, 3, 13))
	require.NoError(t, err)
	assert.Len(t, events, 3)
}

func TestEventsForWeek_Sunday(t *testing.T) {
	svc, _ := setupService()
	ctx := context.Background()

	events, err := svc.EventsForWeek(ctx, 1, date(2026, 3, 17))
	require.NoError(t, err)
	assert.Len(t, events, 3)
}

func TestEventsForMonth(t *testing.T) {
	svc, _ := setupService()
	ctx := context.Background()

	events, err := svc.EventsForMonth(ctx, 1, date(2026, 3, 15))
	require.NoError(t, err)
	assert.Len(t, events, 4)
}

func TestCreate(t *testing.T) {
	svc, repo := setupService()
	ctx := context.Background()

	event, err := svc.Create(ctx, &domains.Event{
		UserID: 1,
		Title:  "new event",
		Date:   date(2024, 5, 1),
	})
	require.NoError(t, err)
	assert.NotZero(t, event.ID)
	assert.Len(t, repo.events, 7)
}

func TestUpdate(t *testing.T) {
	svc, _ := setupService()
	ctx := context.Background()

	updated, err := svc.Update(ctx, &domains.Event{
		ID:     1,
		UserID: 1,
		Title:  "updated monday",
		Date:   date(2026, 3, 11),
	})
	require.NoError(t, err)
	assert.Equal(t, "updated monday", updated.Title)
}

func TestUpdate_NotFound(t *testing.T) {
	svc, _ := setupService()
	ctx := context.Background()

	_, err := svc.Update(ctx, &domains.Event{ID: 999, Title: "nope"})
	assert.ErrorIs(t, err, domains.ErrEventNotFound)
}

func TestDelete(t *testing.T) {
	svc, repo := setupService()
	ctx := context.Background()

	err := svc.Delete(ctx, 1)
	require.NoError(t, err)
	assert.Len(t, repo.events, 5)
}

func TestDelete_NotFound(t *testing.T) {
	svc, _ := setupService()
	ctx := context.Background()

	err := svc.Delete(ctx, 999)
	assert.ErrorIs(t, err, domains.ErrEventNotFound)
}
