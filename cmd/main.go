package main

import (
	"fmt"
	"log"
	"net/http"

	"Task-tracker/interal/db"
	"Task-tracker/interal/db/handlers"
	taskservice "Task-tracker/interal/db/taskService"
)

func main() {
	db.InitDB()

	repo := taskservice.NewTaskRepo(db.DB)
	h := handlers.NewHandler(repo)

	http.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetTasks(w, r)
		case http.MethodPost:
			h.PostTask(w, r)
		case http.MethodPatch:
			h.PatchTask(w, r)
		case http.MethodDelete:
			h.DeleteTask(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
