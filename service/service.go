package service

import "errors"

type Service struct {
	prod Producer
	pres Presenter
}

func NewService(prod Producer, pres Presenter) *Service { // конструктор для сервиса
	return &Service{prod, pres}
}

// функция маскирования ссылок
func (s *Service) replaceLinks(data []string) []string {
	for i, line := range data {
		buffer := []byte(line)    // буфер - байтовый срез
		link := []byte("http://") // откуда маскировать
		for i := 0; i < len(buffer); i++ {
			if i+len(link) <= len(buffer) && string(buffer[i:i+len(link)]) == "http://" { // поиск "http://"
				k := i + len(link)
				for k < len(buffer) && buffer[k] != ' ' {
					buffer[k] = '*'
					k++
				}
			}
		}
		data[i] = string(buffer)
	}
	return data
}

// запуск сервиса
func (s *Service) Run() error {
	data, err := s.prod.Produce()
	if err != nil {
		return errors.New("Producer error: " + err.Error())
	}
	if data == nil {
		data = []string{}
	}
	maskedData := s.replaceLinks(data)
	err = s.pres.Present(maskedData)
	if err != nil {
		return errors.New("Presenter error: " + err.Error())
	}
	return nil
}
