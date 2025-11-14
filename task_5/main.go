package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	// ... do something
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}

/*
Опять про интерфейс. Функция test возвращает значение с конкретным типом,
а значит переменная err с интерфейсом error не может быть равна nil.(интерфейс равен nil, когда и значение и тип равны nil)
*/
