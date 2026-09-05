package repository

import (
	"taskmanager/model"
	"time"
	"sync"
)

type TaskRepository struct {
	mu 			sync.RWMutex
	tasks		map[int]model.Task
	nextId		int
}

func NewTaskRepository() *TaskRepository{
	return &TaskRepository {
		tasks: map[int]model.Task{},
		nextId: 1,
	}
}

func (repo *TaskRepository) Create(task model.Task) model.Task {

	repo.mu.Lock()
	defer repo.mu.Unlock()
	
	task.ID = repo.nextId
	task.Completed = false
	
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	repo.tasks[repo.nextId] = task
	repo.nextId += 1

	return task
}

func (repo *TaskRepository) GetAll() []model.Task {

	repo.mu.RLock()
	defer repo.mu.RUnlock()

	slice := []model.Task{}

	for _, task := range repo.tasks {
		slice = append(slice, task)
	}

	return slice
}

func (repo *TaskRepository) GetById(searchId int) (model.Task, bool) {
	
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	task, exists := repo.tasks[searchId]

	return task, exists
}

func (repo *TaskRepository) Delete(searchId int) bool {

	repo.mu.Lock()
	defer repo.mu.Unlock()

	_, exists := repo.tasks[searchId]
	if !exists {
		return false
	}
	
	delete(repo.tasks, searchId)

	return true
}

func (repo *TaskRepository) Update(searchId int, task model.Task) (model.Task, bool) {

	repo.mu.Lock()
	defer repo.mu.Unlock()

	_, exists := repo.tasks[searchId]
	if !exists {
		return model.Task{}, false
	}
	task.ID = searchId
	task.CreatedAt = repo.tasks[searchId].CreatedAt
	task.UpdatedAt = time.Now()

	repo.tasks[searchId] = task
	
	return task, true
}

func (repo *TaskRepository) Patch(searchId int, req model.TaskPatchRequest) (model.Task, bool) {

	repo.mu.Lock()
	defer repo.mu.Unlock()

	task, exists := repo.tasks[searchId]
	if !exists {
		return model.Task{}, false
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Completed != nil {
		task.Completed = *req.Completed
	}

	task.UpdatedAt = time.Now()
	repo.tasks[searchId] = task
	
	return task, true
}
