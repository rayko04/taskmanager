package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"taskmanager/model"
	"taskmanager/repository"
)
func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, taskmanager.")
}

func trimAndSplit(str string) []string {
		pathstr := strings.TrimRight(str, "/")
		path := strings.Split(pathstr, "/")

		return path
}

func endsInTasks(path []string) bool {
	return len(path) == 2 && path[len(path)-1] == "tasks"
}

func endsInTaskId(path []string) bool {
	return len(path) == 3 && path[len(path)-2] == "tasks"
}

//handles getall()
func getTasksHandler(repo *repository.TaskRepository) http.HandlerFunc{

	return func (w http.ResponseWriter, r *http.Request) {

		tasks := repo.GetAll()
		
		data, err := json.Marshal(tasks)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func getTaskHandler(repo *repository.TaskRepository) http.HandlerFunc{

	return func (w http.ResponseWriter, r *http.Request) {

		path := trimAndSplit(r.URL.Path)

		id, error := strconv.Atoi(path[len(path)-1])
		if error != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		
		task, exists := repo.GetById(id)
		if !exists {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		data, err := json.Marshal(task)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

//handles POST
func createTaskHandler(repo *repository.TaskRepository) http.HandlerFunc{

	return func (w http.ResponseWriter, r *http.Request) {

		//json to go
		decoder := json.NewDecoder(r.Body)
		req := &model.TaskRequest{}		//Decode requires a pointer
		err := decoder.Decode(req)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		task := model.Task {
			Title: req.Title,
			Description: req.Description,
		}
		task = repo.Create(task)

		//send back response in json
		data, err := json.Marshal(task)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	}
}

//handles PUT
func updateTaskHandler(repo *repository.TaskRepository) http.HandlerFunc{

	return func (w http.ResponseWriter, r *http.Request) {

		path := trimAndSplit(r.URL.Path)

		
		id, error := strconv.Atoi(path[len(path)-1])
		if error != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(r.Body)
		req := &model.TaskUpdateRequest{}	//Decode requires a pointer
		err := decoder.Decode(req)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		task := model.Task {
			Title: req.Title,
			Description: req.Description,
			Completed: req.Completed,
		}
		updated := false

		task, updated = repo.Update(id, task)
		if !updated {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		//send back response in json
		data, err := json.Marshal(task)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func deleteTaskHandler(repo *repository.TaskRepository) http.HandlerFunc{

	return func (w http.ResponseWriter, r *http.Request) {

		path := trimAndSplit(r.URL.Path)

		id, error := strconv.Atoi(path[len(path)-1])
		if error != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		deleted := repo.Delete(id)
		if !deleted {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func collectionDispatcher(repo *repository.TaskRepository) http.HandlerFunc {
	
	return func (w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTasksHandler(repo)(w, r)
		case http.MethodPost:
			createTaskHandler(repo)(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func individualDispatcher(repo*repository.TaskRepository) http.HandlerFunc {

	return func (w http.ResponseWriter, r *http.Request) {

		path := trimAndSplit(r.URL.Path)
		if endsInTasks(path) {
			http.Redirect(w, r, "/tasks", http.StatusMovedPermanently)
			return
		} else if endsInTaskId(path) {
			switch r.Method {
			case http.MethodGet:
				getTaskHandler(repo)(w, r)
			case http.MethodDelete:
				deleteTaskHandler(repo)(w, r)
			case http.MethodPut:
				updateTaskHandler(repo)(w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}
		} else {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
	}
}

func main() {

	repo := repository.NewTaskRepository()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/tasks", collectionDispatcher(repo))
	http.HandleFunc("/tasks/", individualDispatcher(repo))

	port := ":8080"
 
	//test: no POST yet
	for i := 0; i < 5; i++ {
		task := model.Task {
			Title: "abc",
			Description: "def",
		}
		repo.Create(task)
	}

	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Failed to establish connection:", err)
		return
	}

}
