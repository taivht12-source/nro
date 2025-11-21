package network

import (
	"fmt"
	"nro-go/internal/app/services"
	"nro-go/internal/core/domain"
	"nro-go/pkg/protocol"
)

// handleNPCInteract xử lý tương tác với NPC.
func (c *Controller) handleNPCInteract(sess *protocol.Session, msg *protocol.Message) {
	req := &protocol.NPCInteractRequest{}
	if err := req.Decode(msg); err != nil {
		fmt.Printf("[NPC] Error decoding request: %v\n", err)
		return
	}

	npcID := int(req.NPCID)

	player := sess.Player.(*domain.Player)
	menuService := services.GetMenuService()

	// Get dynamic menu info
	menuInfo := menuService.GetMenuInfo(player, int(npcID))
	if menuInfo == nil {
		fmt.Printf("[NPC] NPC %d not found or no menu\n", npcID)
		return
	}

	// Send Menu Packet
	c.sendMenu(sess, menuInfo)
}

// sendMenu gửi gói tin hiển thị menu.
func (c *Controller) sendMenu(sess *protocol.Session, info *services.MenuInfo) {
	msg := protocol.NewMessage(38) // OPEN_UI / MENU
	msg.WriteShort(int16(info.NPCID))
	msg.WriteUTF(info.Text)

	msg.WriteByte(int8(len(info.Options)))
	for _, opt := range info.Options {
		msg.WriteUTF(opt) // Menu Text
		// NRO might send ID/Action code for each option too.
		// msg.WriteShort(id)
	}

	sess.SendMessage(msg)
	fmt.Printf("[MENU] Sent menu for NPC %d to %s\n", info.NPCID, sess.Player.(*domain.Player).Name)
}

// handleTaskAccept xử lý nhận nhiệm vụ.
func (c *Controller) handleTaskAccept(sess *protocol.Session, msg *protocol.Message) {
	req := &protocol.TaskAcceptRequest{}
	if err := req.Decode(msg); err != nil {
		fmt.Printf("[TASK] Error decoding request: %v\n", err)
		return
	}

	taskID := int(req.TaskID)

	if sess.Player == nil {
		return
	}

	player, ok := sess.Player.(*domain.Player)
	if !ok {
		return
	}

	taskService := services.GetTaskService()
	err := taskService.AcceptTask(player, int(taskID))

	if err != nil {
		c.sendTaskError(sess, err.Error())
		return
	}

	// Send updated task list
	c.sendTaskList(sess, player)
}

// handleTaskComplete xử lý hoàn thành nhiệm vụ.
func (c *Controller) handleTaskComplete(sess *protocol.Session, msg *protocol.Message) {
	req := &protocol.TaskCompleteRequest{}
	if err := req.Decode(msg); err != nil {
		fmt.Printf("[TASK] Error decoding request: %v\n", err)
		return
	}

	taskID := int(req.TaskID)

	if sess.Player == nil {
		return
	}

	player, ok := sess.Player.(*domain.Player)
	if !ok {
		return
	}

	taskService := services.GetTaskService()

	// Claim reward
	err := taskService.ClaimReward(player, int(taskID))
	if err != nil {
		c.sendTaskError(sess, err.Error())
		return
	}

	// Send updated stats and task list
	c.sendPlayerStats(sess, player)
	c.sendTaskList(sess, player)

	fmt.Printf("[TASK] %s completed task %d\n", player.Name, taskID)
}

// sendTaskList gửi danh sách nhiệm vụ của player.
func (c *Controller) sendTaskList(sess *protocol.Session, player *domain.Player) {
	response := protocol.NewMessage(-23) // TASK_LIST
	response.WriteByte(int8(len(player.Tasks)))

	taskService := services.GetTaskService()
	for _, pt := range player.Tasks {
		task := taskService.GetTask(pt.TaskID)
		if task != nil {
			response.WriteInt(int32(pt.TaskID))
			response.WriteUTF(task.Name)
			response.WriteByte(int8(pt.Status))

			// Send progress for each objective
			response.WriteByte(int8(len(task.Objectives)))
			for i, obj := range task.Objectives {
				progress := pt.Progress[i]
				response.WriteInt(int32(progress))
				response.WriteInt(int32(obj.Count))
			}
		}
	}

	sess.SendMessage(response)
}

// sendTaskError gửi lỗi nhiệm vụ.
func (c *Controller) sendTaskError(sess *protocol.Session, errorMsg string) {
	response := protocol.NewMessage(-25) // TASK_ERROR
	response.WriteUTF(errorMsg)
	sess.SendMessage(response)
}
