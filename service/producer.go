package service

import (
	"bufio"
	"os"
)

type FileProducer struct {
	filePath string
}

func NewFileProducer(filePath string) *FileProducer { // конструктор для поставщика
	return &FileProducer{filePath}
}

func (p *FileProducer) Produce() ([]string, error) { // метод чтения строк из файла
	file, err := os.Open(p.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}
