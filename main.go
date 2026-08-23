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
		case http.MethodGet:
			
			tasks := repo.GetAll()

			data, err := json.Marshal(tasks)
			if err != nil {
				fmt.Println("error marshalling")
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		
		case http.MethodPost:
			fmt.Fprintln(w, "post")
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
