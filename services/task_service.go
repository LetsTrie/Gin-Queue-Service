package services

type Task struct {
	ID          int       `json:"id"`
	IsCompleted bool      `json:"is_completed"`
	Name        string    `json:"name"`
	StoryPoint  int       `json:"story_point"`
}

var tasks []Task
var idCounter int

func GetTasks() []Task {
	return tasks
}

func GetTask(id int) *Task {
	for _, task := range tasks {
		if task.ID == id {
			return &task
		}
	}
	return nil
}

func CreateTask(task Task) {
	task.ID = idCounter
	tasks = append(tasks, task)
	idCounter++
}