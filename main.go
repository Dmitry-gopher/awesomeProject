package main

import (
	"awesomeProject/service"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:] // получение аргументов из терминала
	inputFile := args[0]
	outputFile := "output.txt" // значение по умолчанию

	if len(args) > 1 {
		outputFile = args[1]
	}

	producer := service.NewFileProducer(inputFile)
	presenter := service.NewFilePresenter(outputFile)
	service := service.NewService(producer, presenter)

	if err := service.Run(); err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Println("Данные успешно обработаны")
}
