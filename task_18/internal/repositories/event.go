package repositories

import (
	"context"
	"sync"
	"time"

	"github.com/M-kos/wb_level2/task_18/internal/domains"
)

type EventRepository struct {
	currentId int
	mu        sync.RWMutex
	store     map[int]*domains.Event
}

func NewEventRepository() *EventRepository {
	return &EventRepository{
		currentId: 1,
		mu:        sync.RWMutex{},
		store:     make(map[int]*domains.Event),
	}
}

func (er *EventRepository) Event(ctx context.Context, id int) (*domains.Event, error) {
	er.mu.RLock()
	defer er.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		event, ok := er.store[id]
		if !ok {
			return nil, domains.ErrEventNotFound
		}

		return event, nil
	}

}

func (er *EventRepository) List(ctx context.Context, userId int, from, to time.Time) ([]*domains.Event, error) {
	er.mu.RLock()
	defer er.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		events := make([]*domains.Event, 0)

		for _, event := range er.store {
			if event.UserID == userId && !event.Date.Before(from) && event.Date.Before(to) {
				events = append(events, event)
			}
		}

		return events, nil
	}
}

func (er *EventRepository) Create(ctx context.Context, newEvent *domains.Event) (*domains.Event, error) {
	er.mu.Lock()
	defer er.mu.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		newEvent.ID = er.currentId
		er.store[er.currentId] = newEvent
		er.currentId++

		return newEvent, nil
	}
}

func (er *EventRepository) Update(ctx context.Context, event *domains.Event) (*domains.Event, error) {
	er.mu.Lock()
	defer er.mu.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		_, ok := er.store[event.ID]
		if !ok {
			return nil, domains.ErrEventNotFound
		}

		er.store[event.ID] = event

		return event, nil
	}
}

func (er *EventRepository) Delete(ctx context.Context, id int) error {
	er.mu.Lock()
	defer er.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, ok := er.store[id]
		if !ok {
			return domains.ErrEventNotFound
		}

		delete(er.store, id)

		return nil

	}
}
