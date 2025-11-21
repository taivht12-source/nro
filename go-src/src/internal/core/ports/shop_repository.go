package ports

import "nro/src/internal/core/domain"

type ShopRepository interface {
	GetShop(id int) (*domain.Shop, error)
	GetShopByNPC(npcID int) (*domain.Shop, error)
}
