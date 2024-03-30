package app

import (
	"fmt"
	"konzek-jun/globalerror"
	"konzek-jun/loggerx"
	"konzek-jun/models"
	"konzek-jun/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TaskHandler struct {
	Service      services.TaskService
	WorkerPool   chan struct{}
	MaxWorkerNum int
}

func NewTaskHandler(service services.TaskService, maxWorkerNum int) *TaskHandler {
	return &TaskHandler{
		Service:      service,
		MaxWorkerNum: maxWorkerNum,
		WorkerPool:   make(chan struct{}, maxWorkerNum),
	}
}

func (h *TaskHandler) acquireWorker() {
	h.WorkerPool <- struct{}{}
}

func (h *TaskHandler) releaseWorker() {
	<-h.WorkerPool
}

// @Summary Retrieves all tasks
// @Description Retrieves all tasks
// @Tags Tasks
// @Accept json
// @Produce json
// @Success 200 {object} []models.Task "List of tasks"
// @Failure 500 {object} globalerror.ErrorResponse "Internal server error"
// @Router /tasks [get]
func (h *TaskHandler) GetAllTask(c *fiber.Ctx) error {
	loggerx.Info("GetAllTask function called")

	resultChan := make(chan *[]models.Task)

	defer close(resultChan)

	go func() {

		h.acquireWorker()
		defer func() {
			h.releaseWorker()

		}()

		result, err := h.Service.TaskGetAll()

		if err != nil {
			resultChan <- nil
			return
		}

		resultChan <- &result
	}()
	result := <-resultChan
	if result == nil {

		log.Println("Error fetching tasks")

		return c.Status(http.StatusInternalServerError).JSON(globalerror.ErrorResponse{
			Status: http.StatusInternalServerError,
			ErrorDetail: []globalerror.ErrorResponseDetail{
				{
					FieldName:   "Task",
					Description: "An error occurred while fetching the tasks",
				},
			},
		})
	}

	loggerx.Info("Tasks fetched successfully")
	return c.Status(http.StatusOK).JSON(result)

}

// @Summary Creates a new task
// @Description Creates a new task
// @Tags Tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Task object to create"
// @Success 201 {object} EmptyResponse "Empty response"
// @Failure 400 {object} globalerror.ErrorResponse "Bad request"
// @Failure 500 {object} globalerror.ErrorResponse "Internal server error"
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	loggerx.Info("CreateTask function called")
	var task models.Task
	resultChan := make(chan *[]models.Task)
	defer close(resultChan)

	if err := c.BodyParser(&task); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	if errors := globalerror.Validate(task); len(errors) > 0 && errors[0].HasError {
		return globalerror.HandleValidationErrors(c, errors)
	}
	go func() {
		h.acquireWorker()
		defer h.releaseWorker()

		err := h.Service.TaskInsert(task)
		if err != nil {
			resultChan <- nil
			return

		}
		resultChan <- &[]models.Task{}

	}()
	result := <-resultChan

	if result == nil {
		return c.Status(http.StatusInternalServerError).JSON(globalerror.ErrorResponse{
			Status: http.StatusInternalServerError,
			ErrorDetail: []globalerror.ErrorResponseDetail{
				{
					FieldName:   "Task",
					Description: " An error occurred while adding the task",
				},
			},
		})
	}
	loggerx.Info("Task created successfully")

	return c.Status(http.StatusCreated).JSON(nil)
}

// @Summary Deletes a task by its ID
// @Description Deletes a task by its ID
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path integer true "Task ID to delete"
// @Success 200 {object} EmptyResponse "Empty response"
// @Failure 400 {object} globalerror.ErrorResponse "Bad request"
// @Failure 500 {object} globalerror.ErrorResponse "Internal server error"
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	loggerx.Info("DeleteTask function called")
	resultChan := make(chan *[]models.Task)
	defer close(resultChan)

	strId := c.Params("id")
	id, err := strconv.Atoi(strId)
	fmt.Println(id)
	if err != nil {
		return err
	}

	go func() {
		h.acquireWorker()
		defer h.releaseWorker()

		err := h.Service.TaskDelete(id)
		if err != nil {
			resultChan <- nil
			return
		}
		resultChan <- &[]models.Task{}
	}()

	result := <-resultChan

	if result == nil {
		return c.Status(http.StatusBadRequest).JSON(globalerror.ErrorResponse{
			Status: http.StatusBadRequest,
			ErrorDetail: []globalerror.ErrorResponseDetail{
				{
					FieldName:   "Task",
					Description: "An error occurred while deleting the task",
				},
			},
		})
	}
	loggerx.Info("Task deleted successfully")

	return c.Status(http.StatusOK).JSON(fiber.Map{"success": true})
}

