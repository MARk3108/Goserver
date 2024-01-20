package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	filePath = "result.txt"
	mutex    sync.Mutex
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
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Записываем данные в файл
	if _, err := file.WriteString(data + "\n"); err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func readAndPrintFileContents() {
	// Захватываем мьютекс для предотвращения одновременного чтения из нескольких потоков
	mutex.Lock()
	defer mutex.Unlock()

	// Читаем содержимое файла
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	lines := strings.Split(string(content), "\n")
	if len(lines) == 2 {
		a, err := strconv.Atoi(lines[0][0 : len(lines[0])-1])
		b, err := strconv.Atoi(lines[1])
		if err != nil {
			fmt.Println("Convert error ", err)
		}
		emptyContent := []byte("")
		err = ioutil.WriteFile(filePath, emptyContent, 0644)
		if err != nil {
			fmt.Println("Clearing file error ", err)
		}

		fmt.Println("A = ", a)
		fmt.Println("B = ", b)
		fmt.Println("Summ is ", a+b)
	} else {
		// Выводим содержимое файла
		fmt.Println("File Contents:")
		fmt.Println(string(content))
	}
}
