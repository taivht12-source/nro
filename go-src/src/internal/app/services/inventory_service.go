package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"nro-go/internal/core/domain"
	"nro-go/internal/core/ports"
	"sync"
)

// InventoryService quản lý inventory của player.
//
// EXPLANATION:
// InventoryService chịu trách nhiệm:
// 1. Quản lý túi đồ (Bag), rương (Box), và trang bị trên người (Body/Equipment).
// 2. Xử lý logic thêm/xóa vật phẩm, sắp xếp túi đồ.
// 3. Xử lý logic sử dụng vật phẩm (Use Item) và trang bị/tháo trang bị (Equip/Unequip).
type InventoryService struct {
	itemRepo ports.ItemRepository
}

var inventoryServiceInstance *InventoryService
var inventoryOnce sync.Once

// GetInventoryService returns inventory service instance.
func GetInventoryService() *InventoryService {
	inventoryOnce.Do(func() {
		inventoryServiceInstance = &InventoryService{}
	})
	return inventoryServiceInstance
}

func (inv *InventoryService) SetItemRepository(repo ports.ItemRepository) {
	inv.itemRepo = repo
}

// AddItem thêm item vào hành trang (Bag).
func (inv *InventoryService) AddItem(player *domain.Player, item *domain.Item) error {
	if player.Inventory == nil {
		return errors.New("player inventory is nil")
	}

	// Check if bag is full (hardcoded limit for now, e.g., 20 slots)
	if len(player.Inventory.ItemsBag) >= 20 {
		return errors.New("inventory bag is full")
	}

	// Check if item is stackable (consumables)
	if item.Template.Type == 3 || item.Template.Type == 4 { // Consumables or Dragon Balls
		// Find existing stack
		for _, existingItem := range player.Inventory.ItemsBag {
			if existingItem.Template.ID == item.Template.ID {
				existingItem.Quantity += item.Quantity
				fmt.Printf("[INV] Stacked %s x%d for %s\n", item.Template.Name, item.Quantity, player.Name)
				return nil
			}
		}
	}

	// Add new item
	item.ID = len(player.Inventory.ItemsBag) // Simple ID assignment
	player.Inventory.ItemsBag = append(player.Inventory.ItemsBag, item)
	fmt.Printf("[INV] Added %s to %s's bag\n", item.Template.Name, player.Name)
	return nil
}

// RemoveItem xóa item khỏi hành trang (Bag).
func (inv *InventoryService) RemoveItem(player *domain.Player, index int) error {
	if player.Inventory == nil {
		return errors.New("player inventory is nil")
	}

	if index < 0 || index >= len(player.Inventory.ItemsBag) {
		return errors.New("invalid item index")
	}

	item := player.Inventory.ItemsBag[index]

	// If stackable and quantity > 1, just decrease quantity
	if item.Quantity > 1 {
		item.Quantity--
		fmt.Printf("[INV] Decreased %s quantity to %d for %s\n", item.Template.Name, item.Quantity, player.Name)
		return nil
	}

	// Remove item completely
	player.Inventory.ItemsBag = append(player.Inventory.ItemsBag[:index], player.Inventory.ItemsBag[index+1:]...)
	fmt.Printf("[INV] Removed %s from %s's bag\n", item.Template.Name, player.Name)
	return nil
}

// EquipItem trang bị item (Move from Bag to Body).
func (inv *InventoryService) EquipItem(player *domain.Player, index int) error {
	if player.Inventory == nil {
		return errors.New("player inventory is nil")
	}

	if index < 0 || index >= len(player.Inventory.ItemsBag) {
		return errors.New("invalid item index")
	}

	item := player.Inventory.ItemsBag[index]

	// Check if item can be equipped (based on type/part)
	// Type: 0=Ao, 1=Quan, 2=Gang, 3=Giay, 4=Rada, 5=Pet, 6=Mount, 7=PhuKien
	// This mapping depends on NRO logic. Assuming standard types.
	// Let's assume Part indicates the slot index in ItemsBody.
	// Or we use Type to determine slot.
	// Standard NRO: 0:Ao, 1:Quan, 2:Gang, 3:Giay, 4:Rada, 5:Pet...

	slot := -1
	switch item.Template.Type {
	case 0:
		slot = 0 // Áo
	case 1:
		slot = 1 // Quần
	case 2:
		slot = 2 // Găng
	case 3:
		slot = 3 // Giày
	case 4:
		slot = 4 // Rada
	default:
		return errors.New("item cannot be equipped")
	}

	// Ensure Body has enough slots
	if len(player.Inventory.ItemsBody) <= slot {
		// Expand Body slots if needed (should be initialized with enough nil slots)
		for len(player.Inventory.ItemsBody) <= slot {
			player.Inventory.ItemsBody = append(player.Inventory.ItemsBody, nil)
		}
	}

	// Swap items
	currentEquip := player.Inventory.ItemsBody[slot]

	// Put current equip back to bag (if any)
	if currentEquip != nil {
		player.Inventory.ItemsBag = append(player.Inventory.ItemsBag, currentEquip)
	}

	// Equip new item
	player.Inventory.ItemsBody[slot] = item

	// Remove from bag
	player.Inventory.ItemsBag = append(player.Inventory.ItemsBag[:index], player.Inventory.ItemsBag[index+1:]...)

	fmt.Printf("[INV] Equipped %s to slot %d for %s\n", item.Template.Name, slot, player.Name)
	return nil
}

