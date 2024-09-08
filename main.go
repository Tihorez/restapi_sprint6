package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func getAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tasksList := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		tasksList = append(tasksList, task)
	}

	if err := json.NewEncoder(w).Encode(tasksList); err != nil {
		http.Error(w, "Ошибка при сериализации задач", http.StatusInternalServerError)
		return
	}
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if newTask.ID == "" || newTask.Description == "" {
		http.Error(w, "Отсутствуют обязательные поля", http.StatusBadRequest)
		return
	}

	tasks[newTask.ID] = newTask
	w.WriteHeader(http.StatusCreated)
}

func getTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	task, exists := tasks[id]
	if !exists {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Ошибка при сериализации задачи", http.StatusInternalServerError)
		return
	}
}

func deleteTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	_, exists := tasks[id]
	if !exists {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	delete(tasks, id)
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getAllTasksHandler)
	r.Post("/tasks", createTaskHandler)
	r.Get("/tasks/{id}", getTaskByIDHandler)
	r.Delete("/tasks/{id}", deleteTaskByIDHandler)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}