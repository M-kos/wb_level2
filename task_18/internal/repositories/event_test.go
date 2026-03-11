package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/M-kos/wb_level2/task_18/internal/domains"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func TestCreate(t *testing.T) {
	repo := NewEventRepository()
	ctx := context.Background()

	event, err := repo.Create(ctx, &domains.Event{
		UserID: 1,
		Title:  "test",
		Date:   date(2026, 3, 11),
	})
	require.NoError(t, err)
	assert.Equal(t, 1, event.ID)
}

func TestEvent(t *testing.T) {
	repo := NewEventRepository()
	ctx := context.Background()

	repo.Create(ctx, &domains.Event{UserID: 1, Title: "test", Date: date(2026, 3, 11)})

	event, err := repo.Event(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, "test", event.Title)
}

func TestEvent_NotFound(t *testing.T) {
	repo := NewEventRepository()
	ctx := context.Background()

	_, err := repo.Event(ctx, 999)
	assert.ErrorIs(t, err, domains.ErrEventNotFound)
}

func TestList(t *testing.T) {
	repo := NewEventRepository()
	ctx := context.Background()

	repo.Create(ctx, &domains.Event{UserID: 1, Title: "mar 11", Date: date(2026, 3, 11)})
	repo.Create(ctx, &domains.Event{UserID: 1, Title: "mar 12", Date: date(2026, 3, 12)})
	repo.Create(ctx, &domains.Event{UserID: 1, Title: "mar 13", Date: date(2026, 3, 13)})
	repo.Create(ctx, &domains.Event{UserID: 2, Title: "other user", Date: date(2026, 3, 11)})

	events, err := repo.List(ctx, 1, date(2026, 3, 11), date(2026, 3, 12))
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, "mar 11", events[0].Title)
}

func TestUpdate(t *testing.T) {
	repo := NewEventRepository()
	ctx := context.Background()

	repo.Create(ctx, &domains.Event{UserID: 1, Title: "original", Date: date(2026, 3, 11)})

	updated, err := repo.Update(ctx, &domains.Event{
		ID:     1,
		UserID: 1,
		Title:  "modified",
		Date:   date(2026, 3, 11),
	})
	require.NoError(t, err)
	assert.Equal(t, "modified", updated.Title)

	event, _ := repo.Event(ctx, 1)
	assert.Equal(t, "modified", event.Title)
}

func TestUpdate_NotFound(t *testing.T) {
	repo := NewEventRepository()
	ctx := context.Background()

	_, err := repo.Update(ctx, &domains.Event{ID: 999, Title: "nope"})
	assert.ErrorIs(t, err, domains.ErrEventNotFound)
}

func TestDelete(t *testing.T) {
	repo := NewEventRepository()
	ctx := context.Background()

	repo.Create(ctx, &domains.Event{UserID: 1, Title: "to delete", Date: date(2026, 3, 11)})

	err := repo.Delete(ctx, 1)
	require.NoError(t, err)

	_, err = repo.Event(ctx, 1)
	assert.ErrorIs(t, err, domains.ErrEventNotFound)
}

func TestDelete_NotFound(t *testing.T) {
	repo := NewEventRepository()
	ctx := context.Background()

	err := repo.Delete(ctx, 999)
	assert.ErrorIs(t, err, domains.ErrEventNotFound)
}
