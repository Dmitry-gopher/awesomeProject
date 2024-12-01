package main

import (
	"awesomeProject/service"
	"fmt"
	"os"
)

func main() {
	// Создаём экземпляры Producer и Presenter
	producer := service.NewFileProducer("input.txt")
	presenter := service.NewFilePresenter("output.txt")
	srv := service.NewService(producer, presenter)

	// Запуск сервиса
	if err := srv.Run(); err != nil {
		fmt.Printf("Ошибка при запуске сервиса: %v\n", err)
		os.Exit(1) // Завершаем программу с кодом ошибки
	}

	fmt.Println("Данные успешно обработаны!")
}
