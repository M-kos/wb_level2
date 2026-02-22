package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/M-kos/wb_level2/task_17/internal/telnet"
)

const DefaultTimeout = 10 * time.Second

func main() {
	timeout := flag.Duration("timeout", DefaultTimeout, "connection timeout")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Printf("Usage: %s [options] <host> <port>\n", os.Args[0])
		return
	}

	config := telnet.NewConfig(args[0], args[1], *timeout)
	client := telnet.NewClient(config)

	if err := client.Connect(); err != nil {
		fmt.Printf("error connecting to telnet: %v\n", err)
		return
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	if err := client.Run(ctx, os.Stdin, os.Stdout); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
