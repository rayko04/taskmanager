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

//closure so handler can access repo
func createTasksHandler(repo *repository.TaskRepository) http.HandlerFunc{

	return func (w http.ResponseWriter, r *http.Request) {
		
		//trim ending "/"
		pathstr := strings.TrimRight(r.URL.Path, "/")
		path := strings.Split(pathstr, "/")

		switch r.Method {

			//GET
		case http.MethodGet:
			
			//if ends in tasks
			if len(path) == 2 && path[len(path)-1] == "tasks" {
				tasks := repo.GetAll()
				
				data, err := json.Marshal(tasks)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.Write(data)
				
			} else if len(path) == 3 && path[len(path)-2] == "tasks" {
				//if ends in tasks/[smthin]

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
			} else {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}

			//POST
		case http.MethodPost:

			if len(path) != 2 || path[len(path)-1] != "tasks" {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
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

		case http.MethodDelete: 

			//if ends in tasks
			if len(path) == 2 && path[len(path)-1] == "tasks" {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			} else if len(path) == 3 && path[len(path)-2] == "tasks" {
				//if ends in tasks/[smthin]

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
			} else {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}	
}

func main() {

	repo := repository.NewTaskRepository()
	tasksHandler := createTasksHandler(repo)	//dependency injection through closure

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/tasks", tasksHandler)
	http.HandleFunc("/tasks/", tasksHandler)

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
