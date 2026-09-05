package repository

import (
	"sync"
	"taskmanager/model"
	"testing"
)

func TestConcurrentCreate(t *testing.T) {

	var wg sync.WaitGroup
	wg.Add(100)

	repo := NewTaskRepository()

	for i:= 0; i < 100; i++ {
		go func() {
			defer wg.Done()

			repo.Create(model.Task{Title: "test",})
		}()
	}

	wg.Wait()

	tasks := repo.GetAll()
	if len(tasks) != 100 {
		t.Errorf("expected 100 tasks, got %d", len(tasks))
	}
}
