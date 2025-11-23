package main

import (
	"fmt"
	"github.com/M-kos/wb_level2/task_13/internal/mycut"
)

func main() {
	err := mycut.Run()
	if err != nil {
		fmt.Println("error: ", err)
	}
}
