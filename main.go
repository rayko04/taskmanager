package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"taskmanager/database"
	"taskmanager/model"
	"taskmanager/repository"
	"time"

	"github.com/joho/godotenv"
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

func jsonResponse(w http.ResponseWriter, status int, data any) {

	bytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(bytes)
}

//a valid post request must conatin a valid title 
func validReq(req *model.TaskRequest) bool {
	
	if strings.TrimSpace(req.Title) == "" {
		return false
	}
	return true
}

//a valid update request must contain all fields
func validUpdateReq(req *model.TaskUpdateRequest) bool {

	if req.Title == nil || req.Description == nil || req.Completed == nil {
		return false
	}
	return true
}

//a valid patch request must contain atleast one field
func validPatchReq(req *model.TaskPatchRequest) bool {

	if req.Title == nil && req.Description == nil && req.Completed == nil {
		return false
	}
	return true
}

//decode JSON into Golang
func decodeJSON(read io.ReadCloser, dest any) error {
	
	decoder := json.NewDecoder(read)
	decoder.DisallowUnknownFields()	//err if any unknown value inside json
	
	err := decoder.Decode(dest)
	if err != nil {
		return err
	}

	err = decoder.Decode(&struct{}{})	// err if any non json trailing value
	if err == io.EOF {
		return nil
	}
	return err
}

//handles GET all
func getTasksHandler(repo *repository.TaskRepository) http.HandlerFunc{

	return func (w http.ResponseWriter, r *http.Request) {

		tasks, err := repo.GetAll()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		jsonResponse(w, http.StatusOK, tasks)
	}
}

//handles GET by id
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

		jsonResponse(w, http.StatusOK, task)
	}
}

//handles POST
func createTaskHandler(repo *repository.TaskRepository) http.HandlerFunc{

	return func (w http.ResponseWriter, r *http.Request) {

		//json to golang

		req := &model.TaskRequest{}		//Decode requires a pointer
		err := decodeJSON(r.Body, req)

		if err != nil || !validReq(req) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		task := model.Task {
			Title: req.Title,
			Description: req.Description,
		}
		task, err = repo.Create(task)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		//send back response in json
		jsonResponse(w, http.StatusCreated, task)
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

		req := &model.TaskUpdateRequest{}	//Decode requires a pointer
		err := decodeJSON(r.Body, req)

		if err != nil || !validUpdateReq(req) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		task := model.Task {
			Title: *req.Title,
			Description: *req.Description,
			Completed: *req.Completed,
		}
		updated := false

		task, updated = repo.Update(id, task)
		if !updated {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		//send back response in json
		jsonResponse(w, http.StatusOK, task)
	}
}

//handles PATCH
func patchTaskHandler(repo *repository.TaskRepository) http.HandlerFunc{

	return func (w http.ResponseWriter, r *http.Request) {

		path := trimAndSplit(r.URL.Path)

		id, error := strconv.Atoi(path[len(path)-1])
		if error != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		req := &model.TaskPatchRequest{}	//Decode requires a pointer
		err := decodeJSON(r.Body, req)
		if err != nil || !validPatchReq(req) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		task, patched := repo.Patch(id, *req)
		if !patched {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		//send back response in json
		jsonResponse(w, http.StatusOK, task)
	}
}

//handles delete
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
		case http.MethodDelete, http.MethodPut, http.MethodPatch:
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func individualDispatcher(repo*repository.TaskRepository) http.HandlerFunc {

	return func (w http.ResponseWriter, r *http.Request) {
		
		switch r.Method {
		case http.MethodGet:
			getTaskHandler(repo)(w, r)
		case http.MethodDelete:
			deleteTaskHandler(repo)(w, r)
		case http.MethodPut:
			updateTaskHandler(repo)(w, r)
		case http.MethodPatch:
			patchTaskHandler(repo)(w, r)
		case http.MethodPost:
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		
	}
}

func dispatcher(repo*repository.TaskRepository) http.HandlerFunc {

	return func (w http.ResponseWriter, r *http.Request) {
		path := trimAndSplit(r.URL.Path)
		if endsInTasks(path) {
			collectionDispatcher(repo)(w, r)
		} else if endsInTaskId(path) {
			individualDispatcher(repo)(w, r)
		} else {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	}
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Failed to load .env:", err)
		return
	}

	pool, err := database.NewPool()
	if err != nil {
		fmt.Println("Failed to create pool:", err)
		return
	}

	repo := repository.NewTaskRepository(pool)

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/tasks", dispatcher(repo))
	http.HandleFunc("/tasks/", dispatcher(repo))
 
	//test
	// for i := 0; i < 5; i++ {
	// 	task := model.Task {
	// 		Title: "abc",
	// 		Description: "def",
	// 	}
	// 	repo.Create(task)
	// }

	port := ":8080"

	serv := http.Server{
		Addr: port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	err = serv.ListenAndServe()
	if err != nil {
		fmt.Println("Failed to establish connection:", err)
		return
	}
}
