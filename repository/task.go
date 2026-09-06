package repository

import (
	"context"
	"sync"
	"taskmanager/model"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	mu 			sync.RWMutex
	pool 		*pgxpool.Pool
	tasks		map[int]model.Task
	nextId		int
}

func NewTaskRepository(pl *pgxpool.Pool) *TaskRepository{
	return &TaskRepository {
		pool:  pl,
		tasks: map[int]model.Task{},
		nextId: 1,
	}
}

func (repo *TaskRepository) Create(task model.Task) (model.Task, error) {

	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	err := repo.pool.QueryRow(
		context.Background(), 
		"INSERT INTO tasks (title, description, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id", 
		task.Title, task.Description, task.CreatedAt, task.UpdatedAt,
	).Scan(&task.ID)
	
	return task, err
}

func (repo *TaskRepository) GetAll() ([]model.Task, error) {

	rows, err := repo.pool.Query(
		context.Background(),
		"SELECT * FROM tasks",
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	
	slice := []model.Task{}
	for rows.Next() {
		var task model.Task

		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}

		slice = append(slice, task)
	}

	if err := rows.Err(); err != nil {
		return nil, nil
	}

	return slice, err
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
