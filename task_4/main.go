package main

func main() {
	ch := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()

	for n := range ch {
		println(n)
	}
}

/*
Выведется от 0 до 9 включительно, а потом случится deadlock, так как пишущая горутина завершилась,
но канал не закрыт, поэтому с main будет бесконечно ждать новых данных
*/
