package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/M-kos/wb_level2/task_18/internal/config"
	"github.com/M-kos/wb_level2/task_18/internal/handlers"
	"github.com/M-kos/wb_level2/task_18/internal/middlewares"
	"github.com/M-kos/wb_level2/task_18/internal/repositories"
	"github.com/M-kos/wb_level2/task_18/internal/services"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		slog.Error("error loading config", "error", err)
		return
	}

	router := http.NewServeMux()

	eventRepository := repositories.NewEventRepository()
	eventService := services.NewEventService(eventRepository)

	handlers.NewEventHandler(router, eventService, middlewares.LoggingMiddleware)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Port),
		Handler: router,
	}

	if err = server.ListenAndServe(); err != nil {
		slog.Error("error starting server", "error", err)
		return
	}
}
