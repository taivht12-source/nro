package network

import (
	"fmt"
	"nro/src/internal/app/commands"
	"nro/src/internal/app/services"
	"nro/src/internal/core/domain"
	"nro/src/internal/infrastructure/session"
	"nro/src/pkg/protocol"
)

// Controller xử lý các gói tin từ Client.
type Controller struct {
	loginHandler *commands.LoginHandler
	charService  *services.CharacterService
}

func NewController(loginHandler *commands.LoginHandler, charService *services.CharacterService) *Controller {
	return &Controller{
		loginHandler: loginHandler,
		charService:  charService,
	}
}

// OnMessage được gọi khi Session nhận được một Message.
func (c *Controller) OnMessage(sess *protocol.Session, msg *protocol.Message) {
	defer msg.Cleanup()

	switch protocol.Cmd(msg.Command) {
	case protocol.CMD_LOGIN:
		c.handleLogin(sess, msg)
	case protocol.CMD_PLAYER_MOVE:
		c.handlePlayerMove(sess, msg)
	case protocol.CMD_USE_SKILL:
		c.handleUseSkill(sess, msg)
	case protocol.CMD_NPC_INTERACT:
		c.handleNPCInteract(sess, msg)
	case protocol.CMD_TASK_ACCEPT:
		c.handleTaskAccept(sess, msg)
	case protocol.CMD_TASK_COMPLETE:
		c.handleTaskComplete(sess, msg)
	case protocol.CMD_SHOP_OPEN:
		c.handleShopOpen(sess, msg)
	case protocol.CMD_SHOP_BUY:
		c.handleShopBuy(sess, msg)
	case protocol.CMD_SKILL_SELECT:
		c.handleSkillSelect(sess, msg)
	default:
		fmt.Printf("Unhandled CMD: %d (%s)\n", msg.Command, protocol.Cmd(msg.Command))
	}
}

// handleLogin xử lý yêu cầu đăng nhập.
func (c *Controller) handleLogin(sess *protocol.Session, msg *protocol.Message) {
	req := &protocol.LoginRequest{}
	if err := req.Decode(msg); err != nil {
		fmt.Printf("[LOGIN] Error decoding request: %v\n", err)
		c.sendLoginError(sess, "Invalid login data")
		return
	}

	username := req.Username
	password := req.Password
	fmt.Printf("[LOGIN] Attempt: username=%s\n", username)

	// Kiểm tra nếu không có loginHandler (NO-DB mode)
	if c.loginHandler == nil {
		fmt.Println("[LOGIN] NO-DB mode - using mock login")
		c.handleMockLogin(sess, username)
		return
	}

	// Xác thực với LoginHandler
	user, err := c.loginHandler.Handle(commands.LoginCommand{
		Username: username,
		Password: password,
	})

	if err != nil {
		fmt.Printf("[LOGIN] Failed for %s: %v\n", username, err)
		c.sendLoginError(sess, err.Error())
		return
	}

	fmt.Printf("[LOGIN] ✓ User %s authenticated\n", user.Username)

	// Lấy danh sách nhân vật
	characters, err := c.charService.GetCharacterList(user.ID)
	if err != nil {
		fmt.Printf("[LOGIN] Error getting characters: %v\n", err)
		c.sendLoginError(sess, "Failed to load characters")
		return
	}

	// Gửi danh sách nhân vật cho client
	c.sendCharacterList(sess, characters)

	// Lưu user vào session (để sau này chọn nhân vật)
	sess.UserID = user.ID
}

// handleMockLogin xử lý đăng nhập giả (cho NO-DB mode).
func (c *Controller) handleMockLogin(sess *protocol.Session, username string) {
	// Tạo mock character list
	mockChars := []*domain.Player{
		{
			ID:     1,
			UserID: 1,
			Name:   "Goku",
			Gender: 0,
			Power:  1000,
		},
		{
			ID:     2,
			UserID: 1,
			Name:   "Vegeta",
			Gender: 0,
			Power:  900,
		},
	}

	c.sendCharacterList(sess, mockChars)
	sess.UserID = 1
	fmt.Printf("[LOGIN] ✓ Mock login successful for %s\n", username)
}

// sendCharacterList gửi danh sách nhân vật cho client.
func (c *Controller) sendCharacterList(sess *protocol.Session, characters []*domain.Player) {
	response := protocol.NewMessage(-1) // LOGIN response
	response.WriteByte(0)               // Success code

	// Ghi số lượng nhân vật
	response.WriteByte(int8(len(characters)))

	// Ghi thông tin từng nhân vật
	for _, char := range characters {
		response.WriteInt(int32(char.ID))
		response.WriteUTF(char.Name)
		response.WriteByte(char.Gender)
		response.WriteShort(char.Head)
		response.WriteShort(char.Body)
		response.WriteShort(char.Leg)
		response.WriteInt(int32(char.Power >> 32)) // High 32 bits
		response.WriteInt(int32(char.Power))       // Low 32 bits
	}

	sess.SendMessage(response)
	fmt.Printf("[LOGIN] Sent %d characters to client\n", len(characters))
}

