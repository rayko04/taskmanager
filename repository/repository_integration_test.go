package repository

import (
	"context"
	"os"
	"taskmanager/database"
	"taskmanager/model"
	"testing"
	"time"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func initialize() (*pgxpool.Pool, *TaskRepository, error) {

	err := godotenv.Load("../.env")
	if err != nil {
		return nil, nil, err
	}

	pool, err := database.NewPool(os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		return nil, nil, err
	}

	_, err = pool.Exec(context.Background(), "DELETE FROM tasks")
	if err != nil {
		pool.Close()
    	return nil, nil, err
	}

	repo := NewTaskRepository(pool)
	return pool, repo, nil
}

func TestCreate(t *testing.T) {

	pool, repo, err := initialize()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	title := "Test"
	desc := "test"

	task := model.Task{
		Title: title,
		Description: desc,
	}
	task, err =  repo.Create(task)
	if err != nil {
		t.Fatal(err)
	}

	if task.ID == 0 {
		t.Errorf("ID returned 0")
	} 
	if task.Description != desc {
		t.Errorf("Mismatched Description")
	} 
	if task.Title != title {
		t.Errorf("Mismatched Title")
	} 
	if task.Completed != false {
		t.Error("Completed returned true")
	}
	if task.CreatedAt.IsZero() {
		t.Errorf("CreatedAt not set")
	} 
	if task.UpdatedAt.IsZero() {
		t.Errorf("UpdatedAt not set")
	}

	task, err = repo.GetById(task.ID)
	if err != nil {
		t.Fatal(err)
	}

	if task.ID == 0 {
		t.Errorf("ID returned 0")
	} 
	if task.Description != desc {
		t.Errorf("Mismatched Description")
	} 
	if task.Title != title {
		t.Errorf("Mismatched Title")
	} 
	if task.Completed != false {
		t.Error("Completed returned true")
	}
	if task.CreatedAt.IsZero() {
		t.Errorf("CreatedAt not set")
	} 
	if task.UpdatedAt.IsZero() {
		t.Errorf("UpdatedAt not set")
	}

}


func TestGetAll(t *testing.T) {

	pool, repo, err := initialize()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	title := "Test"
	desc := "test"

	task := model.Task {
		Title: title,
		Description: desc,
	}

	expectedCount := 5
	for i := 0; i < expectedCount; i++ {
		
	 	_, err = repo.Create(task)
		if err != nil {
			t.Fatal(err)
		}
	}

	tasks, err := repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) != expectedCount {
		t.Errorf("got %d tasks, want %d", len(tasks), expectedCount)
	}

	for _, task := range tasks {
		if task.ID == 0 {
		t.Errorf("ID returned 0")
		} 
		if task.Description != desc {
			t.Errorf("Mismatched Description")
		} 
		if task.Title != title {
			t.Errorf("Mismatched Title")
		} 
		if task.Completed != false {
			t.Error("Completed returned true")
		}
		if task.CreatedAt.IsZero() {
			t.Errorf("CreatedAt not set")
		} 
		if task.UpdatedAt.IsZero() {
			t.Errorf("UpdatedAt not set")
		}
	}

}


func TestGetById(t *testing.T) {

	pool, repo, err := initialize()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	title := "Test"
	desc := "test"

	task := model.Task {
		Title: title,
		Description: desc,
	}

	task, err = repo.Create(task)
	if err != nil {
		t.Fatal(err)
	}

	task, err = repo.GetById(task.ID)
	if err != nil {
		t.Fatal(err)
	}

	if task.ID == 0 {
	t.Errorf("ID returned 0")
	} 
	if task.Description != desc {
		t.Errorf("Mismatched Description")
	} 
	if task.Title != title {
		t.Errorf("Mismatched Title")
	} 
	if task.Completed != false {
		t.Error("Completed returned true")
	}
	if task.CreatedAt.IsZero() {
		t.Errorf("CreatedAt not set")
	} 
	if task.UpdatedAt.IsZero() {
		t.Errorf("UpdatedAt not set")
	}

}

func TestGetByIdNotFound(t *testing.T) {

	pool, repo, err := initialize()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	_, err = repo.GetById(9999)
	if err != pgx.ErrNoRows {
		t.Fatalf("expected: %v, got: %v", pgx.ErrNoRows, err)
	}
}

func TestUpdate(t *testing.T) {

	pool, repo, err := initialize()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	task := model.Task {
		Title: "abc",
		Description: "efg",
	}

	task, err = repo.Create(task)
	if err != nil {
		t.Fatal(err)
	}

	title := "Test"
	desc := "test"
	comp := true
	id := task.ID
	createTime := task.CreatedAt
	updateTime := task.UpdatedAt

	update := model.Task{
		Title: title,
		Description: desc,
		Completed: comp,
	}

	task, err = repo.Update(task.ID, update)
	if err != nil {
		t.Fatal(err)
	}
	if task.ID == 0 || task.ID != id {
    	t.Errorf("ID changed: expected %d, got %d", id, task.ID)
	}
	if task.Description != desc {
		t.Errorf("Description mismatch: expected %q, got %q", desc, task.Description)
	}
	if task.Title != title {
		t.Errorf("Title mismatch: expected %q, got %q", title, task.Title)
	}
	if task.Completed != comp {
		t.Errorf("Completed mismatch: expected %t, got %t", comp, task.Completed)
	}
	if !task.CreatedAt.Truncate(time.Microsecond).Equal(createTime.Truncate(time.Microsecond)) {
    	t.Errorf("CreatedAt changed: expected %v, got %v", createTime, task.CreatedAt)
	}
	if !task.UpdatedAt.After(updateTime) {
		t.Errorf("UpdatedAt was not advanced: previous %v, got %v", updateTime, task.UpdatedAt)
	}

}


