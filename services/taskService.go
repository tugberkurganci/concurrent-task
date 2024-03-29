package services

import (
	"fmt"
	"konzek-jun/loggerx"
	"konzek-jun/models"
	"konzek-jun/repository"
)

//go:generate mockgen -destination=../mocks//service/mockTaskservice.go -package=services konzek-jun/services TaskService
type TaskService interface {
	TaskInsert(Task models.Task) error
	TaskGetAll() ([]models.Task, error)
	TaskDelete(id int) error
	TaskUpdate(task models.Task) error
	TaskGetByID(id int) (models.Task, error)
	GetAllTaskWithPagination(page, pageSize int) ([]models.Task, error)
}

type DefaultTaskService struct {
	Repo repository.TaskRepository
}

func NewTaskService(Repo repository.TaskRepository) DefaultTaskService {

	return DefaultTaskService{
		Repo: Repo,
	}
}

func (t DefaultTaskService) TaskInsert(task models.Task) error {
	_, err := t.Repo.Insert(task)
	if err != nil {
		loggerx.Error(fmt.Sprintf("Error while inserting task: %s", err))
		return err
	}
	loggerx.Info("Task inserted successfully")
	return nil
}

func (t DefaultTaskService) TaskGetAll() ([]models.Task, error) {
	result, err := t.Repo.GetAll()
	if err != nil {
		loggerx.Error(fmt.Sprintf("Error while getting all tasks: %s", err))
		return nil, err
	}
	loggerx.Info("Retrieved all tasks successfully")
	return result, nil
}

func (t DefaultTaskService) TaskDelete(id int) error {
	err := t.Repo.Delete(id)
	if err != nil {
		loggerx.Error(fmt.Sprintf("Error while deleting task: %s", err))
		return err
	}
	loggerx.Info("Task deleted successfully")
	return nil
}

func (t DefaultTaskService) TaskUpdate(task models.Task) error {
	err := t.Repo.Update(task)
	if err != nil {
		loggerx.Error(fmt.Sprintf("Error while updating task: %s", err))
		return err
	}
	loggerx.Info("Task updated successfully")
	return nil
}

func (t DefaultTaskService) TaskGetByID(id int) (models.Task, error) {
	task, err := t.Repo.GetByID(id)
	if err != nil {
		loggerx.Error(fmt.Sprintf("Error while getting task by ID: %s", err))
		return models.Task{}, err
	}
	loggerx.Info("Retrieved task by ID successfully")
	return task, nil
}

func (s DefaultTaskService) GetAllTaskWithPagination(page, pageSize int) ([]models.Task, error) {
	offset := (page - 1) * pageSize
	limit := pageSize
	tasks, err := s.Repo.GetTasksWithPagination(offset, limit)
	if err != nil {
		loggerx.Error(fmt.Sprintf("Error while getting tasks with pagination: %s", err))
		return nil, err
	}
	loggerx.Info("Retrieved tasks with pagination successfully")
	return tasks, nil
}