// sendLoginError gửi thông báo lỗi đăng nhập.
func (c *Controller) sendLoginError(sess *protocol.Session, errorMsg string) {
	response := protocol.NewMessage(-1)
	response.WriteByte(1) // Error code
	response.WriteUTF(errorMsg)
	sess.SendMessage(response)
}

// handleUseSkill xử lý sử dụng skill.
func (c *Controller) handleUseSkill(sess *protocol.Session, msg *protocol.Message) {
	req := &protocol.UseSkillRequest{}
	if err := req.Decode(msg); err != nil {
		fmt.Printf("[SKILL] Error decoding request: %v\n", err)
		return
	}

	skillIndex := int(req.SkillIndex)
	targetID := int(req.TargetID)

	// Get player from session
	if sess.Player == nil {
		fmt.Println("[SKILL] No player in session")
		return
	}

	player, ok := sess.Player.(*domain.Player)
	if !ok {
		fmt.Println("[SKILL] Invalid player type")
		return
	}

	// Use skill via CombatService
	combatService := services.GetCombatService()
	damage, err := combatService.UseSkill(player, int(skillIndex), int(targetID))

	if err != nil {
		fmt.Printf("[SKILL] %s failed to use skill: %v\n", player.Name, err)
		c.sendSkillError(sess, err.Error())
		return
	}

	// Broadcast skill effect to zone
	c.broadcastSkillEffect(player, int(skillIndex), int(targetID), damage)

	// Send updated HP/MP to player
	c.sendPlayerStats(sess, player)
}

// broadcastSkillEffect broadcasts skill usage to players in zone.
func (c *Controller) broadcastSkillEffect(player *domain.Player, skillIndex int, targetID int, damage int) {
	if skillIndex < 0 || skillIndex >= len(player.Skills) {
		return
	}

	skill := player.Skills[skillIndex]
	zoneService := services.GetZoneService()

	// Create skill effect message
	msg := protocol.NewMessage(-14) // SKILL_EFFECT
	msg.WriteInt(int32(player.ID))
	msg.WriteByte(int8(skillIndex))
	msg.WriteInt(int32(skill.Template.ID))
	msg.WriteInt(int32(targetID))
	msg.WriteInt(int32(damage))

	// Broadcast to zone
	zoneService.BroadcastToZone(player.MapID, player.ZoneID, msg)
	msg.Cleanup()

	fmt.Printf("[SKILL] Broadcasted %s's %s to zone (damage: %d)\n",
		player.Name, skill.Template.Name, damage)
}

// sendPlayerStats sends HP/MP update to player.
func (c *Controller) sendPlayerStats(sess *protocol.Session, player *domain.Player) {
	msg := protocol.NewMessage(-15) // PLAYER_STATS
	msg.WriteInt(int32(player.HP))
	msg.WriteInt(int32(player.MaxHP))
	msg.WriteInt(int32(player.MP))
	msg.WriteInt(int32(player.MaxMP))
	sess.SendMessage(msg)
}

// sendSkillError sends skill error message.
func (c *Controller) sendSkillError(sess *protocol.Session, errorMsg string) {
	msg := protocol.NewMessage(-16) // SKILL_ERROR
	msg.WriteUTF(errorMsg)
	sess.SendMessage(msg)
}

