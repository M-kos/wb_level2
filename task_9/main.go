package main

import (
	"fmt"
	"github.com/M-kos/wb_level2/task_9/internal/unpacker"
)

func main() {
	result, err := unpacker.Unpack("a4bc2d5e")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}
