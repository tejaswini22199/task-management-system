package services

import (
	"github.com/tejaswini22199/task-management-system/taskservice/models"
	"github.com/tejaswini22199/task-management-system/taskservice/repository"
)

func CreateTask(input models.TaskInput, userId int) (models.Task, error) {
	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		CreatedBy:   userId,
	}

	return repository.CreateTask(task, input.UserIDs) // Pass user IDs
}

func GetTasks() ([]models.Task, error) {
	return repository.GetTasks()
}

func GetTaskByID(id int) (models.Task, error) {
	return repository.GetTaskByID(id)
}

func GetTasksByStatus(status string) ([]models.Task, error) {
	return repository.GetTasksByStatus(status)
}

func UpdateTask(id int, input models.TaskInput) (models.Task, error) {
	task, err := repository.GetTaskByID(id)
	if err != nil {
		return models.Task{}, err
	}

	task.Title = input.Title
	task.Description = input.Description
	task.Status = input.Status

	return repository.UpdateTask(task)
}

func DeleteTask(id int) error {
	return repository.DeleteTask(id)
}