func (c *Controller) handlePlayerMove(sess *protocol.Session, msg *protocol.Message) {
	// Đọc dữ liệu di chuyển
	// Format: [1 byte status] [2 byte X] [2 byte Y] (Thường là vậy, cần check kỹ code Java)
	// Trong Java: byte b = msg.reader().readByte(); short x = msg.reader().readShort(); short y = msg.reader().readShort();

	// Lưu ý: msg.Reader đã được khởi tạo trong Session.readLoop khi tạo Message

	status, err := msg.ReadByte()
	if err != nil {
		return
	}
	x, err := msg.ReadShort()
	if err != nil {
		return
	}
	y, err := msg.ReadShort()
	if err != nil {
		return
	}

	fmt.Printf("[MOVE] Session %d: Status=%d, X=%d, Y=%d\n", sess.ID, status, x, y)

	var player *domain.Player

	// Cập nhật vị trí vào Player Entity
	if sess.Player != nil {
		if p, ok := sess.Player.(*domain.Player); ok {
			player = p

			// Validate movement
			validator := services.GetMovementValidator()

			// 1. Kiểm tra vị trí có hợp lệ không
			if !validator.ValidatePosition(x, y, player.MapID) {
				fmt.Printf("[MOVE] REJECTED - Invalid position for Player %s: (%d, %d) on Map %d\n",
					player.Name, x, y, player.MapID)
				return
			}

			// 2. Kiểm tra tốc độ di chuyển (anti-cheat)
			if !validator.ValidateSpeed(player, x, y) {
				fmt.Printf("[MOVE] REJECTED - Speed hack detected for Player %s: from (%d,%d) to (%d,%d)\n",
					player.Name, player.X, player.Y, x, y)
				return
			}

			// 3. Cập nhật vị trí
			validator.UpdateLastMove(player, x, y)
			fmt.Printf("[MOVE] ✓ Player %s moved to (%d, %d)\n", player.Name, x, y)

			// 4. Kiểm tra waypoint (chuyển map)
			zoneService := services.GetZoneService()
			waypoint := zoneService.CheckWaypoint(player)
			if waypoint != nil {
				// Teleport player to new map
				zoneService.ChangeMap(player, waypoint)

				// Send new map info
				c.sendMapInfo(sess, int(waypoint.GoMapID))

				// Send player list in new zone
				c.sendPlayerList(sess, player)

				fmt.Printf("[WAYPOINT] Player %s teleported to map %d\n", player.Name, waypoint.GoMapID)
				return // Don't broadcast move if teleported
			}
		}
	} else {
		// Mock login nếu chưa có player (chỉ để test Phase 3)
		mockPlayer := &domain.Player{
			ID:     1,
			Name:   "Goku Test",
			X:      x,
			Y:      y,
			MapID:  0, // Map mặc định
			ZoneID: 0, // Zone mặc định
		}
		sess.Player = mockPlayer
		player = mockPlayer

		// Đăng ký vào SessionManager và ZoneService (Mock Login)
		session.GetManager().Add(mockPlayer.ID, sess)
		services.GetZoneService().EnterZone(mockPlayer, mockPlayer.MapID, mockPlayer.ZoneID)

		fmt.Printf("[MOVE] ✓ Mocked Player %s login and moved to (%d, %d)\n", mockPlayer.Name, x, y)
	}

	if player != nil {
		// Broadcast cho các người chơi khác trong cùng Zone
		// Lưu ý: msg lúc này đã được đọc hết (ReadIndex ở cuối),
		// nếu muốn gửi lại nguyên văn gói tin này thì cần reset reader hoặc tạo gói tin mới.
		// Tuy nhiên, logic chuẩn là Server nhận tin -> Xử lý -> Tạo gói tin phản hồi mới.
		// Nhưng với gói tin di chuyển, thường Server chỉ forward (chuyển tiếp) cho người khác.

		// Cách đơn giản nhất: Tạo message mới để broadcast
		msgMove := protocol.NewMessage(-7)
		msgMove.WriteByte(int8(player.ID)) // Gửi ID người di chuyển (giả sử protocol cần vậy)
		msgMove.WriteShort(x)
		msgMove.WriteShort(y)

		services.GetZoneService().BroadcastMove(player, msgMove)
		msgMove.Cleanup()
	}
}

// sendMapInfo gửi thông tin map cho client.
func (c *Controller) sendMapInfo(sess *protocol.Session, mapID int) {
	mapService := services.GetMapService()
	mapTemplate := mapService.GetMap(mapID)

	if mapTemplate == nil {
		fmt.Printf("[MAP] Error: Map %d not found\n", mapID)
		return
	}

	msg := protocol.NewMessage(-24) // MAP_INFO
	msg.WriteInt(int32(mapTemplate.ID))
	msg.WriteUTF(mapTemplate.Name)
	msg.WriteByte(mapTemplate.Type)
	msg.WriteByte(mapTemplate.PlanetID)
	msg.WriteByte(mapTemplate.TileID)
	msg.WriteByte(mapTemplate.BgID)
	msg.WriteByte(mapTemplate.BgType)
	msg.WriteByte(int8(mapTemplate.Zones))

	sess.SendMessage(msg)
	fmt.Printf("[MAP] Sent map info for %s (ID: %d)\n", mapTemplate.Name, mapTemplate.ID)
}

