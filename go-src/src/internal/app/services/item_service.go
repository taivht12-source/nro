package services

import (
	"fmt"
	"nro/src/internal/core/domain"
	"nro/src/internal/core/ports"
	"sync"
	"time"
)

// ItemService quản lý item templates và tạo items.
//
// EXPLANATION:
// ItemService chịu trách nhiệm:
// 1. Load và cache thông tin mẫu vật phẩm (ItemTemplate) từ Database.
// 2. Tạo ra các instance vật phẩm (Item) mới từ template.
// 3. Cung cấp thông tin chi tiết về vật phẩm dựa trên ID.
type ItemService struct {
	templates map[int]*domain.ItemTemplate
	repo      ports.ItemRepository
	mu        sync.RWMutex
}

var itemServiceInstance *ItemService
var itemOnce sync.Once

// GetItemService returns singleton instance.
func GetItemService() *ItemService {
	itemOnce.Do(func() {
		itemServiceInstance = &ItemService{
			templates: make(map[int]*domain.ItemTemplate),
		}
		// Load mock items initially, but will be overridden if repo is set
		// itemServiceInstance.loadMockItems()
	})
	return itemServiceInstance
}

func (is *ItemService) SetRepository(repo ports.ItemRepository) {
	is.mu.Lock()
	defer is.mu.Unlock()
	is.repo = repo
}

// loadMockItems loads default item templates.
func (is *ItemService) loadMockItems() {
	mockItems := []*domain.ItemTemplate{
		// Weapons (Type 0)
		{ID: 1, Name: "Gậy Như Ý", Type: 0, Gender: -1, Description: "Vũ khí huyền thoại", IconID: 1, Part: 0},
		{ID: 2, Name: "Kiếm Z", Type: 0, Gender: -1, Description: "Thanh kiếm sắc bén", IconID: 2, Part: 0},
		{ID: 3, Name: "Búa Thor", Type: 0, Gender: -1, Description: "Búa thần sấm", IconID: 3, Part: 0},

		// Armor (Type 1)
		{ID: 10, Name: "Áo Kame", Type: 1, Gender: 0, Description: "Áo của Quy Lão", IconID: 10, Part: 1},
		{ID: 11, Name: "Áo Namek", Type: 1, Gender: 1, Description: "Áo của người Namek", IconID: 11, Part: 1},
		{ID: 12, Name: "Áo Saiyan", Type: 1, Gender: 2, Description: "Áo chiến binh Saiyan", IconID: 12, Part: 1},

		// Accessories (Type 2)
		{ID: 20, Name: "Mắt Kính Scouter", Type: 2, Gender: -1, Description: "Đo sức mạnh", IconID: 20, Part: 2},
		{ID: 21, Name: "Tai Nghe Kaio", Type: 2, Gender: -1, Description: "Tai nghe thần linh", IconID: 21, Part: 2},
		{ID: 22, Name: "Nhẫn Potara", Type: 2, Gender: -1, Description: "Nhẫn hợp thể", IconID: 22, Part: 2},

		// Consumables (Type 3)
		{ID: 30, Name: "Đậu Thần", Type: 3, Gender: -1, Description: "Hồi phục HP/MP", IconID: 30},
		{ID: 31, Name: "Bí Ngô", Type: 3, Gender: -1, Description: "Tăng HP tạm thời", IconID: 31},
		{ID: 32, Name: "Cà Rốt", Type: 3, Gender: -1, Description: "Tăng MP tạm thời", IconID: 32},

		// Dragon Balls (Type 4)
		{ID: 40, Name: "Ngọc Rồng 1 Sao", Type: 4, Gender: -1, Description: "Viên ngọc rồng", IconID: 40},
		{ID: 41, Name: "Ngọc Rồng 2 Sao", Type: 4, Gender: -1, Description: "Viên ngọc rồng", IconID: 41},
		{ID: 42, Name: "Ngọc Rồng 3 Sao", Type: 4, Gender: -1, Description: "Viên ngọc rồng", IconID: 42},
	}

	for _, template := range mockItems {
		is.templates[template.ID] = template
	}

	fmt.Printf("[ITEM] Loaded %d mock item templates\n", len(is.templates))
}

// GetTemplate returns item template by ID.
func (is *ItemService) GetTemplate(id int) *domain.ItemTemplate {
	is.mu.RLock()
	template, ok := is.templates[id]
	is.mu.RUnlock()

	if ok {
		return template
	}

	// If not in cache and repo is set, try to load from repo
	if is.repo != nil {
		t, err := is.repo.GetTemplate(id)
		if err == nil && t != nil {
			is.mu.Lock()
			is.templates[id] = t
			is.mu.Unlock()
			return t
		}
	}

	return nil
}

// GetAllTemplates returns all item templates.
func (is *ItemService) GetAllTemplates() []*domain.ItemTemplate {
	is.mu.RLock()
	defer is.mu.RUnlock()

	templates := make([]*domain.ItemTemplate, 0, len(is.templates))
	for _, t := range is.templates {
		templates = append(templates, t)
	}
	return templates
}

// CreateItem creates a new item instance from template.
func (is *ItemService) CreateItem(templateID int, quantity int) *domain.Item {
	template := is.GetTemplate(templateID)
	if template == nil {
		fmt.Printf("[ITEM] Template %d not found\n", templateID)
		return nil
	}

	return &domain.Item{
		ID:       0, // Will be set when added to inventory
		Template: template,
		Quantity: quantity,
		Options:  make([]*domain.ItemOption, 0),
		CreateAt: time.Now().Unix(),
	}
}

// CreateItemWithOptions creates item with specific options (stats).
func (is *ItemService) CreateItemWithOptions(templateID int, quantity int, options []*domain.ItemOption) *domain.Item {
	item := is.CreateItem(templateID, quantity)
	if item != nil {
		item.Options = options
	}
	return item
}

// GetItemsByType returns all templates of a specific type.
func (is *ItemService) GetItemsByType(itemType int8) []*domain.ItemTemplate {
	is.mu.RLock()
	defer is.mu.RUnlock()

	items := make([]*domain.ItemTemplate, 0)
	for _, t := range is.templates {
		if t.Type == itemType {
			items = append(items, t)
		}
	}
	return items
}
