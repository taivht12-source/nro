package main

import (
	"fmt"
	"nro/src/internal/app/commands"
	"nro/src/internal/app/services"
	"nro/src/internal/core/ports"
	"nro/src/internal/infrastructure/config"
	"nro/src/internal/infrastructure/network"
	"nro/src/internal/infrastructure/persistence"
)

func main() {
	fmt.Println("Starting NRO Go Server...")

	// 1. Load Config
	// Giả sử file Config.properties nằm ở thư mục gốc (cùng cấp với go-src hoặc thư mục chạy)
	// Cần trỏ đúng đường dẫn. Ở đây hardcode tạm đường dẫn tuyệt đối hoặc tương đối.
	// Trong thực tế nên dùng flag hoặc env.
	err := config.Load("../../../Config.properties") // Relative path từ cmd/server
	if err != nil {
		fmt.Printf("Warning: Cannot load config: %v. Using defaults.\n", err)
	}

	// 2. Connect Database (Optional)
	cfg := config.Get()
	var userRepo ports.UserRepository
	var playerRepo ports.PlayerRepository
	var loginHandler *commands.LoginHandler
	var charService *services.CharacterService

	if cfg.RequireDB {
		if err := persistence.ConnectDB(); err != nil {
			panic(fmt.Sprintf("Cannot connect to DB: %v", err))
		}
		defer persistence.CloseDB()
		fmt.Println("✓ Database connected")

		// Initialize Components with DB
		db := persistence.GetDB()
		itemRepo := persistence.NewMySQLItemRepository(db)
		skillRepo := persistence.NewMySQLSkillRepository(db)
		userRepo = persistence.NewMySQLUserRepository(db)
		playerRepo = persistence.NewMySQLPlayerRepository(db, itemRepo, skillRepo)
		mapRepo := persistence.NewMySQLMapRepository(db)

		loginHandler = commands.NewLoginHandler(userRepo)
		charService = services.NewCharacterService(playerRepo)

		// Inject MapRepo into MapService
		services.GetMapService().SetRepository(mapRepo)

		// Inject ItemRepo into ItemService and InventoryService
		services.GetItemService().SetRepository(itemRepo)
		services.GetInventoryService().SetItemRepository(itemRepo)

		// Initialize Shop System
		shopRepo := persistence.NewMySQLShopRepository(db, itemRepo)
		services.GetShopService().SetRepository(shopRepo)

		fmt.Printf("✓ LoginHandler, CharacterService, ZoneService, ItemService, and ShopService initialized\n")
	} else {
		fmt.Println("⚠ Running in NO-DB mode (testing only)")
		fmt.Println("  Database connection is disabled. Login will use mock data.")
		// CharacterService will be nil, controller will handle mock login
	}

	// 3. Start TCP Server
	// Init Controller with dependencies
	controller := network.NewController(loginHandler, charService)

	// Initialize services (load mock data)
	services.GetMapService()   // Load maps
	services.GetItemService()  // Load items
	services.GetSkillService() // Load skills
	services.GetNPCService()   // Load NPCs
	services.GetTaskService()  // Load tasks

	server := network.NewTCPServer(cfg.ServerPort, controller)

	fmt.Printf("🚀 Starting NRO Go Server on port %d...\n", cfg.ServerPort)
	if err := server.Start(); err != nil {
		panic(err)
	}
}
