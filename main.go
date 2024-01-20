package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Определяем обработчик для пути "/"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Привет, мир!")
	})

	// Запускаем HTTP сервер на порту 8080
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
