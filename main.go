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

var tasks []requestBody

func getTask(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(tasks)
}

func postTask(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    var body requestBody
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, "invalid json", http.StatusBadRequest)
        return
    }

    body.ID = uuid.NewString()
    tasks = append(tasks, body)

    w.WriteHeader(http.StatusCreated)
    fmt.Fprintf(w, "task created\n")
}

func patchTask(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPatch {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    var body requestBody
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, "invalid json", http.StatusBadRequest)
        return
    }

    if body.ID == "" {
        http.Error(w, "id is required", http.StatusBadRequest)
        return
    }

    for i := range tasks {
        if tasks[i].ID == body.ID {
            if body.Task != "" {
                tasks[i].Task = body.Task
            }
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, "task updated\n")
            return
        }
    }

    http.Error(w, "task not found", http.StatusNotFound)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    var body requestBody
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, "invalid json", http.StatusBadRequest)
        return
    }

    if body.ID == "" {
        http.Error(w, "id is required", http.StatusBadRequest)
        return
    }

    for i := range tasks {
        if tasks[i].ID == body.ID {
            tasks = append(tasks[:i], tasks[i+1:]...)
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, "task deleted\n")
            return
        }
    }

    http.Error(w, "task not found", http.StatusNotFound)
}

func main() {
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
    http.ListenAndServe(":8080", nil)
}
