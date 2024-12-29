package service

import (
	"fmt"
	"sync"
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

	sem := make(chan struct{}, 10) // Создание семафора для ограничения количества горутин

	inputChan := make(chan string)
	outputChan := make(chan string)

	var wg sync.WaitGroup

	go func() {
		defer close(inputChan)
		lines, err := s.prod.Produce()
		if err != nil {
			fmt.Println("Error reading lines:", err)
			return
		}
		for _, line := range lines {
			inputChan <- line
		}
	}()

	// Запуск основных горутин
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range inputChan {
				sem <- struct{}{}
				go func(line string) {
					defer func() { <-sem }()
					maskedLine := s.replaceLinks(line)
					outputChan <- maskedLine
				}(line)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(outputChan)
	}()

	var results []string
	go func() {
		for maskedLine := range outputChan {
			results = append(results, maskedLine)
		}
	}()

	wg.Wait()

	err := s.pres.Present(results)
	if err != nil {
		return err
	}

	return nil
}
