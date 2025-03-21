package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
)

var db *sql.DB

func init() {
	var err error
	connStr := "user=youruser password=yourpassword dbname=taskdb sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

type Task struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	json.NewDecoder(r.Body).Decode(&task)

	if task.UserID == 0 {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO tasks (user_id, title, status) VALUES ($1, $2, $3) RETURNING id"
	err := db.QueryRow(query, task.UserID, task.Title, task.Status).Scan(&task.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	userID := r.URL.Query().Get("user_id")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if limit == 0 {
		limit = 20 // Default pagination for web
	}
	offset := (page - 1) * limit

	var rows *sql.Rows
	var err error

	if status != "" && userID != "" {
		rows, err = db.Query("SELECT id, user_id, title, status FROM tasks WHERE status=$1 AND user_id=$2 LIMIT $3 OFFSET $4", status, userID, limit, offset)
	} else if userID != "" {
		rows, err = db.Query("SELECT id, user_id, title, status FROM tasks WHERE user_id=$1 LIMIT $2 OFFSET $3", userID, limit, offset)
	} else {
		rows, err = db.Query("SELECT id, user_id, title, status FROM tasks LIMIT $1 OFFSET $2", limit, offset)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Status)
		tasks = append(tasks, task)
	}

	json.NewEncoder(w).Encode(tasks)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var task Task
	json.NewDecoder(r.Body).Decode(&task)

	_, err := db.Exec("UPDATE tasks SET user_id=$1, title=$2, status=$3 WHERE id=$4", task.UserID, task.Title, task.Status, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	_, err := db.Exec("DELETE FROM tasks WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/tasks", createTask).Methods("POST")
	r.HandleFunc("/tasks", getTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")

	log.Println("Task Management Service running on port 8080")
	http.ListenAndServe(":8080", r)
}
