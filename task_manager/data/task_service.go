package data

import (
	"sync"
	"task_manager/models"
)

var (
	tasks  = make(map[int]models.Task)
	nextID = 1
	mu     sync.Mutex
)

func GetAllTasks() []models.Task {
	mu.Lock()
	defer mu.Unlock()
	result := make([]models.Task, 0, len(tasks))
	for _, task := range tasks {
		result = append(result, task)
	}
	return result
}

func GetTaskByID(id int) (models.Task, bool) {
	mu.Lock()
	defer mu.Unlock()
	task, found := tasks[id]
	return task, found
}

func CreateTask(task models.Task) models.Task {
	mu.Lock()
	defer mu.Unlock()
	task.ID = nextID
	tasks[nextID] = task
	nextID++
	return task
}

func UpdateTask(id int, updated models.Task) (models.Task, bool) {
	mu.Lock()
	defer mu.Unlock()
	_, found := tasks[id]
	if !found {
		return models.Task{}, false
	}
	updated.ID = id
	tasks[id] = updated
	return updated, true
}

func DeleteTask(id int) bool {
	mu.Lock()
	defer mu.Unlock()
	_, found := tasks[id]
	if !found {
		return false
	}
	delete(tasks, id)
	return true
}
