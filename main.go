package main

import (
	"bufio"
	"fmt"
	"os"
)

func replaceLinks(mes string) string {
	buffer := []byte(mes)  // буфер - байтовый срез
	s := []byte("http://") // откуда маскировать

	found := false // поиск начала ссылки
	i := 0         // индекс, с которого начнётся замена
	for i < len(buffer) {
		for j := 0; j < len(s); j++ {
			if buffer[j+i] == s[j] {
				found = true
			} else {
				found = false
				break
			}
		}
		if found { // нашли http://
			k := i + len(s) // текст ссылку
			for k < len(buffer) && buffer[k] != ' ' {
				buffer[k] = '*'
				k++
			}
		}
		i++
	}
	return string(buffer)
}

func main() {

	fmt.Println("Введите текст сообщения:") // Запрос ввода
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n') // считывание строки

	// Обработка и вывод результата
	textWithoutLinks := replaceLinks(text)
	fmt.Println("Результат:")
	fmt.Println(textWithoutLinks)

}
