package service

import (
	"fmt"
)

type Producer interface {
	Produce() ([]string, error)
}

type Presenter interface {
	Present([]string) error
}

type Service struct {
	prod Producer
	pres Presenter
}

func NewService(prod Producer, pres Presenter) *Service { // конструктор для сервиса
	return &Service{prod, pres}
}

// Функция маскирования ссылок
func (s *Service) replaceLinks(data string) string {
	if len(data) == 0 { // Если строка пустая, просто возвращаем её
		return data
	}
	buffer := []byte(data)    // Преобразуем строку в срез байтов
	link := []byte("http://") // Что будем искать и маскировать

	for i := 0; i < len(buffer); i++ {
		if i+len(link) <= len(buffer) && string(buffer[i:i+len(link)]) == "http://" {
			k := i + len(link)
			for k < len(buffer) && buffer[k] != ' ' {
				buffer[k] = '*'
				k++
			}
		}
	}
	return string(buffer)
}

// Запуск сервиса
func (s *Service) Run() error {

	lines, err := s.prod.Produce()
	if err != nil {
		fmt.Println("Error reading lines:", err)
	}

	sem := make(chan struct{}, 10) // Создание семафора для ограничения количества горутин
	done := make(chan struct{})    // Синхронизация завершения всех горутин
	inputChan := make(chan string)
	outputChan := make(chan string)

	// Запуск горутины для сбора результатов
	go func() {
		results := make([]string, len(lines))
		for i := 0; i < len(lines); i++ {
			results[i] = <-outputChan
		}
		_ = s.pres.Present(results)
		close(done)
	}()

	// Отправление данных в канал для обработки
	go func() {
		defer close(inputChan)
		for _, line := range lines {
			inputChan <- line
		}
	}()

	// Обработка строк
	for line := range inputChan {
		sem <- struct{}{}
		go func(l string) {
			defer func() { <-sem }()
			maskedLine := s.replaceLinks(l)
			outputChan <- maskedLine
		}(line)
	}

	// Ожидание завершения всех горутин
	<-done

	return nil
}
