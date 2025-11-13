package main

import "fmt"

func test() (x int) {
	defer func() {
		x++
	}()
	x = 1
	return
}

func anotherTest() int {
	var x int
	defer func() {
		x++
	}()
	x = 1
	return x
}

func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())
}

/*
defer'ы выполнятся в порядке запуков функций
в первом выведется 2, потому что x - именнованное возвращаемое значение
во втором 1, так как возвращается неименованное значение
*/
