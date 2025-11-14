package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v, ok := <-a:
				if ok {
					c <- v
				} else {
					a = nil
				}
			case v, ok := <-b:
				if ok {
					c <- v
				} else {
					b = nil
				}
			}
			if a == nil && b == nil {
				close(c)
				return
			}
		}
	}()
	return c
}

func main() {
	rand.Seed(time.Now().Unix())
	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4, 6, 8)
	c := merge(a, b)
	for v := range c {
		fmt.Print(v)
	}
}

/*
Не очень понятно что объяснить. Создаются два канала, куда пишутся числа с рандомной задержкой. Затем эти два канаоа
мерджатся в один. В select пытаемся прочитать из этих двух каналов и проверяем, если он закрыт то присваиваем переменным nil.
Это для того, чтобы select больше не пфтался из него прочитать.
И в main, из этого смердженного канала, читаем и выводим в консоль, все что туда приходит.
В итоге, вывод будет от 1 до 8 в случайном порядке.
*/
