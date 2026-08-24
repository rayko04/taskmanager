package repository

import (
	"taskmanager/model"
	"time"
)

type TaskRepository struct {
	tasks		map[int]model.Task
	nextId	int
}

func NewTaskRepository() *TaskRepository{
	return &TaskRepository {
		tasks: map[int]model.Task{},
		nextId: 1,
	}
}

func (repo *TaskRepository) Create(task model.Task) model.Task {
	
	task.ID = repo.nextId
	task.Completed = false
	
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	repo.tasks[repo.nextId] = task
	repo.nextId += 1

	return task
}

func (repo TaskRepository) GetAll() []model.Task {
	slice := []model.Task{}

	for _, task := range repo.tasks {
		slice = append(slice, task)
	}

	return slice
}

func (repo TaskRepository) GetById(searchId int) (model.Task, bool) {
	task, exists := repo.tasks[searchId]
	return task, exists
}
