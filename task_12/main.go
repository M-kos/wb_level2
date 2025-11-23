package main

import (
	"fmt"
	"github.com/M-kos/wb_level2/task_12/internal/mygrep"
)

func main() {
	err := mygrep.Run()
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
