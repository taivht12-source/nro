package domain

// Shop cửa hàng.
type Shop struct {
	ID      int
	Name    string
	Type    int8 // 0: Normal, 1: Consign?
	NPCID   int  // NPC quản lý shop
	Items   []*ShopItem
	TabName []string // Tên các tab trong shop
}

// ShopItem vật phẩm trong shop.
type ShopItem struct {
	ID         int
	ShopID     int
	TemplateID int
	Template   *ItemTemplate
	Price      int64
	BuyType    int8 // 0: Gold, 1: Gem
	Quantity   int  // Số lượng (nếu có giới hạn, thường là -1 = vô hạn)
	IsNew      bool
}
