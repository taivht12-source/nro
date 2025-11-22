```
# Microservice Component Diagram - NRO Go Game Server

## Tổng quan thành phần

- Game Client
- Gateway (TCP + Protocol)
- Auth Service
- Session Service
- World Service
- Item/Inventory Service
- Combat Service
- Boss Service
- Quest/NPC Service
- NATS (Message Bus)
- Các DB: Auth DB, Session DB, World DB, Item DB, Combat/Log DB, Quest DB

---

## Sơ đồ component dạng text

```

[Game Client]
|
| TCP Packets (LOGIN, SELECT_CHAR, ENTER_MAP, MOVE, ATTACK/USE_SKILL, ...)
v
[Gateway - TCP + Protocol]

- Decode/encode NRO packet
- Map CMD_* -> gRPC call
- Quản lý connection vật lý (socket, heartbeat)

```
    +-- gRPC --> [Auth Service]
    |               Interfaces:
    |               - Login(credentials) -> {userId, accessToken, characters[]}
    |               - ValidateToken(accessToken) -> {userId, valid}
    |               - SelectCharacter(userId, characterId) -> {characterState}
    |               - Register(...)
    |               - ListCharacters(userId) -> {characters[]}
    |
    |               Internal components:
    |               - UserController / AuthController
    |               - UserService
    |               - PlayerService
    |               - TokenManager (JWT / Session token)
    |               - Repositories:
    |                   * UserRepository
    |                   * PlayerRepository
    |
    |               Persistence:
    |               - Auth DB (users, players, credentials)
    |
    +-- gRPC --> [Session Service]
    |               Interfaces:
    |               - CreateSession(userId, characterId, initialState) -> {sessionId}
    |               - GetSession(sessionId) -> {playerId, mapId, zoneId, pos, flags}
    |               - UpdateSessionState(sessionId, statePatch)
    |               - EndSession(sessionId)
    |               - KickSession(playerId)
    |
    |               Internal components:
    |               - SessionManager (in-memory + cache + DB)
    |               - SessionValidator
    |               - SessionEventPublisher
    |               - Repositories:
    |                   * SessionRepository
    |
    |               Persistence:
    |               - Session DB (session snapshot, last map/pos, status)
    |
    |               Events (publish to NATS):
    |               - SessionCreated {sessionId, playerId, mapId, zoneId}
    |               - SessionEnded   {sessionId, playerId}
    |               - PlayerMoved    {playerId, mapId, zoneId, pos}
    |
    |               Events (subscribe from NATS):
    |               - PlayerDamaged  {playerId, hp}
    |               - PlayerDied     {playerId}
    |
    +-- gRPC --> [World Service]
    |               Interfaces:
    |               - EnterMap(sessionId, mapId) -> {zoneId, spawnEntities[]}
    |               - Move(sessionId, fromPos, toPos) -> {finalPos, visibleEntities[]}
    |               - ListPlayersInZone(zoneId) -> {players[]}
    |               - GetMapInfo(mapId) -> {zones, npcs, mobs, bossSpawn}
    |               - UpdateMobState(mobId, statePatch)
    |               - QueryMob(mobId) -> {mobState}
    |
    |               Internal components:
    |               - MapManager
    |               - ZoneManager
    |               - MovementValidator
    |               - VisibilityService (AOI)
    |               - EntityRegistry (players, mobs, npcs, drops)
    |               - Repositories:
    |                   * MapRepository
    |                   * ZoneRepository
    |                   * MobRepository
    |
    |               Persistence:
    |               - World DB (maps, zones, mob templates, npc, spawn points)
    |
    |               Events (publish to NATS):
    |               - PlayerEnteredMap {playerId, mapId, zoneId}
    |               - PlayerLeftMap    {playerId, mapId, zoneId}
    |               - MobSpawned       {mobId, mapId, zoneId, pos}
    |               - MobDespawned     {mobId}
    |
    |               Events (subscribe from NATS):
    |               - BossSpawned      {bossId, mapId, zoneId, pos}
    |               - BossKilled       {bossId}
    |
    +-- gRPC --> [Item / Inventory Service]
    |               Interfaces:
    |               - GetInventory(playerId) -> {items[]}
    |               - AddItem(playerId, itemId, quantity)
    |               - RemoveItem(playerId, itemId, quantity)
    |               - MoveItem(playerId, fromSlot, toSlot)
    |               - OpenShop(playerId, shopId) -> {shopItems[]}
    |               - BuyItem(playerId, shopId, itemId, quantity)
    |               - SellItem(playerId, itemId, quantity)
    |               - GenerateLoot(mobId, playerId) -> {lootItems[]}
    |
    |               Internal components:
    |               - InventoryManager
    |               - ShopManager
    |               - LootGenerator
    |               - ItemTemplateProvider
    |               - Repositories:
    |                   * InventoryRepository
    |                   * ItemTemplateRepository
    |                   * ShopRepository
    |
    |               Persistence:
    |               - Item DB (inventory, item templates, shop configs)
    |
    |               Events (publish to NATS):
    |               - PlayerItemGained {playerId, items[]}
    |               - PlayerItemLost   {playerId, items[]}
    |
    |               Events (subscribe from NATS):
    |               - MobKilled {mobId, killerId, mapId, zoneId}
    |                   -> GenerateLoot + update inventory
    |
    +-- gRPC --> [Combat Service]
    |               Interfaces:
    |               - UseSkill(sessionId, skillId, targetInfo)
    |               - Attack(sessionId, targetInfo)
    |               - GetCombatState(sessionId)
    |
    |               Internal components:
    |               - CombatEngine
    |               - SkillEngine
    |               - EffectEngine
    |               - DamageCalculator
    |               - CooldownManager
    |               - Repositories:
    |                   * SkillRepository
    |                   * EffectRepository
    |
    |               Collaboration (gRPC outbound):
    |               - -> WorldService.QueryMob(...)
    |               - -> WorldService.UpdateMobState(...)
    |               - -> ItemService.GenerateLoot(...) (hoặc publish event cho Item)
    |
    |               Persistence:
    |               - Combat/Log DB (combat log, damage stats, anti-cheat signals)
    |
    |               Events (publish to NATS):
    |               - PlayerDamaged   {playerId, deltaHp, hpLeft}
    |               - PlayerDied      {playerId, killerId, mapId, zoneId}
    |               - MobDamaged      {mobId, deltaHp, hpLeft}
    |               - MobKilled       {mobId, killerId, mapId, zoneId}
    |               - PlayerGainedExp {playerId, exp, newLevel?}
    |               - PlayerLevelUp   {playerId, newLevel}
    |
    +-- gRPC --> [Boss Service]
    |               Interfaces:
    |               - SpawnBoss(bossId, mapId, zoneId, pos)
    |               - DespawnBoss(bossInstanceId)
    |               - GetBossState(bossInstanceId)
    |               - AdminControlBoss(...)
    |
    |               Internal components:
    |               - BossManager
    |               - BossAIController
    |               - BossScheduler (spawn timing)
    |               - Repositories:
    |                   * BossTemplateRepository
    |                   * BossSpawnConfigRepository
    |
    |               Persistence:
    |               - World DB / Boss tables (spawn config, template)
    |
    |               Events (publish to NATS):
    |               - BossSpawned      {bossInstanceId, bossId, mapId, zoneId, pos}
    |               - BossPhaseChanged {bossInstanceId, phase}
    |               - BossKilled       {bossInstanceId, bossId, killerId}
    |
    +-- gRPC --> [Quest / NPC Service]
                    Interfaces:
                    - InteractNPC(sessionId, npcId, action) -> {dialog, menu, questOptions}
                    - AcceptQuest(playerId, questId)
                    - CompleteQuest(playerId, questId)
                    - GetPlayerQuests(playerId) -> {quests[]}
                    - TriggerMenuAction(playerId, menuId, optionId)
    
                    Internal components:
                    - NPCManager
                    - QuestManager
                    - ConditionChecker (kill count, item, level, boss kill, location)
                    - RewardApplier (call Item/World/Session nếu cần)
                    - Repositories:
                        * NPCRepository
                        * QuestRepository
                        * PlayerQuestRepository
    
                    Persistence:
                    - Quest DB (quest templates, npc dialog, player-quest state)
    
                    Events (publish to NATS):
                    - QuestUpdated   {playerId, questId, progress}
                    - QuestCompleted {playerId, questId, rewards[]}
    
                    Events (subscribe from NATS):
                    - MobKilled       {mobId, killerId}
                    - PlayerGainedExp {playerId, exp}
                    - BossKilled      {bossId, killerId}
                    - PlayerItemGained {playerId, items[]}
```

---

## Hạ tầng chung - Message Bus NATS

```

[NATS Message Bus]
- Channel / Subject chính:
* session.*          (SessionCreated, SessionEnded, PlayerMoved)
* combat.*           (PlayerDamaged, PlayerDied, MobKilled, PlayerGainedExp, PlayerLevelUp)
* world.*            (PlayerEnteredMap, PlayerLeftMap, MobSpawned, MobDespawned)
* boss.*             (BossSpawned, BossPhaseChanged, BossKilled)
* quest.*            (QuestUpdated, QuestCompleted)
* item.*             (PlayerItemGained, PlayerItemLost)

    - Consumers:
        * Session Service    (session.*, combat.PlayerDamaged, combat.PlayerDied)
        * World Service      (boss.BossSpawned, boss.BossKilled)
        * Item Service       (combat.MobKilled)
        * Quest Service      (combat.*, item.PlayerItemGained, boss.BossKilled)
        * Analytics/Log      (subscribe all *)
    ```

---

Bạn có thể tải file này về để sử dụng. Nếu cần, có thể hỗ trợ tạo file trên một nền tảng chia sẻ hoặc gửi qua email cho bạn.```

