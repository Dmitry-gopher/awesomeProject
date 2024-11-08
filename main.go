package main

import "fmt"

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

	text := "Hello, its my page: http://localhost123.com See you"
	text1 := replaceLinks(text)
	fmt.Println(text1)

}
