package main

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"taskmanager/model"

	"testing"
	"time"

	"github.com/jackc/pgx/v5"
)

// VALIDATE REQ
func TestValidReq(t *testing.T) {

	tests := []struct {
		name string
		req  model.TaskRequest
		want bool
	}{
		{
			name: "valid title",
			req: model.TaskRequest{
				Title: "Get Job",
			},
			want: true,
		},
		{
			name: "empty title",
			req: model.TaskRequest{
				Title: "",
			},
			want: false,
		},
		{
			name: "whitespace title",
			req: model.TaskRequest{
				Title: "    ",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validReq(&tt.req)

			if got != tt.want {
				t.Errorf("validReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

// UPDATE VALIDATION
func TestValidUpdateReq(t *testing.T) {

	tests := []struct {
		name string
		req  model.TaskUpdateRequest
		want bool
	}{
		{
			name: "No Title",
			req: model.TaskUpdateRequest{
				Title:       nil,
				Description: new("test"),
				Completed:   new(false),
			},
			want: false,
		},
		{
			name: "No Description",
			req: model.TaskUpdateRequest{
				Title:       new("test"),
				Description: nil,
				Completed:   new(false),
			},
			want: false,
		},
		{
			name: "No Completed",
			req: model.TaskUpdateRequest{
				Title:       new("test"),
				Description: new("test"),
				Completed:   nil,
			},
			want: false,
		},
		{
			name: "None",
			req: model.TaskUpdateRequest{
				Title:       nil,
				Description: nil,
				Completed:   nil,
			},
			want: false,
		},
		{
			name: "Valid Req",
			req: model.TaskUpdateRequest{
				Title:       new("test"),
				Description: new("test"),
				Completed:   new(false),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validUpdateReq(&tt.req)

			if got != tt.want {
				t.Errorf("validUpdateReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

// PATCH VALIDATION
func TestValidPatchReq(t *testing.T) {

	tests := []struct {
		name string
		req  model.TaskPatchRequest
		want bool
	}{
		{
			name: "No Title",
			req: model.TaskPatchRequest{
				Title:       nil,
				Description: new("test"),
				Completed:   new(false),
			},
			want: true,
		},
		{
			name: "No Description",
			req: model.TaskPatchRequest{
				Title:       new("test"),
				Description: nil,
				Completed:   new(false),
			},
			want: true,
		},
		{
			name: "No Completed",
			req: model.TaskPatchRequest{
				Title:       new("test"),
				Description: new("test"),
				Completed:   nil,
			},
			want: true,
		},
		{
			name: "Invalid Req",
			req: model.TaskPatchRequest{
				Title:       nil,
				Description: nil,
				Completed:   nil,
			},
			want: false,
		},
		{
			name: "All Fields",
			req: model.TaskPatchRequest{
				Title:       new("test"),
				Description: new("test"),
				Completed:   new(false),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validPatchReq(&tt.req)

			if got != tt.want {
				t.Errorf("validPatchReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

// COLLECTION PATH
func TestIsCollectionPath(t *testing.T) {

	tests := []struct {
		name string
		path []string
		want bool
	}{
		{
			name: "valid path",
			path: []string{"", "tasks"},
			want: true,
		},
		{
			name: "invalid path",
			path: []string{"api", "tasks"},
			want: false,
		},
		{
			name: "invalid path",
			path: []string{"api", "", "tasks"},
			want: false,
		},
		{
			name: "invalid path",
			path: []string{"", "task"},
			want: false,
		},
		{
			name: "invalid path",
			path: []string{"", "tasks", "1"},
			want: false,
		},
		{
			name: "invalid path",
			path: []string{},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isCollectionPath(tt.path)

			if got != tt.want {
				t.Errorf("isCollectionPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

// INDIVIDUAL PATH
func TestIsIndividualPath(t *testing.T) {

	tests := []struct {
		name string
		path []string
		want bool
	}{
		{
			name: "valid path",
			path: []string{"", "tasks", "5"},
			want: true,
		},
		{
			name: "invalid path",
			path: []string{"api", "tasks"},
			want: false,
		},
		{
			name: "invalid path",
			path: []string{"api", "", "tasks"},
			want: false,
		},
		{
			name: "invalid path",
			path: []string{"", "tasks", "5", "8"},
			want: false,
		},
		{
			name: "invalid path",
			path: []string{},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isIndividualPath(tt.path)

			if got != tt.want {
				t.Errorf("isIndividualPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

// JSON DECODE
func TestDecodeJSON(t *testing.T) {

	tests := []struct {
		name string
		read io.ReadCloser
		want bool
	}{
		{
			name: "valid json",
			read: io.NopCloser(strings.NewReader(`{"title":"Buy milk","description":"From the store"}`)),
			want: false,
		},
		{
			name: "empty",
			read: io.NopCloser(strings.NewReader(`{}`)),
			want: false,
		},
		{
			name: "invalid json",
			read: io.NopCloser(strings.NewReader(`{"title":"Buy milk"`)),
			want: true,
		},
		{
			name: "unknown field",
			read: io.NopCloser(strings.NewReader(`{"title":"Buy milk","foo":"From the store"}`)),
			want: true,
		},
		{
			name: "multiple json",
			read: io.NopCloser(strings.NewReader(`{"title":"Buy milk"}{"description":"From the store"}`)),
			want: true,
		},
		{
			name: "trailing garbage",
			read: io.NopCloser(strings.NewReader(`{"title":"Buy milk","description":"From the store"} garbage`)),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := false
			err := decodeJSON(tt.read, &model.TaskRequest{})

			if err != nil {
				got = true
			}

			if got != tt.want {
				t.Errorf("decodeJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ROOT RESPONSE
func TestRootHandler(t *testing.T) {

	request := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	rootHandler(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("Expected: %v, Got: %v", http.StatusOK, recorder.Code)
	}
	if recorder.Body.String() != "Hello, taskmanager.\n" { //used println in handler so '\n'
		t.Errorf("Mismatched Body Error")
	}
}

// CREATE FAIL
func TestCreateTaskHandlerFail(t *testing.T) {

	tests := []struct {
		name     string
		body     io.Reader
		wantCode int
	}{
		{
			name:     "Invalid Json",
			body:     strings.NewReader(`{json`),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Empty Title",
			body:     strings.NewReader(`{"title":""}`),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Whitespace Title",
			body:     strings.NewReader(`{"title":" "}`),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Unknown Field",
			body:     strings.NewReader(`{"foo":"bar"}`),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Multiple Json",
			body:     strings.NewReader(`{"title":"x"}{"title":"y"}`),
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request := httptest.NewRequest("POST", "/tasks", tt.body)
			recorder := httptest.NewRecorder()
			createTaskHandler(nil)(recorder, request)

			if recorder.Code != tt.wantCode {
				t.Errorf("Expected: %v, Got: %v", tt.wantCode, recorder.Code)
			}
		})
	}
}

type fakeRepo struct {
	tasks   []model.Task
	task    model.Task
	err     error
	deleted bool
}

func (repo *fakeRepo) Create(task model.Task) (model.Task, error) {
	return repo.task, nil
}

// CREATE SUCCESS
func TestCreateTaskHandlerSuccess(t *testing.T) {
	now := time.Now()

	expected := model.Task{
		ID:          1,
		Title:       "Learn Go",
		Description: "HTTP handler testing",
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	repo := &fakeRepo{
		task: expected,
	}

	request := httptest.NewRequest(
		"POST",
		"/tasks",
		strings.NewReader(`{"title":"Learn Go","description":"HTTP handler testing"}`),
	)
	recorder := httptest.NewRecorder()

	createTaskHandler(repo)(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("Expected: %v, Got: %v", http.StatusCreated, recorder.Code)
	}

	if recorder.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Wrong Header Content-Type")
	}

	var task model.Task
	body := io.NopCloser(strings.NewReader(recorder.Body.String()))

	if err := decodeJSON(body, &task); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if task.ID != expected.ID ||
		task.Title != expected.Title ||
		task.Description != expected.Description ||
		task.Completed != expected.Completed {
		t.Errorf("Response mismatch\nExpected: %+v\nGot: %+v", expected, task)
	}

	if task.CreatedAt.IsZero() || task.UpdatedAt.IsZero() {
		t.Errorf("Timestamps were not set")
	}
}

func (repo *fakeRepo) GetAll() ([]model.Task, error) {
	return repo.tasks, repo.err
}

// GET TASKS
func TestGetTasksHandlerSuccess(t *testing.T) {
	now := time.Now()

	expected := []model.Task{
		{
			ID:          1,
			Title:       "Learn Go",
			Description: "HTTP handler testing",
			Completed:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          2,
			Title:       "Learn PostgreSQL",
			Description: "Database testing",
			Completed:   true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	repo := &fakeRepo{
		tasks: expected,
	}

	request := httptest.NewRequest("GET", "/tasks", nil)
	recorder := httptest.NewRecorder()

	getTasksHandler(repo)(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("Expected: %v, Got: %v", http.StatusOK, recorder.Code)
	}

	if recorder.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Wrong Header Content-Type")
	}

	var tasks []model.Task
	body := io.NopCloser(strings.NewReader(recorder.Body.String()))

	if err := decodeJSON(body, &tasks); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(tasks) != len(expected) {
		t.Fatalf("Expected %d tasks, Got: %d", len(expected), len(tasks))
	}

	for i, task := range tasks {
		if task.ID != expected[i].ID ||
			task.Title != expected[i].Title ||
			task.Description != expected[i].Description ||
			task.Completed != expected[i].Completed {
			t.Errorf("Task %d mismatch\nExpected: %+v\nGot: %+v", i, expected[i], task)
		}
	}
}

// GET TASKS FAIL
func TestGetTasksHandlerFail(t *testing.T) {

	repo := &fakeRepo{
		err: errors.New("GET Err"),
	}

	request := httptest.NewRequest("GET", "/tasks", nil)
	recorder := httptest.NewRecorder()

	getTasksHandler(repo)(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("Expected: %v, Got: %v", http.StatusInternalServerError, recorder.Code)
	}
}

func (repo *fakeRepo) GetById(searchId int) (model.Task, error) {
	return repo.task, repo.err
}

// TASK LOOKUP
func TestGetTaskHandlerSuccess(t *testing.T) {

	now := time.Now()

	expected := model.Task{
		ID:          1,
		Title:       "Learn Go",
		Description: "HTTP handler testing",
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	repo := &fakeRepo{
		task: expected,
	}

	request := httptest.NewRequest("GET", "/tasks/1", nil) //id doesnt really matter here
	recorder := httptest.NewRecorder()

	getTaskHandler(repo)(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("Expected: %v, Got: %v", http.StatusOK, recorder.Code)
	}

	if recorder.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Wrong Header Content-Type")
	}

	var task model.Task
	body := io.NopCloser(strings.NewReader(recorder.Body.String()))

	if err := decodeJSON(body, &task); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if task.ID != expected.ID ||
		task.Title != expected.Title ||
		task.Description != expected.Description ||
		task.Completed != expected.Completed {
		t.Errorf("Task mismatch\nExpected: %+v\nGot: %+v", expected, task)
	}

}

// INVALID ID
func TestGetTaskHandlerInvalidID(t *testing.T) {
	request := httptest.NewRequest("GET", "/tasks/abc", nil)
	recorder := httptest.NewRecorder()

	getTaskHandler(nil)(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("Expected: %v, Got: %v", http.StatusBadRequest, recorder.Code)
	}
}

// TASK NOTFOUND
func TestGetTaskHandlerNotFound(t *testing.T) {

	repo := &fakeRepo{
		err: pgx.ErrNoRows,
	}

	request := httptest.NewRequest("GET", "/tasks/999", nil)
	recorder := httptest.NewRecorder()

	getTaskHandler(repo)(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("Expected: %v, Got: %v", http.StatusNotFound, recorder.Code)
	}
}

// TASK LOOKUP FAIL
func TestGetTaskHandlerFail(t *testing.T) {

	repo := &fakeRepo{
		err: errors.New("GET Err"),
	}

	request := httptest.NewRequest("GET", "/tasks/1", nil)
	recorder := httptest.NewRecorder()

	getTaskHandler(repo)(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("Expected: %v, Got: %v", http.StatusInternalServerError, recorder.Code)
	}
}

func (repo *fakeRepo) Delete(id int) (bool, error) {
	return repo.deleted, repo.err
}

// DELETE TASK
func TestDeleteTaskHandler(t *testing.T) {
	tests := []struct {
		name       string
		deleted    bool
		err        error
		wantStatus int
	}{
		{
			name:       "Success",
			deleted:    true,
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "Not Found",
			deleted:    false,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Repository Error",
			err:        errors.New("DELETE Err"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{
				deleted: tt.deleted,
				err:     tt.err,
			}

			request := httptest.NewRequest("DELETE", "/tasks/1", nil)
			recorder := httptest.NewRecorder()

			deleteTaskHandler(repo)(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Fatalf(
					"Expected: %v, Got: %v",
					tt.wantStatus,
					recorder.Code,
				)
			}
		})
	}
}

func (repo *fakeRepo) Update(id int, task model.Task) (model.Task, error) {
	return repo.task, repo.err
}

// UPDATE SUCCESS
func TestUpdateTaskHandlerSuccess(t *testing.T) {
	now := time.Now()

	expected := model.Task{
		ID:          1,
		Title:       "Learn Go",
		Description: "Updated description",
		Completed:   true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	repo := &fakeRepo{
		task: expected,
	}

	request := httptest.NewRequest(
		"PUT",
		"/tasks/1",
		strings.NewReader(`{
			"title": "Learn Go",
			"description": "Updated description",
			"completed": true
		}`),
	)
	recorder := httptest.NewRecorder()

	updateTaskHandler(repo)(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("Expected: %v, Got: %v", http.StatusOK, recorder.Code)
	}

	if recorder.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Wrong Header Content-Type")
	}

	var task model.Task
	body := io.NopCloser(strings.NewReader(recorder.Body.String()))

	if err := decodeJSON(body, &task); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if task.ID != expected.ID ||
		task.Title != expected.Title ||
		task.Description != expected.Description ||
		task.Completed != expected.Completed {
		t.Errorf("Task mismatch\nExpected: %+v\nGot: %+v", expected, task)
	}
}

// UPDATE FAIL
func TestUpdateTaskHandlerFail(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		err        error
		wantStatus int
	}{
		{
			name:       "Invalid JSON",
			body:       `{json`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Missing Title",
			body: `{
				"description": "Updated",
				"completed": true
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Missing Description",
			body: `{
				"title": "Learn Go",
				"completed": true
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Missing Completed",
			body: `{
				"title": "Learn Go",
				"description": "Updated"
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Unknown Field",
			body: `{
				"title": "Learn Go",
				"description": "Updated",
				"completed": true,
				"foo": "bar"
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Multiple JSON",
			body: `{
				"title": "Learn Go",
				"description": "Updated",
				"completed": true
			}{
				"title": "Another"
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Task Not Found",
			body:       `{"title":"Learn Go","description":"Updated","completed":true}`,
			err:        pgx.ErrNoRows,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Repository Error",
			body:       `{"title":"Learn Go","description":"Updated","completed":true}`,
			err:        errors.New("UPDATE Err"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{
				err: tt.err,
			}

			request := httptest.NewRequest(
				"PUT",
				"/tasks/1",
				strings.NewReader(tt.body),
			)
			recorder := httptest.NewRecorder()

			updateTaskHandler(repo)(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Fatalf(
					"Expected: %v, Got: %v",
					tt.wantStatus,
					recorder.Code,
				)
			}
		})
	}
}

func (repo *fakeRepo) Patch(id int, req model.TaskPatchRequest) (model.Task, error) {
	return repo.task, repo.err
}

// PATCH SUCCESS
func TestPatchTaskHandlerSuccess(t *testing.T) {
	now := time.Now()

	expected := model.Task{
		ID:          1,
		Title:       "Learn Go",
		Description: "Updated description",
		Completed:   true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	repo := &fakeRepo{task: expected}

	request := httptest.NewRequest(
		"PATCH",
		"/tasks/1",
		strings.NewReader(`{
            "title": "Learn Go",
            "description": "Updated description",
            "completed": true
        }`),
	)

	recorder := httptest.NewRecorder()

	patchTaskHandler(repo)(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("Expected: %v, Got: %v", http.StatusOK, recorder.Code)
	}

	if recorder.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected JSON Content-Type")
	}

	var task model.Task

	body := io.NopCloser(strings.NewReader(recorder.Body.String()))

	if err := decodeJSON(body, &task); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if task.ID != expected.ID ||
		task.Title != expected.Title ||
		task.Description != expected.Description ||
		task.Completed != expected.Completed {
		t.Errorf("Response mismatch\nExpected: %+v\nGot: %+v", expected, task)
	}

	if task.CreatedAt.IsZero() || task.UpdatedAt.IsZero() {
		t.Errorf("Timestamps were not set")
	}
}

// PATCH FAIL
func TestPatchTaskHandlerFail(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		body       string
		err        error
		wantStatus int
	}{
		{
			name:       "Invalid JSON",
			path:       "/tasks/1",
			body:       `{json`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Unknown Field",
			path: "/tasks/1",
			body: `{
                "title": "Learn Go",
                "foo": "bar"
            }`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Multiple JSON",
			path:       "/tasks/1",
			body:       `{"title":"Learn Go"}{"completed":true}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "No Fields",
			path:       "/tasks/1",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Invalid ID",
			path:       "/tasks/abc",
			body:       `{"title":"Learn Go"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Task Not Found",
			path:       "/tasks/999",
			body:       `{"title":"Learn Go"}`,
			err:        pgx.ErrNoRows,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Repository Error",
			path:       "/tasks/1",
			body:       `{"title":"Learn Go"}`,
			err:        errors.New("PATCH Err"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{err: tt.err}

			request := httptest.NewRequest(
				"PATCH",
				tt.path,
				strings.NewReader(tt.body),
			)

			recorder := httptest.NewRecorder()

			patchTaskHandler(repo)(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Fatalf(
					"Expected: %v, Got: %v",
					tt.wantStatus,
					recorder.Code,
				)
			}
		})
	}
}
