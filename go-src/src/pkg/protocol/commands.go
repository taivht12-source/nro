package protocol

// Cmd defines the command ID for network messages.
type Cmd int8

const (
	CMD_LOGIN         Cmd = -1
	CMD_REGISTER      Cmd = -2
	CMD_PLAYER_MOVE   Cmd = -7
	CMD_USE_SKILL     Cmd = -11
	CMD_SKILL_EFFECT  Cmd = -14 // Server -> Client
	CMD_PLAYER_STATS  Cmd = -15 // Server -> Client
	CMD_SKILL_ERROR   Cmd = -16 // Server -> Client
	CMD_NPC_INTERACT  Cmd = -20
	CMD_TASK_ACCEPT   Cmd = -21
	CMD_TASK_COMPLETE Cmd = -22
	CMD_TASK_LIST     Cmd = -23 // Server -> Client
	CMD_SKILL_SELECT  Cmd = -25
	CMD_NPC_DIALOGUE  Cmd = -26 // Server -> Client? Or -20 response?
	CMD_SHOP_OPEN     Cmd = -44
	CMD_SHOP_BUY      Cmd = -45
	CMD_OPEN_UI       Cmd = 38 // Server -> Client (Menu)
)

// String returns the string representation of the command.
func (c Cmd) String() string {
	switch c {
	case CMD_LOGIN:
		return "LOGIN"
	case CMD_REGISTER:
		return "REGISTER"
	case CMD_PLAYER_MOVE:
		return "PLAYER_MOVE"
	case CMD_USE_SKILL:
		return "USE_SKILL"
	case CMD_NPC_INTERACT:
		return "NPC_INTERACT"
	case CMD_TASK_ACCEPT:
		return "TASK_ACCEPT"
	case CMD_TASK_COMPLETE:
		return "TASK_COMPLETE"
	case CMD_SKILL_SELECT:
		return "SKILL_SELECT"
	case CMD_SHOP_OPEN:
		return "SHOP_OPEN"
	case CMD_SHOP_BUY:
		return "SHOP_BUY"
	default:
		return "UNKNOWN"
	}
}
