package persistence

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"nro-go/internal/core/domain"
	"nro-go/internal/core/ports"
)

type MySQLShopRepository struct {
	db       *sql.DB
	itemRepo ports.ItemRepository
}

func NewMySQLShopRepository(db *sql.DB, itemRepo ports.ItemRepository) ports.ShopRepository {
	return &MySQLShopRepository{db: db, itemRepo: itemRepo}
}

func (r *MySQLShopRepository) GetShop(id int) (*domain.Shop, error) {
	query := `SELECT id, name, type, npc_id, items, tab_name FROM shop WHERE id = ?`
	row := r.db.QueryRow(query, id)
	return r.scanShop(row)
}

func (r *MySQLShopRepository) GetShopByNPC(npcID int) (*domain.Shop, error) {
	query := `SELECT id, name, type, npc_id, items, tab_name FROM shop WHERE npc_id = ?`
	row := r.db.QueryRow(query, npcID)
	return r.scanShop(row)
}

func (r *MySQLShopRepository) scanShop(row *sql.Row) (*domain.Shop, error) {
	var shop domain.Shop
	var itemsJSON, tabNameJSON string

	err := row.Scan(&shop.ID, &shop.Name, &shop.Type, &shop.NPCID, &itemsJSON, &tabNameJSON)
	if err != nil {
		return nil, err
	}

	shop.Items = r.parseShopItems(itemsJSON, shop.ID)
	json.Unmarshal([]byte(tabNameJSON), &shop.TabName)

	return &shop, nil
}

func (r *MySQLShopRepository) parseShopItems(data string, shopID int) []*domain.ShopItem {
	var rawList []string
	if err := json.Unmarshal([]byte(data), &rawList); err != nil {
		return nil
	}

	var items []*domain.ShopItem
	for _, s := range rawList {
		var itemData []interface{}
		if err := json.Unmarshal([]byte(s), &itemData); err != nil {
			continue
		}

		if len(itemData) < 4 {
			continue
		}

		// Format: [TemplateID, Price, BuyType, Quantity, IsNew]
		templateID := int(itemData[0].(float64))
		price := int64(itemData[1].(float64))
		buyType := int8(itemData[2].(float64))
		quantity := int(itemData[3].(float64))
		isNew := false
		if len(itemData) > 4 {
			isNew = itemData[4].(bool)
		}

		template, _ := r.itemRepo.GetTemplate(templateID)
		if template == nil {
			template = &domain.ItemTemplate{ID: templateID, Name: fmt.Sprintf("Unknown %d", templateID)}
		}

		items = append(items, &domain.ShopItem{
			ID:         len(items), // Generate temp ID
			ShopID:     shopID,
			TemplateID: templateID,
			Template:   template,
			Price:      price,
			BuyType:    buyType,
			Quantity:   quantity,
			IsNew:      isNew,
		})
	}
	return items
}
