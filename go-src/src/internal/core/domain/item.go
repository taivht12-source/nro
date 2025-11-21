package domain

// ItemTemplate mẫu vật phẩm.
type ItemTemplate struct {
	ID           int
	Type         int8
	Gender       int8
	Name         string
	Description  string
	Level        int
	IconID       int
	Part         int16
	IsUpToUp     bool
	PowerRequire int64
	Gold         int
	Gem          int
	Head         int
	Body         int
	Leg          int
}

// Item vật phẩm cụ thể (trong hành trang).
type Item struct {
	ID         int // Unique ID in inventory (optional, or just index)
	TemplateID int
	Template   *ItemTemplate
	Quantity   int
	Options    []*ItemOption
	CreateAt   int64
}

// ItemOption chỉ số của vật phẩm.
type ItemOption struct {
	ID    int
	Param int
}

// Inventory chứa thông tin hành trang của người chơi.
type Inventory struct {
	Gold int64
	Gem  int

	ItemsBody []*Item // Trang bị đang mặc
	ItemsBag  []*Item // Hành trang
	ItemsBox  []*Item // Rương đồ
}
