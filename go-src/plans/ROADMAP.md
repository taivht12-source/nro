upda# Implementation Plan - Full NRO Conversion to Go (Hexagonal + CQRS)

## Goal

Completely convert the NRO Java game server to a Go-based Microservices architecture using Hexagonal (Ports & Adapters) and CQRS patterns.

## User Review Required
>
> [!IMPORTANT]
> This is a massive migration. The plan is broken down into **7 Phases**. We will execute **Phase 1** first.
> The architecture will be strict:
>
> - **Domain**: Pure Go structs, no dependencies.
> - **Application**: Use Cases / CQRS Commands & Queries.
> - **Infrastructure**: Database (MySQL), Network (TCP), External APIs.
> - **Interfaces (Ports)**: Define how layers communicate.

## Phase 1: Foundation & Network Layer (Completed)

**Goal**: Establish the project skeleton and enable Client-Server communication (Handshake).

- [ ] **Movement Logic**: Handle `PLAYER_MOVE` packet.
- [ ] **State Management**: Sync player position to other players in the same zone.

## Phase 2: Core Domain & Authentication (Completed)

**Goal**: Define core entities and implement Login flow.

- [x] **Entities**: `User`, `Player`, `Session`.
- [x] **Database**: MySQL Connection, `UserRepository`.
- [x] **Login**: `LoginCommand`, `LoginHandler`, Password hashing.

## Phase 3: World & Movement System (Completed)

**Goal**: Manage Maps, Zones, and Player Movement.

- [x] **Entities**: `MapTemplate`, `Zone`, `Waypoint`, `Mob`.
- [x] **Movement Logic**: Handle `PLAYER_MOVE` (-7) in `Controller`.
- [x] **State Management**: `ZoneService` & `SessionManager` for broadcasting position.
- [x] **Persistence**: Implement `MySQLMapRepository` (Currently Stub).
- [x] **Map Transition**: Handle `Waypoint` logic.

## Phase 4: Item & Inventory System (Completed)

**Goal**: Manage items, equipment, and shops.

- [x] **Item Templates**: Load `ItemTemplate` from DB/Config.
- [x] **Inventory**: Port `Inventory`, `Item` logic (`inventory_service.go` created).
- [x] **Shop System**: Implement `ShopService`, `ConsignShopService`.

## Phase 5: Skills & Combat System (Completed)

**Goal**: Implement fighting mechanics.

- [x] **Skill System**: Port `SkillService`, `SkillTemplate` (`skill_service.go` created).
- [x] **Combat Logic**: Calculate damage, HP/MP updates (`combat_service.go` created).
- [x] **Effect System**: Handle buffs/debuffs.

## Phase 6: NPCs & Tasks (Completed)

**Goal**: Interaction with the world.

- [x] **NPC System**: Port `NpcService`, `MenuController` (`npc_service.go` created).
- [x] **Task System**: Port `TaskService` (`task_service.go` created).
- [x] **Menu System**: Implement dynamic menus and dialogue.
- [x] **Network Refactoring**: Introduce Command Enums and Packet DTOs.

## Phase 7: Advanced Features (Completed)

**Goal**: Complete the game features.

- [x] **Bosses**: Port `BossManager` with AI system.
- [x] **Clans**: Port `ClanService`.
- [x] **Events**: Port `Shenron`, `Budokai`, etc.
- [x] **Boss AI**: Implement specific boss behaviors (Broly, Cell, Frieza, Android, Black Goku).

## Phase 8: Microservices & gRPC (New)

**Goal**: Split monolithic server into microservices communicating via gRPC.

- [ ] **Protobuf Definition**: Define `.proto` files for inter-service communication (Auth, World, Item).
- [ ] **gRPC Server**: Implement gRPC servers for each domain.
- [ ] **gRPC Client**: Implement clients to consume services.
- [ ] **API Gateway**: Create a gateway to route client packets to appropriate services.

---

## Proposed Changes (Phase 1 Detail)

### [NEW] Project Structure

```
go-src/src/
├── cmd/
│   └── server/          # Main entry point
├── internal/
│   ├── core/            # Hexagonal Core
│   │   ├── domain/      # Entities (Player, Session)
│   │   └── ports/       # Interfaces (Repository, Network)
│   ├── app/             # Application Layer
│   │   ├── commands/    # CQRS Commands (Login, Move)
│   │   └── queries/     # CQRS Queries (GetInfo)
│   └── infrastructure/  # Adapters
│       ├── persistence/ # MySQL Repositories
│       └── network/     # TCP Handler (Controller)
└── pkg/
    └── protocol/        # Low-level Network (Message, Cipher)
```

### [NEW] [go.mod](file:///c:/Users/PC/Downloads/NRO_BLACK_GOKU/NRO_BLACK_GOKU/go-src/src/go.mod)

- Initialize Go module `nro-go`.

### [NEW] [main.go](file:///c:/Users/PC/Downloads/NRO_BLACK_GOKU/NRO_BLACK_GOKU/go-src/src/cmd/server/main.go)

- Initialize `TCPServer` and Dependency Injection container.

### [NEW] [message.go](file:///c:/Users/PC/Downloads/NRO_BLACK_GOKU/NRO_BLACK_GOKU/go-src/src/pkg/protocol/message.go)

- Implement `Read/Write` methods matching Java's `DataInputStream/DataOutputStream`.

### [NEW] [session.go](file:///c:/Users/PC/Downloads/NRO_BLACK_GOKU/NRO_BLACK_GOKU/go-src/src/pkg/protocol/session.go)

- Manage socket connection and encryption keys.

## Verification Plan (Phase 1)

- **Unit Test**: Verify `Message` encoding/decoding matches Java's byte output.
- **Integration Test**: Start Go Server -> Connect with Java Client (or Test Client) -> Verify Handshake success.
