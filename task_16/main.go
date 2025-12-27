package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/M-kos/wb_level2/task_16/internal/collector"
	"github.com/M-kos/wb_level2/task_16/pkg/semaphore"
)

const (
	resultDir      = "result"
	requestTimeout = 10 * time.Second
	workerCount    = 5
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("need <url>")
		return
	}

	urls := make([]*url.URL, 0)

	for _, arg := range os.Args[1:] {
		u, err := url.ParseRequestURI(arg)
		if err != nil {
			u, err = url.Parse("https://" + arg)
			if err != nil {
				fmt.Printf("invalid url %s: %v\n", arg, err)
				continue
			}
		}

		urls = append(urls, u)
	}

	if len(urls) == 0 {
		fmt.Println("no valid urls provided")
		return
	}

	sem := semaphore.NewSemaphore(workerCount)

	ctr := collector.NewCollector(urls, sem, resultDir, requestTimeout)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	err := ctr.Start(ctx)
	if err != nil {
		fmt.Println("error:", err)
	}
}