// @Summary Updates an existing task
// @Description Updates an existing task
// @Tags Tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Updated task object"
// @Success 201 {object} EmptyResponse "Empty response"
// @Failure 400 {object} globalerror.ErrorResponse "Bad request"
// @Failure 500 {object} globalerror.ErrorResponse "Internal server error"
// @Router /tasks [put]
func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	loggerx.Info("UpdateTask function called")
	var updatedTask models.Task
	resultChan := make(chan *[]models.Task)
	defer close(resultChan)

	if err := c.BodyParser(&updatedTask); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Geçersiz gövde")
	}

	if errors := globalerror.Validate(updatedTask); len(errors) > 0 && errors[0].HasError {
		return globalerror.HandleValidationErrors(c, errors)
	}
	go func() {
		h.acquireWorker()
		defer h.releaseWorker()

		err := h.Service.TaskUpdate(updatedTask)
		if err != nil {
			resultChan <- nil
			return
		}
		resultChan <- &[]models.Task{}
	}()
	result := <-resultChan
	if result == nil {
		return c.Status(http.StatusInternalServerError).JSON(globalerror.ErrorResponse{
			Status: http.StatusInternalServerError,
			ErrorDetail: []globalerror.ErrorResponseDetail{
				{
					FieldName:   "Task",
					Description: "An error occurred while updating the task",
				},
			},
		})
	}

	loggerx.Info("Task updated successfully")
	return c.Status(http.StatusOK).JSON(fiber.Map{"success": true})
}

// @Summary Retrieves a task by its ID
// @Description Retrieves a task by its ID
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path integer true "Task ID to retrieve"
// @Success 200 {object} models.Task "Task object"
// @Failure 404 {object} globalerror.ErrorResponse "Not found"
// @Failure 500 {object} globalerror.ErrorResponse "Internal server error"
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetByID(c *fiber.Ctx) error {
	loggerx.Info("GetByID function called")

	strID := c.Params("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(globalerror.ErrorResponse{
			Status: http.StatusNotFound,
			ErrorDetail: []globalerror.ErrorResponseDetail{
				{
					FieldName:   "Task",
					Description: "invalid task id",
				},
			},
		})
	}

	resultChan := make(chan *models.Task)

	defer close(resultChan)

	go func() {
		h.acquireWorker()
		defer h.releaseWorker()

		model, err := h.Service.TaskGetByID(id)

		if err != nil {
			resultChan <- nil
			return
		}

		resultChan <- &model
	}()

	result := <-resultChan
	if result != nil {
		loggerx.Info("Task loaded successfully")
		return c.Status(http.StatusOK).JSON(result)
	} else {
		return c.Status(http.StatusInternalServerError).JSON(globalerror.ErrorResponse{
			Status: http.StatusInternalServerError,
			ErrorDetail: []globalerror.ErrorResponseDetail{
				{
					FieldName:   "Task",
					Description: "An error occurred while loading the task",
				},
			},
		})
	}
}

// @Summary Retrieves all tasks with pagination
// @Description Retrieves all tasks with pagination
// @Tags Tasks
// @Accept json
// @Produce json
// @Param page query integer false "Page number"
// @Param pageSize query integer false "Number of tasks per page"
// @Success 200 {object} EmptyResponse "Empty response"
// @Failure 400 {object} globalerror.ErrorResponse "Bad request"
// @Failure 500 {object} globalerror.ErrorResponse "Internal server error"
// @Router /tasks/page [get]
func (h *TaskHandler) GetAllTaskWithPagination(c *fiber.Ctx) error {

	loggerx.Info("GetAllTaskWithPagination function called")

	params := new(PaginationParams)
	if err := c.QueryParser(params); err != nil {
		return c.Status(http.StatusBadRequest).JSON(globalerror.ErrorResponse{
			Status: http.StatusBadRequest,
			ErrorDetail: []globalerror.ErrorResponseDetail{
				{
					FieldName:   "Pagination",
					Description: "Invalid pagination parameters",
				},
			},
		})
	}

	tasks, err := h.Service.GetAllTaskWithPagination(params.Page, params.PageSize)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(globalerror.ErrorResponse{
			Status: http.StatusInternalServerError,
			ErrorDetail: []globalerror.ErrorResponseDetail{
				{
					FieldName:   "Task",
					Description: "An error occurred while fetching the tasks",
				},
			},
		})
	}

	loggerx.Info("Tasks fetched successfully")

	return c.JSON(fiber.Map{
		"tasks": tasks,
	})
}

type PaginationParams struct {
	Page     int `query:"page"`
	PageSize int `query:"pageSize"`
}

type EmptyResponse struct {
	Success bool `json:"success"`
}
