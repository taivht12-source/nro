package ports

import "nro-go/internal/core/domain"

type ShopRepository interface {
	GetShop(id int) (*domain.Shop, error)
	GetShopByNPC(npcID int) (*domain.Shop, error)
}
