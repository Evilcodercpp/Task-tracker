package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type requestBody struct {
	ID   string `json:"id"`
	Task string `json:"task"`
}

// Хранилище задач
var tasks []requestBody

// GET /task
func getTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Возвращаем весь список задач
	json.NewEncoder(w).Encode(tasks)
}

// POST /task
func postTask(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 1. создаём пустую структуру
	var body requestBody

	// 2. декодируем JSON
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// 3. генерируем уникальный ID
	body.ID = uuid.NewString()

	// 4. добавляем в хранилище
	tasks = append(tasks, body)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "task saved")
}

func PatchTask(w http.ResponseWriter, r *http.Request){
	
	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 1. создаём пустую структуру
	var body requestBody

	// 2. декодируем JSON
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// ищем айди который нам нужно обновить 
	for i, arr := range tasks{
		if tasks[i].ID == body.ID 
		// 3. обновляем данные в хранилище
		tasks = append(tasks, body)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "task update")
	}
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var body requestBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// ищем задачу по ID
	for i := range tasks {
		if tasks[i].ID == body.ID {

			// удаляем элемент из slice
			tasks = append(tasks[:i], tasks[i+1:]...)

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "task deleted")
			return
		}
	}

	// если не нашли
	http.Error(w, "task not found", http.StatusNotFound)
}

func main() {
	http.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getTask(w, r)
		} else if r.Method == http.MethodPost {
			postTask(w, r)
		} else if r.Method == http.MethodPatch{
			PatchTask(w, r)
		} else if r.Method == http.MethodDelete {
			DeleteTask(w, r)
		}else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}