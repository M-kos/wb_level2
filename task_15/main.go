package main

import (
	"fmt"
	"github.com/M-kos/wb_level2/task_15/internal/myshell"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	muShell := myshell.NewMyShell()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		for range sigs {
			fmt.Println("\n^C")
		}
	}()

	muShell.Run()
}
