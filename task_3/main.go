package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)
	fmt.Println(err == nil)
}

/*
Интерфейс - под капотом, это структура с двумя полями: указатель на тип и указатель на данные.
Интерфейс равен nil только тогда, когда оба этих указателя равны nil.
В свою ч
очередь, в пустом интерфейсе (interface{}) нет указателя на тип, поэтому он равен nil, когда указатель на данные равен nil.
*/
