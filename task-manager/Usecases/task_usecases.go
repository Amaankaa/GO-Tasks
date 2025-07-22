package usecases

import (
	"task-manager/Domain"
)

type TaskUsecase struct {
	taskRepo domain.TaskRepository
}

func NewTaskUsecase(taskRepo domain.TaskRepository) *TaskUsecase {
	return &TaskUsecase{
		taskRepo: taskRepo,
	}
}

func (tu *TaskUsecase) GetAllTasks() ([]domain.Task, error) {
	return tu.taskRepo.GetAllTasks()
}

func (tu *TaskUsecase) GetTaskByID(id string) (domain.Task, error) {
	return tu.taskRepo.GetTaskByID(id)
}

func (tu *TaskUsecase) CreateTask(task domain.Task) (domain.Task, error) {
	return tu.taskRepo.CreateTask(task)
}

func (tu *TaskUsecase) UpdateTask(id string, task domain.Task) (domain.Task, error) {
	return tu.taskRepo.UpdateTask(id, task)
}

func (tu *TaskUsecase) DeleteTask(id string) error {
	return tu.taskRepo.DeleteTask(id)
}