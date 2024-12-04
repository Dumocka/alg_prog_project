package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Task представляет собой задание для тестирования
type Task struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Difficulty  string   `json:"difficulty"`
	Tags        []string `json:"tags"`
}

// Простое хранилище данных в памяти
var tasks = make(map[string]Task)

// enableCORS добавляет CORS заголовки к ответу
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next(w, r)
    }
}

func main() {
	// Добавим тестовые данные
	tasks["1"] = Task{
		ID:          "1",
		Title:       "Первое задание",
		Description: "Описание первого задания",
		Difficulty:  "easy",
		Tags:        []string{"math", "basic"},
	}

	// Настраиваем маршруты с CORS
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Write([]byte("Testing Platform API"))
			return
		}
		http.NotFound(w, r)
	})
	http.HandleFunc("/tasks", enableCORS(handleTasks))
	http.HandleFunc("/task/", enableCORS(handleTask))

	// Запускаем сервер
	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// handleTasks обрабатывает GET и POST запросы для списка заданий
func handleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		// Получаем все задания
		taskList := make([]Task, 0, len(tasks))
		for _, task := range tasks {
			taskList = append(taskList, task)
		}
		json.NewEncoder(w).Encode(taskList)

	case "POST":
		// Создаем новое задание
		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Генерируем простой ID
		task.ID = fmt.Sprintf("%d", len(tasks)+1)
		tasks[task.ID] = task

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(task)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleTask обрабатывает GET, PUT и DELETE запросы для конкретного задания
func handleTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Получаем ID задания из URL
	id := r.URL.Path[len("/task/"):]
	if id == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		// Получаем задание по ID
		task, exists := tasks[id]
		if !exists {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(task)

	case "PUT":
		// Обновляем задание
		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		task.ID = id
		tasks[id] = task
		json.NewEncoder(w).Encode(task)

	case "DELETE":
		// Удаляем задание
		if _, exists := tasks[id]; !exists {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		delete(tasks, id)
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
