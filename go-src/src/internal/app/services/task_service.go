package services

import (
	"errors"
	"fmt"
	"nro-go/internal/core/domain"
	"nro-go/internal/core/ports"
	"sync"
)

// TaskService quản lý task/quest system.
//
// EXPLANATION:
// TaskService chịu trách nhiệm:
// 1. Load và quản lý hệ thống nhiệm vụ (TaskTemplate).
// 2. Xử lý logic nhận nhiệm vụ (Accept), cập nhật tiến độ (Update Progress), và trả nhiệm vụ (Complete).
// 3. Trao thưởng (Rewards) khi hoàn thành nhiệm vụ.
type TaskService struct {
	repo  ports.TaskRepository
	tasks map[int]*domain.Task
	mu    sync.RWMutex
}

var taskServiceInstance *TaskService
var taskOnce sync.Once

// GetTaskService returns singleton instance.
func GetTaskService() *TaskService {
	taskOnce.Do(func() {
		taskServiceInstance = &TaskService{
			tasks: make(map[int]*domain.Task),
		}
	})
	return taskServiceInstance
}

func (s *TaskService) SetRepository(repo ports.TaskRepository) {
	s.repo = repo
}

// LoadTasks loads all task templates from the repository.
func (s *TaskService) LoadTasks() error {
	if s.repo == nil {
		return fmt.Errorf("task repository not set")
	}

	tasks, err := s.repo.GetTasks()
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, t := range tasks {
		s.tasks[t.ID] = t
	}
	fmt.Printf("[TASK] Loaded %d task templates from DB\n", len(tasks))
	return nil
}

// GetTask returns task template by ID.
func (s *TaskService) GetTask(id int) *domain.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tasks[id]
}

// AcceptTask adds task to player's task list.
func (s *TaskService) AcceptTask(player *domain.Player, taskID int) error {
	task := s.GetTask(taskID)
	if task == nil {
		return errors.New("task not found")
	}

	// Check level requirement
	if player.Level < task.RequireLevel {
		return fmt.Errorf("require level %d", task.RequireLevel)
	}

	// Check if already has this task
	for _, pt := range player.Tasks {
		if pt.TaskID == taskID {
			return errors.New("already have this task")
		}
	}

	// Add task
	playerTask := &domain.PlayerTask{
		TaskID:   taskID,
		Progress: make(map[int]int),
		Status:   0, // In Progress
	}

	player.Tasks = append(player.Tasks, playerTask)
	fmt.Printf("[TASK] %s accepted task: %s\n", player.Name, task.Name)
	return nil
}

// UpdateProgress updates task progress.
func (s *TaskService) UpdateProgress(player *domain.Player, taskID int, objectiveIndex int, count int) error {
	// Find player task
	var playerTask *domain.PlayerTask
	for _, pt := range player.Tasks {
		if pt.TaskID == taskID {
			playerTask = pt
			break
		}
	}

	if playerTask == nil {
		return errors.New("task not found in player's task list")
	}

	// Update progress
	playerTask.Progress[objectiveIndex] += count

	// Check if task is complete
	task := s.GetTask(taskID)
	if task != nil && s.checkTaskComplete(playerTask, task) {
		playerTask.Status = 1 // Completed
		fmt.Printf("[TASK] %s completed task: %s\n", player.Name, task.Name)
	}

	return nil
}

// checkTaskComplete checks if all objectives are done.
func (s *TaskService) checkTaskComplete(playerTask *domain.PlayerTask, task *domain.Task) bool {
	for i, obj := range task.Objectives {
		if playerTask.Progress[i] < obj.Count {
			return false
		}
	}
	return true
}

// ClaimReward gives rewards to player.
func (s *TaskService) ClaimReward(player *domain.Player, taskID int) error {
	// Find player task
	var playerTask *domain.PlayerTask
	for _, pt := range player.Tasks {
		if pt.TaskID == taskID {
			playerTask = pt
			break
		}
	}

	if playerTask == nil {
		return errors.New("task not found")
	}

	if playerTask.Status != 1 {
		return errors.New("task not completed")
	}

	// Get rewards
	task := s.GetTask(taskID)
	if task == nil {
		return errors.New("task template not found")
	}

	// Give rewards
	player.Exp += task.Rewards.Exp
	// TODO: Add gold to player
	// TODO: Add items to inventory

	playerTask.Status = 2 // Claimed
	fmt.Printf("[TASK] %s claimed rewards for: %s (Exp: %d)\n",
		player.Name, task.Name, task.Rewards.Exp)

	return nil
}