func TestPatch(t *testing.T) {

	pool, repo, err := initialize()
	if err != nil {
		t.Fatalf("failed to initialize test database: %v", err)
	}
	defer pool.Close()

	tests := []struct {
		name              string
		initialTitle      string
		initialDesc       string
		initialCompleted  bool
		patch             model.TaskPatchRequest
		expectedTitle     string
		expectedDesc      string
		expectedCompleted bool
	}{
		{
			name:              "single field",
			initialTitle:      "abc",
			initialDesc:       "efg",
			initialCompleted:  false,
			patch: model.TaskPatchRequest{
				Title: new("Test"),
			},
			expectedTitle:     "Test",
			expectedDesc:      "efg",
			expectedCompleted: false,
		},
		{
			name:              "multiple fields",
			initialTitle:      "abc",
			initialDesc:       "efg",
			initialCompleted:  false,
			patch: model.TaskPatchRequest{
				Title:     new("Test"),
				Completed: new(true),
			},
			expectedTitle:     "Test",
			expectedDesc:      "efg",
			expectedCompleted: true,
		},
		{
			name:              "explicit false",
			initialTitle:      "abc",
			initialDesc:       "efg",
			initialCompleted:  true,
			patch: model.TaskPatchRequest{
				Completed: new(false),
			},
			expectedTitle:     "abc",
			expectedDesc:      "efg",
			expectedCompleted: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := pool.Exec(context.Background(), "DELETE FROM tasks")
			if err != nil {
				t.Fatalf("failed to clear tasks: %v", err)
			}

			task := model.Task{
				Title:       tt.initialTitle,
				Description: tt.initialDesc,
				Completed:   tt.initialCompleted,
			}

			task, err = repo.Create(task)
			if err != nil {
				t.Fatalf("failed to create task: %v", err)
			}

			id := task.ID
			createTime := task.CreatedAt
			updateTime := task.UpdatedAt

			task, err = repo.Patch(task.ID, tt.patch)
			if err != nil {
				t.Fatalf("failed to patch task: %v", err)
			}

			if task.ID != id {
				t.Errorf("ID changed: expected %d, got %d", id, task.ID)
			}

			if task.Title != tt.expectedTitle {
				t.Errorf("Title mismatch: expected %q, got %q",
					tt.expectedTitle, task.Title)
			}

			if task.Description != tt.expectedDesc {
				t.Errorf("Description mismatch: expected %q, got %q",
					tt.expectedDesc, task.Description)
			}

			if task.Completed != tt.expectedCompleted {
				t.Errorf("Completed mismatch: expected %t, got %t",
					tt.expectedCompleted, task.Completed)
			}

			if !task.CreatedAt.Truncate(time.Microsecond).
				Equal(createTime.Truncate(time.Microsecond)) {
				t.Errorf("CreatedAt changed: expected %v, got %v",
					createTime, task.CreatedAt)
			}

			if !task.UpdatedAt.After(updateTime) {
				t.Errorf("UpdatedAt was not advanced: previous %v, got %v",
					updateTime, task.UpdatedAt)
			}
		})
	}
}

func TestPatchNotFound(t *testing.T) {

	pool, repo, err := initialize()
	if err != nil {
		t.Fatalf("failed to initialize test database: %v", err)
	}
	defer pool.Close()

	title := "Test"

	patch := model.TaskPatchRequest{
		Title: new(title),
	}

	_, err = repo.Patch(9999, patch)
	if err != pgx.ErrNoRows {
		t.Fatalf("expected pgx.ErrNoRows, got %v", err)
	}
}

func TestDelete(t *testing.T) {

	pool, repo, err := initialize()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	task := model.Task {
		Title: "abc",
		Description: "efg",
	}

	task, err = repo.Create(task)
	if err != nil {
		t.Fatal(err)
	}

	deleted, err := repo.Delete(task.ID)
	if err != nil {
		t.Errorf("Error while Deletion: %v", err)
	}
	if !deleted {
		t.Errorf("task was not deleted: expected true, got false")
	}

	_, err = repo.GetById(task.ID)
	if err != pgx.ErrNoRows {
		t.Errorf("expected pgx.ErrNoRows, got %v", err)
	}
}

func TestDeleteNotFound(t *testing.T) {

	pool, repo, err := initialize()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	deleted, err := repo.Delete(9999)
	if err != nil {
		t.Errorf("Error while Deletion: %v", err)
	}
	if deleted {
		t.Errorf("expected false for nonexistent task, got true")
	}
}
