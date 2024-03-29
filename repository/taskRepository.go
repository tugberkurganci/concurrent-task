package repository

import (
	"context"
	"database/sql"
	"fmt"
	"konzek-jun/loggerx"
	"konzek-jun/models"
	"time"

	_ "github.com/lib/pq"
)

type TaskRepositoryDb struct {
	DB *sql.DB
}

type TaskRepository interface {
	Insert(todo models.Task) (int64, error)
	GetAll() ([]models.Task, error)
	Delete(id int) error
	GetByID(id int) (models.Task, error)
	Update(task models.Task) error
	GetTasksWithPagination(offset, limit int) ([]models.Task, error)
}

func NewTaskRepository(db *sql.DB) *TaskRepositoryDb {
	return &TaskRepositoryDb{DB: db}
}

func (t *TaskRepositoryDb) Insert(task models.Task) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var lastInsertID int64

	err := withRetry(func() error {
		err := t.DB.QueryRowContext(ctx, "INSERT INTO tasks (title, content, status) VALUES ($1, $2, $3) RETURNING id", task.Title, task.Content, task.Status).Scan(&lastInsertID)

		if err != nil {
			loggerx.Error(fmt.Sprintf("Error while inserting task: %v", err))
			return err
		}

		loggerx.Info("Task inserted successfully")
		return nil
	})
	return lastInsertID, err
}

func (t *TaskRepositoryDb) GetAll() ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var tasks []models.Task
	err := withRetry(func() error {
		rows, err := t.DB.QueryContext(ctx, "SELECT id, title, content, status FROM tasks")
		if err != nil {
			loggerx.Error(fmt.Sprintf("Error while getting all tasks: %v", err))
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var task models.Task
			err := rows.Scan(&task.Id, &task.Title, &task.Content, &task.Status)
			if err != nil {
				return err
			}
			tasks = append(tasks, task)
		}

		loggerx.Info("Retrieved all tasks successfully")
		return nil
	})
	return tasks, err
}

func (t *TaskRepositoryDb) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := withRetry(func() error {
		_, err := t.DB.ExecContext(ctx, "DELETE FROM tasks WHERE id = $1", id)
		if err != nil {
			loggerx.Error(fmt.Sprintf("Error while deleting task: %v", err))
			return err
		}
		loggerx.Info("Task deleted successfully")
		return nil
	})
	return err
}

func (t *TaskRepositoryDb) GetByID(id int) (models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var task models.Task
	err := withRetry(func() error {
		err := t.DB.QueryRowContext(ctx, "SELECT id, title, content, status FROM tasks WHERE id = $1", id).Scan(&task.Id, &task.Title, &task.Content, &task.Status)
		if err != nil {
			loggerx.Error(fmt.Sprintf("Error while getting task by ID: %v", err))
			return err
		}
		loggerx.Info("Retrieved task by ID successfully")
		return nil
	})
	return task, err
}

func (t *TaskRepositoryDb) Update(task models.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := withRetry(func() error {
		_, err := t.DB.ExecContext(ctx, "UPDATE tasks SET title = $1, content = $2, status = $3 WHERE id = $4", task.Title, task.Content, task.Status, task.Id)
		if err != nil {
			loggerx.Error(fmt.Sprintf("Error while updating task: %v", err))
			return err
		}
		loggerx.Info("Task updated successfully")
		return nil
	})
	return err
}

func withRetry(operation func() error) error {
	var err error
	for i := 0; i < 3; i++ {
		err = operation()
		if err == nil {
			return nil
		}
		loggerx.Error(fmt.Sprintf("Error occurred, retrying: %v", err))
		time.Sleep(100 * time.Millisecond)
	}
	return err
}

func (t *TaskRepositoryDb) GetTasksWithPagination(offset, limit int) ([]models.Task, error) {
	query := "SELECT * FROM tasks ORDER BY id LIMIT $1 OFFSET $2"
	rows, err := t.DB.Query(query, limit, offset)
	if err != nil {
		loggerx.Error(fmt.Sprintf("Error while getting tasks with pagination: %v", err))
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.Id, &task.Title, &task.Content, &task.Status); err != nil {
			loggerx.Error(fmt.Sprintf("Error while scanning task: %v", err))
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		loggerx.Error(fmt.Sprintf("Error occurred while iterating over rows: %v", err))
		return nil, err
	}
	loggerx.Info("Retrieved tasks with pagination successfully")
	return tasks, nil
}
