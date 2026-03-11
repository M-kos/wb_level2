package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/M-kos/wb_level2/task_18/internal/domains"
	"github.com/M-kos/wb_level2/task_18/internal/dto"
	"github.com/M-kos/wb_level2/task_18/internal/middlewares"
)

type EventService interface {
	EventsForDay(ctx context.Context, userId int, date time.Time) ([]*domains.Event, error)
	EventsForWeek(ctx context.Context, userId int, date time.Time) ([]*domains.Event, error)
	EventsForMonth(ctx context.Context, userId int, date time.Time) ([]*domains.Event, error)
	Create(ctx context.Context, newEvent *domains.Event) (*domains.Event, error)
	Update(ctx context.Context, event *domains.Event) (*domains.Event, error)
	Delete(ctx context.Context, eventId int) error
}

type EventHandler struct {
	service EventService
}

func NewEventHandler(router *http.ServeMux, service EventService, middleware middlewares.Middleware) {
	handler := &EventHandler{
		service: service,
	}

	router.HandleFunc("GET /events_for_day", middleware(handler.EventsForDay))
	router.HandleFunc("GET /events_for_week", middleware(handler.EventsForWeek))
	router.HandleFunc("GET /events_for_month", middleware(handler.EventsForMonth))
	router.HandleFunc("POST /create_event", middleware(handler.Create))
	router.HandleFunc("POST /update_event/{id}", middleware(handler.Update))
	router.HandleFunc("POST /delete_event/{id}", middleware(handler.Delete))
}

func (eh *EventHandler) EventsForDay(w http.ResponseWriter, r *http.Request) {
	eh.getEvents(w, r, eh.service.EventsForDay)
}

func (eh *EventHandler) EventsForWeek(w http.ResponseWriter, r *http.Request) {
	eh.getEvents(w, r, eh.service.EventsForWeek)
}

func (eh *EventHandler) EventsForMonth(w http.ResponseWriter, r *http.Request) {
	eh.getEvents(w, r, eh.service.EventsForMonth)
}

func (eh *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	var createEventDto dto.CreateEvent
	if err := json.NewDecoder(r.Body).Decode(&createEventDto); err != nil {
		slog.Error("[Create] error decoding create event", "error", err)
		writeErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := createEventDto.Validate(); err != nil {
		slog.Error("[Create] error validating create event", "error", err)
		writeErrorJSON(w, "validation error", http.StatusBadRequest)
		return
	}

	event, err := eh.service.Create(r.Context(), createEventDto.ToDomain())
	if err != nil {
		slog.Error("[Create] error creating event", "error", err)
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	eventDto := dto.EventDtoFromDomain(event)

	writeJSON(w, http.StatusOK, dto.EventResponse{
		Result: eventDto,
	})
}

func (eh *EventHandler) Update(w http.ResponseWriter, r *http.Request) {
	queryEventId := r.PathValue("id")
	eventId, err := strconv.Atoi(queryEventId)
	if err != nil {
		slog.Error("[Update] error converting query event id to int", "error", err)
		writeErrorJSON(w, "invalid event id", http.StatusBadRequest)
		return
	}

	var updateEventDto dto.UpdateEvent
	if err := json.NewDecoder(r.Body).Decode(&updateEventDto); err != nil {
		slog.Error("[Update] error decoding update event", "error", err)
		writeErrorJSON(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := updateEventDto.Validate(); err != nil {
		slog.Error("[Update] error validating update event", "error", err)
		writeErrorJSON(w, "validation error", http.StatusBadRequest)
		return
	}

	event, err := eh.service.Update(r.Context(), updateEventDto.ToDomain(eventId))
	if err != nil {
		if errors.Is(err, domains.ErrEventNotFound) {
			slog.Error("[Update] error updating event, event not found", "error", err)
			writeErrorJSON(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		slog.Error("[Update] error updating event", "error", err)
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	eventDto := dto.EventDtoFromDomain(event)

	writeJSON(w, http.StatusOK, dto.EventResponse{
		Result: eventDto,
	})
}

func (eh *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	queryEventId := r.PathValue("id")
	eventId, err := strconv.Atoi(queryEventId)
	if err != nil {
		slog.Error("[Delete] error converting query event id to int", "error", err)
		writeErrorJSON(w, "invalid event id", http.StatusBadRequest)
		return
	}

	err = eh.service.Delete(r.Context(), eventId)
	if err != nil {
		if errors.Is(err, domains.ErrEventNotFound) {
			slog.Error("[Delete] error deleting event, event not found", "error", err)
			writeErrorJSON(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		slog.Error("[Delete] error deleting event", "error", err)
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (eh *EventHandler) getEvents(w http.ResponseWriter, r *http.Request, serviceFn func(ctx context.Context, userId int, date time.Time) ([]*domains.Event, error)) {
	queryUserId := r.URL.Query().Get("user_id")
	queryDate := r.URL.Query().Get("date")

	userId, err := strconv.Atoi(queryUserId)
	if err != nil {
		slog.Error("[Get Events] error converting query user id to int", "error", err)
		writeErrorJSON(w, "invalid user id", http.StatusBadRequest)
		return
	}

	date, err := time.Parse(time.DateOnly, queryDate)
	if err != nil {
		slog.Error("[Get Events] error converting query date to time", "error", err)
		writeErrorJSON(w, "invalid date, expected format: YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	select {
	case <-r.Context().Done():
		slog.Info("[Get Events] context done")
		writeErrorJSON(w, "request cancelled by the client", http.StatusRequestTimeout)
		return
	default:
		events, err := serviceFn(r.Context(), userId, date)
		if err != nil {
			slog.Error("[Get Events] error getting events", "error", err)
			writeErrorJSON(w, "something went wrong", http.StatusInternalServerError)
			return
		}

		var results []*dto.EventDto

		for _, event := range events {
			eventDto := dto.EventDtoFromDomain(event)

			results = append(results, eventDto)
		}

		writeJSON(w, http.StatusOK, dto.EventsResponse{
			Result: results,
		})
	}
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeErrorJSON(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse{Error: message})
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("error encoding response", "error", err)
	}
}
