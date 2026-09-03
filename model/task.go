package model

import "time"

type Task struct {
	ID				int			`json:"id"`
	Title			string		`json:"title"`
	Description 	string		`json:"description"`
	Completed		bool		`json:"completed"`
	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at"`
}

type TaskRequest struct {
	Title			string		`json:"title"`
	Description		string		`json:"description"`
}

type TaskUpdateRequest struct {
	Title			*string		`json:"title"`
	Description		*string		`json:"description"`
	Completed		*bool		`json:"completed"`
}

type TaskPatchRequest struct {
	Title			*string		`json:"title"`
	Description		*string		`json:"description"`
	Completed		*bool		`json:"completed"`
}
