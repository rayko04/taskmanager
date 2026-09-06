package repository

import (
	"context"
	"taskmanager/model"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	pool 		*pgxpool.Pool
}

func NewTaskRepository(pl *pgxpool.Pool) *TaskRepository{
	return &TaskRepository {
		pool:  pl,
	}
}

func (repo *TaskRepository) Create(task model.Task) (model.Task, error) {

	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	err := repo.pool.QueryRow(
		context.Background(), 
		"INSERT INTO tasks (title, description, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id, completed", 
		task.Title, task.Description, task.CreatedAt, task.UpdatedAt,
	).Scan(&task.ID, &task.Completed)
	
	return task, err
}

func (repo *TaskRepository) GetAll() ([]model.Task, error) {

	rows, err := repo.pool.Query(
		context.Background(),
		"SELECT id, title, description, completed, created_at, updated_at FROM tasks",
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
		return nil, err
	}

	return slice, nil
}

func (repo *TaskRepository) GetById(searchId int) (model.Task, error) {

	task := model.Task{}
	err := repo.pool.QueryRow(
		context.Background(),
		"SELECT id, title, description, completed, created_at, updated_at FROM tasks WHERE id = $1",
		searchId, 
	).Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.CreatedAt, &task.UpdatedAt)

	return task, err
}

func (repo *TaskRepository) Delete(searchId int) (bool, error) {

	result, err := repo.pool.Exec(
		context.Background(),
		"DELETE FROM tasks WHERE id = $1",
		searchId,
	)

	if err != nil {
    	return false, err
	}

	return result.RowsAffected() > 0, nil
}

func (repo *TaskRepository) Update(searchId int, task model.Task) (model.Task, error) {

	err := repo.pool.QueryRow(
		context.Background(),
		"UPDATE tasks SET title = $1, description = $2, completed = $3, updated_at = $4 WHERE id = $5 RETURNING id, title, description, completed, created_at, updated_at;",
		task.Title, task.Description, task.Completed, time.Now(), searchId,
	).Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.CreatedAt, &task.UpdatedAt)

	return task, err
}

func (repo *TaskRepository) Patch(searchId int, req model.TaskPatchRequest) (model.Task, error) {

	var title, desc, comp any
	if req.Title != nil {
		title = *req.Title	//0-value for any is nil
	}

	if req.Description != nil {
		desc = *req.Description
	}

	if req.Completed != nil {
		comp = *req.Completed
	}

	task := model.Task{}
	err := repo.pool.QueryRow(
		context.Background(),
		"UPDATE tasks SET title = COALESCE($1, title), description = COALESCE($2, description), completed = COALESCE($3, completed), updated_at = $4 WHERE id = $5 RETURNING id, title, description, completed, created_at, updated_at;",
		title, desc, comp, time.Now(), searchId,
	).Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.CreatedAt, &task.UpdatedAt)

	return task, err
}
