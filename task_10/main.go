package main

import (
	"fmt"
	"github.com/M-kos/wb_level2/task_10/internal/mysort"
)

func main() {
	res, err := mysort.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, s := range res {
		fmt.Println(s)
	}
}
