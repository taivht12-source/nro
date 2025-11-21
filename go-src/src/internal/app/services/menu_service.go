package services

import (
	"fmt"
	"nro/src/internal/core/domain"
	"sync"
)

// MenuService quản lý hệ thống Menu và Hội thoại (Dialogue).
//
// EXPLANATION (Giải thích):
// MenuService đóng vai trò là "UI Logic" phía Server. Nhiệm vụ của nó là:
//  1. Xây dựng nội dung Menu: Quyết định xem NPC sẽ nói gì và hiển thị những lựa chọn nào cho người chơi.
//     Ví dụ: Nếu người chơi chưa nhận nhiệm vụ -> Hiển thị "Nhận nhiệm vụ".
//     Nếu đã làm xong -> Hiển thị "Trả nhiệm vụ".
//  2. Xử lý lựa chọn (Callback): Khi người chơi chọn một mục (Select), Service này sẽ thực thi hành động tương ứng.
//     Ví dụ: Chọn "Đến Đảo Kame" -> Thực hiện dịch chuyển.
type MenuService struct {
	// menus lưu trữ các menu đang mở hoặc cấu hình menu (nếu cần).
	// Trong mô hình đơn giản, ta có thể xử lý stateless dựa trên NPC ID và Menu ID.
}

var menuServiceInstance *MenuService
var menuOnce sync.Once

func GetMenuService() *MenuService {
	menuOnce.Do(func() {
		menuServiceInstance = &MenuService{}
	})
	return menuServiceInstance
}

// MenuInfo chứa thông tin hiển thị menu
type MenuInfo struct {
	NPCID   int
	Text    string
	Options []string
}

// GetMenuInfo lấy thông tin menu để hiển thị.
func (s *MenuService) GetMenuInfo(player *domain.Player, npcID int) *MenuInfo {
	npcService := GetNPCService()
	npc := npcService.GetTemplate(npcID)
	if npc == nil {
		return nil
	}

	text := npcService.GetDialogue(npcID, 0)
	var options []string

	// Add Task options
	taskService := GetTaskService()
	for _, taskID := range npc.Tasks {
		task := taskService.GetTask(taskID)
		if task != nil {
			// Check status
			// If not accepted: "Nhận nhiệm vụ X"
			// If in progress: "Tiến độ nhiệm vụ X"
			// If completed: "Trả nhiệm vụ X"
			options = append(options, fmt.Sprintf("Nhiệm vụ: %s", task.Name))
		}
	}

	options = append(options, "Đóng")

	return &MenuInfo{
		NPCID:   npcID,
		Text:    text,
		Options: options,
	}
}

// HandleMenuSelect xử lý khi người chơi chọn menu.
func (s *MenuService) HandleMenuSelect(player *domain.Player, npcID int, selectIndex int) {
	fmt.Printf("[MENU] %s selected option %d for NPC %d\n", player.Name, selectIndex, npcID)

	// Logic xử lý:
	// Dựa vào npcID và selectIndex để biết người chơi chọn gì.
	// Đây là phần phức tạp nhất vì cần mapping chính xác.

	npcService := GetNPCService()
	npc := npcService.GetTemplate(npcID)
	if npc == nil {
		return
	}

	// Giả sử options map 1-1 với Tasks
	if selectIndex < len(npc.Tasks) {
		taskID := npc.Tasks[selectIndex]
		taskService := GetTaskService()

		// Thử nhận nhiệm vụ
		err := taskService.AcceptTask(player, taskID)
		if err != nil {
			// Nếu lỗi (đã nhận rồi), thử kiểm tra hoàn thành?
			// Logic này cần chi tiết hơn.
			fmt.Printf("[MENU] Action failed: %v\n", err)
		}
	}
}
