package service

import (
	"os"
)

type Presenter interface {
	Present([]string) error
}
type FilePresenter struct {
	filePath string
}

func NewFilePresenter(filePath string) *FilePresenter { // конструктор для обработчика
	return &FilePresenter{filePath}
}

func (fp *FilePresenter) Present(data []string) error {
	file, err := os.Create(fp.filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, l := range data {
		_, err := file.WriteString(l + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}