// sendPlayerList gửi danh sách người chơi trong zone.
func (c *Controller) sendPlayerList(sess *protocol.Session, currentPlayer *domain.Player) {
	zoneService := services.GetZoneService()
	players := zoneService.GetPlayersInZone(currentPlayer.MapID, currentPlayer.ZoneID)

	msg := protocol.NewMessage(-6) // PLAYER_LIST
	msg.WriteByte(int8(len(players)))

	for _, p := range players {
		msg.WriteInt(int32(p.ID))
		msg.WriteUTF(p.Name)
		msg.WriteByte(p.Gender)
		msg.WriteShort(p.X)
		msg.WriteShort(p.Y)
		msg.WriteShort(p.Head)
		msg.WriteShort(p.Body)
		msg.WriteShort(p.Leg)
	}

	sess.SendMessage(msg)
	fmt.Printf("[ZONE] Sent player list (%d players) to %s\n", len(players), currentPlayer.Name)
}

// handleShopOpen xử lý yêu cầu mở shop.
func (c *Controller) handleShopOpen(sess *protocol.Session, msg *protocol.Message) {
	// Read Type (0: Normal, 1: Learn Skill, 2: ...), NPC ID, Shop ID?
	// NRO protocol: type(byte), npcID(byte/int?)
	// Let's assume: type(byte), npcID(byte)

	opType, _ := msg.ReadByte()
	npcID, _ := msg.ReadByte() // or ReadInt? usually byte for NPC on map

	fmt.Printf("[SHOP] Request Open: Type=%d, NPC=%d\n", opType, npcID)

	// Get Shop by NPC
	shopService := services.GetShopService()
	shop, err := shopService.GetShopByNPC(int(npcID))
	if err != nil {
		fmt.Printf("[SHOP] Error getting shop for NPC %d: %v\n", npcID, err)
		return
	}

	// Send Shop Open Packet (-44)
	c.sendShopOpen(sess, shop)
}

func (c *Controller) sendShopOpen(sess *protocol.Session, shop *domain.Shop) {
	msg := protocol.NewMessage(-44)
	msg.WriteByte(shop.Type)
	msg.WriteByte(int8(len(shop.TabName))) // Number of tabs

	// Write Tabs
	for _, tab := range shop.TabName {
		msg.WriteUTF(tab)

		// Filter items for this tab (simplified: all items in first tab for now or split logic)
		// NRO usually sends all items and client filters? Or sends items per tab?
		// Let's assume simple structure: 1 tab, all items.
		// If multiple tabs, we need logic to assign items to tabs.
		// For now, write all items in first tab, empty for others.

		// Simplified: Just write items count and items
		msg.WriteByte(int8(len(shop.Items)))
		for _, item := range shop.Items {
			msg.WriteInt(int32(item.ID))
			msg.WriteInt(int32(item.TemplateID))
			msg.WriteUTF(item.Template.Name)
			msg.WriteByte(item.Template.Type) // Item Type
			msg.WriteInt(int32(item.Price))
			msg.WriteInt(int32(item.Template.IconID))
			msg.WriteByte(0) // isNew?
			// More item info if needed
		}
	}
	sess.SendMessage(msg)
}

// handleShopBuy xử lý mua vật phẩm.
func (c *Controller) handleShopBuy(sess *protocol.Session, msg *protocol.Message) {
	// Read: ShopID(byte/int), ItemID(int), Quantity(int)
	// NRO: type(byte), shopID(byte), itemID(short/int), quantity(int)

	_, _ = msg.ReadByte() // type?
	shopID, _ := msg.ReadByte()
	itemID, _ := msg.ReadInt() // or Short?
	quantity, _ := msg.ReadInt()

	if quantity <= 0 {
		return
	}

	player := sess.Player.(*domain.Player)
	shopService := services.GetShopService()

	err := shopService.BuyItem(player, int(shopID), int(itemID), int(quantity))
	if err != nil {
		fmt.Printf("[SHOP] Buy failed: %v\n", err)
		c.sendShopError(sess, err.Error())
		return
	}

	// Send success / Update Inventory
	c.sendPlayerStats(sess, player) // Update Gold/Gem
	// Send Inventory Update Packet (Need to implement)
	// c.sendInventoryUpdate(sess, player)

	fmt.Printf("[SHOP] Buy success for %s\n", player.Name)
}

func (c *Controller) sendShopError(sess *protocol.Session, message string) {
	// Send alert
	msg := protocol.NewMessage(-16) // Server Alert
	msg.WriteUTF(message)
	sess.SendMessage(msg)
}

// handleSkillSelect xử lý chọn skill.
func (c *Controller) handleSkillSelect(sess *protocol.Session, msg *protocol.Message) {
	skillID, _ := msg.ReadShort() // Or Int?
	// Update player's selected skill (if we track it)
	// Usually client sends UseSkill with specific skill index, so Select might just be visual or for auto-attack.
	fmt.Printf("[SKILL] Selected Skill ID: %d\n", skillID)
}
