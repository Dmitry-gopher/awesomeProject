package service

import (
	"errors"
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
	data, err := s.prod.Produce()
	if err != nil {
		return errors.New("Producer error: " + err.Error())
	}

	if data == nil || len(data) == 0 { // Если данных нет, обрабатываем как пустой ввод
		data = []string{""} // Добавляем пустую строку для обработки
	}

	// Маскировка ссылок
	var maskedData []string
	for _, line := range data {
		maskedData = append(maskedData, s.replaceLinks(line))
	}

	// Передача замаскированных данных Presenter'у
	err = s.pres.Present(maskedData)
	if err != nil {
		return errors.New("Presenter error: " + err.Error())
	}

	return nil
}
