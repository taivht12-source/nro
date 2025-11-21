package services

import (
	"errors"
	"fmt"
	"nro-go/internal/core/domain"
	"nro-go/internal/core/ports"
	"sync"
)

type ShopService struct {
	//
	// EXPLANATION:
	// ShopService chịu trách nhiệm:
	// 1. Quản lý các cửa hàng (Shop) và danh sách vật phẩm bán trong đó.
	// 2. Xử lý logic mua hàng (Buy Item): Kiểm tra tiền (Vàng/Ngọc), trừ tiền, và thêm vật phẩm vào túi.
	// 3. Liên kết Shop với NPC cụ thể.
	repo ports.ShopRepository
	mu   sync.RWMutex
}

var shopServiceInstance *ShopService
var shopOnce sync.Once

func GetShopService() *ShopService {
	shopOnce.Do(func() {
		shopServiceInstance = &ShopService{}
	})
	return shopServiceInstance
}

func (s *ShopService) SetRepository(repo ports.ShopRepository) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.repo = repo
}

// GetShopByNPC returns shop managed by an NPC.
func (s *ShopService) GetShopByNPC(npcID int) (*domain.Shop, error) {
	if s.repo == nil {
		return nil, errors.New("shop repository not set")
	}
	return s.repo.GetShopByNPC(npcID)
}

// BuyItem handles buying an item from a shop.
func (s *ShopService) BuyItem(player *domain.Player, shopID int, itemID int, quantity int) error {
	if s.repo == nil {
		return errors.New("shop repository not set")
	}

	shop, err := s.repo.GetShop(shopID)
	if err != nil {
		return err
	}

	var shopItem *domain.ShopItem
	for _, item := range shop.Items {
		if item.ID == itemID {
			shopItem = item
			break
		}
	}

	if shopItem == nil {
		return errors.New("item not found in shop")
	}

	// Calculate total price
	totalPrice := shopItem.Price * int64(quantity)

	// Check player funds
	inventoryService := GetInventoryService()
	if shopItem.BuyType == 0 { // Gold
		if player.Inventory.Gold < totalPrice {
			return errors.New("not enough gold")
		}
		player.Inventory.Gold -= totalPrice
	} else { // Gem
		if int64(player.Inventory.Gem) < totalPrice {
			return errors.New("not enough gem")
		}
		player.Inventory.Gem -= int(totalPrice)
	}

	// Create item
	itemService := GetItemService()
	newItem := itemService.CreateItem(shopItem.TemplateID, quantity)
	if newItem == nil {
		return errors.New("failed to create item")
	}

	// Add to inventory
	if err := inventoryService.AddItem(player, newItem); err != nil {
		// Refund if failed
		if shopItem.BuyType == 0 {
			player.Inventory.Gold += totalPrice
		} else {
			player.Inventory.Gem += int(totalPrice)
		}
		return err
	}

	fmt.Printf("[SHOP] %s bought %s x%d for %d %s\n",
		player.Name, newItem.Template.Name, quantity, totalPrice, map[int8]string{0: "Gold", 1: "Gem"}[shopItem.BuyType])

	return nil
}
