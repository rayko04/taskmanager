package main

import (
	"fmt"
	"net/http"
	"taskmanager/model"
	"encoding/json"
	"taskmanager/repository"
)
func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, taskmanager.")
}

//closure so handler can access repo
func createTasksHandler(repo *repository.TaskRepository) http.HandlerFunc{

	return func (w http.ResponseWriter, r *http.Request) {
		switch r.Method {

			//GET
		case http.MethodGet:
			
			tasks := repo.GetAll()

			data, err := json.Marshal(tasks)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		

			//POST
		case http.MethodPost:
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
