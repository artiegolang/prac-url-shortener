package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	endpoint := "http://localhost:8085/"
	// Приглашение в консоли
	fmt.Println("Введите длинный URL")
	// Открываем потоковое чтение из консоли
	reader := bufio.NewReader(os.Stdin)
	// Читаем строку из консоли
	long, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	long = strings.TrimSuffix(long, "\n")
	// Добавляем HTTP-клиент
	client := &http.Client{}
	// Пишем запрос
	// Запрос методом POST должен содержать тело с URL в формате text/plain
	request, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(long))
	if err != nil {
		panic(err)
	}
	// В заголовках запроса указываем Content-Type: text/plain
	request.Header.Add("Content-Type", "text/plain")
	// Отправляем запрос и получаем ответ
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	// Выводим код ответа
	fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()
	// Читаем поток из тела ответа
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	// И печатаем его
	fmt.Println(string(body))
}
