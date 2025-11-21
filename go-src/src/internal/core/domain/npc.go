package domain

// NPCTemplate mẫu NPC.
type NPCTemplate struct {
	ID       int
	Name     string
	MapID    int
	X, Y     int16
	Avatar   int
	Dialogue []string
	Tasks    []int // Task IDs this NPC offers
}

// Task nhiệm vụ.
type Task struct {
	ID           int
	Name         string
	Description  string
	NPCID        int
	RequireLevel int
	Rewards      *TaskReward
	Objectives   []*TaskObjective
}

// TaskObjective mục tiêu nhiệm vụ.
type TaskObjective struct {
	Type   int // 0: Kill mob, 1: Collect item, 2: Talk to NPC
	Target int // Mob ID, Item ID, or NPC ID
	Count  int
}

// TaskReward phần thưởng nhiệm vụ.
type TaskReward struct {
	Exp   int64
	Gold  int
	Items []*Item
}

// PlayerTask nhiệm vụ của người chơi.
type PlayerTask struct {
	TaskID   int
	Progress map[int]int // ObjectiveIndex -> Count
	Status   int         // 0: In Progress, 1: Completed, 2: Claimed
}