// UnequipItem gỡ trang bị (Move from Body to Bag).
func (inv *InventoryService) UnequipItem(player *domain.Player, slot int) error {
	if player.Inventory == nil {
		return errors.New("player inventory is nil")
	}

	if slot < 0 || slot >= len(player.Inventory.ItemsBody) {
		return errors.New("invalid equipment slot")
	}

	item := player.Inventory.ItemsBody[slot]
	if item == nil {
		return errors.New("no item equipped in slot")
	}

	// Add back to bag
	if err := inv.AddItem(player, item); err != nil {
		return err
	}

	// Remove from body
	player.Inventory.ItemsBody[slot] = nil
	fmt.Printf("[INV] Unequipped %s from slot %d for %s\n", item.Template.Name, slot, player.Name)
	return nil
}

// FindItem tìm item trong hành trang theo template ID.
func (inv *InventoryService) FindItem(player *domain.Player, templateID int) *domain.Item {
	if player.Inventory == nil {
		return nil
	}

	for _, item := range player.Inventory.ItemsBag {
		if item.Template.ID == templateID {
			return item
		}
	}
	return nil
}

// ParseItems parses JSON string from DB to []*Item.
// Format: ["[\"TemplateID,Quantity,OptionsJSON,Timestamp\"]", ...]
func (inv *InventoryService) ParseItems(data string) []*domain.Item {
	if inv.itemRepo == nil {
		fmt.Println("[INV] Warning: ItemRepository not set, cannot parse items")
		return nil
	}

	var rawList []string
	if err := json.Unmarshal([]byte(data), &rawList); err != nil {
		fmt.Println("[INV] Error parsing item JSON list:", err)
		return nil
	}

	var items []*domain.Item
	for _, s := range rawList {
		// s is like "[TemplateID,Quantity,OptionsJSON,Timestamp]"
		// Need to parse this inner JSON array
		var itemData []interface{}
		if err := json.Unmarshal([]byte(s), &itemData); err != nil {
			continue
		}

		if len(itemData) < 2 {
			continue
		}

		// 1. TemplateID
		templateID := int(itemData[0].(float64))

		// 2. Quantity
		quantity := int(itemData[1].(float64))

		// 3. Options (String of JSON array)
		var options []*domain.ItemOption
		if len(itemData) > 2 {
			optStr, ok := itemData[2].(string)
			if ok {
				// Parse options: "[[OptionID,Param], ...]"
				var rawOpts [][]int
				if err := json.Unmarshal([]byte(optStr), &rawOpts); err == nil {
					for _, ro := range rawOpts {
						if len(ro) >= 2 {
							options = append(options, &domain.ItemOption{
								ID:    ro[0],
								Param: ro[1],
							})
						}
					}
				}
			}
		}

		// Load Template
		template, err := inv.itemRepo.GetTemplate(templateID)
		if err != nil {
			fmt.Printf("[INV] Warning: Template %d not found\n", templateID)
			// Create dummy template to avoid crash
			template = &domain.ItemTemplate{ID: templateID, Name: fmt.Sprintf("Unknown Item %d", templateID)}
		}

		items = append(items, &domain.Item{
			TemplateID: templateID,
			Template:   template,
			Quantity:   quantity,
			Options:    options,
		})
	}
	return items
}

// UseConsumable sử dụng vật phẩm tiêu hao.
func (inv *InventoryService) UseConsumable(player *domain.Player, index int) error {
	if player.Inventory == nil {
		return errors.New("player inventory is nil")
	}

	if index < 0 || index >= len(player.Inventory.ItemsBag) {
		return errors.New("invalid item index")
	}

	item := player.Inventory.ItemsBag[index]

	// Check if item is consumable
	// Assuming Type 3 is consumable (Food/Potions)
	if item.Template.Type != 3 {
		return errors.New("item is not consumable")
	}

	// Apply item effect based on template ID
	// TODO: Move this logic to ItemEffectService or similar
	switch item.Template.ID {
	case 30: // Đậu Thần - Full HP/MP
		player.HP = player.MaxHP
		player.MP = player.MaxMP
		fmt.Printf("[INV] %s used Đậu Thần, HP/MP fully restored\n", player.Name)

	case 31: // Bí Ngô - +100 HP
		player.HP += 100
		if player.HP > player.MaxHP {
			player.HP = player.MaxHP
		}
		fmt.Printf("[INV] %s used Bí Ngô, HP +100\n", player.Name)

	case 32: // Cà Rốt - +50 MP
		player.MP += 50
		if player.MP > player.MaxMP {
			player.MP = player.MaxMP
		}
		fmt.Printf("[INV] %s used Cà Rốt, MP +50\n", player.Name)
	}

	// Remove item
	return inv.RemoveItem(player, index)
}
