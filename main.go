package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var (
	i     int = 0
	mutex sync.Mutex
)

func main() {
	go handlePostRequests()

	// Запускаем второй поток, который каждую секунду выводит содержимое файла
	go func() {
		for {
			readAndPrintFileContents()
			time.Sleep(time.Second)
		}
	}()

	// Запускаем HTTP сервер
	http.ListenAndServe(":8080", nil)
}

func handlePostRequests() {
	http.HandleFunc("/write", func(w http.ResponseWriter, r *http.Request) {
		// Читаем тело запроса
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		// Записываем результат в файл
		writeToFile(string(body))
		fmt.Fprint(w, "Data written to file")
	})
}

func writeToFile(data string) {
	// Захватываем мьютекс для предотвращения одновременной записи из нескольких потоков
	mutex.Lock()
	defer mutex.Unlock()

	// Открываем файл для добавления данных
	fmt.Println(data)
}

func readAndPrintFileContents() {
	// Захватываем мьютекс для предотвращения одновременного чтения из нескольких потоков
	mutex.Lock()
	defer mutex.Unlock()

	// Читаем содержимое файла
	fmt.Println("Thread 2 iteration ", i)
	i++
}
