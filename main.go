package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var task string

type requestBody struct {
	Task string `json:"task"`
}

func getTask(w http.ResponseWriter, r *http.Request) {
	if task == "" {
		task = "guest"
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hello, %s", task)
}

func postTask(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var body requestBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	task = body.Task

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "task saved")
}

func main() {
	http.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getTask(w, r)
		} else if r.Method == http.MethodPost {
			postTask(w, r)
		}
	})

	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}