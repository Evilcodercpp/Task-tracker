package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type Task struct {
	ID   string `gorm:"primaryKey" json:"id"`
	Task string `json:"task"`
}

func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	}

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Task{}); err != nil {
		log.Fatalf("could not migrate database: %v", err)
	}
}

func getTask(w http.ResponseWriter, r *http.Request) {
	var tasks []Task
	if err := db.Find(&tasks).Error; err != nil {
		http.Error(w, "could not fetch tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func postTask(w http.ResponseWriter, r *http.Request) {
	var body Task
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	body.ID = uuid.NewString()

	if err := db.Create(&body).Error; err != nil {
		http.Error(w, "could not create task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(body)
}

func patchTask(w http.ResponseWriter, r *http.Request) {
	var body Task
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if body.ID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	var task Task
	if err := db.First(&task, "id = ?", body.ID).Error; err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	if body.Task != "" {
		task.Task = body.Task
	}

	if err := db.Save(&task).Error; err != nil {
		http.Error(w, "could not update task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	var body Task
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if body.ID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	result := db.Delete(&Task{}, "id = ?", body.ID)
	if result.Error != nil {
		http.Error(w, "could not delete task", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	initDB()

	http.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTask(w, r)
		case http.MethodPost:
			postTask(w, r)
		case http.MethodPatch:
			patchTask(w, r)
		case http.MethodDelete:
			deleteTask(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
