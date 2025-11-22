<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

# Viết 1 file md hoàn chỉnh mô tả cho component diagram vs deployment diagram

Dưới đây là skeleton 1 file Markdown hoàn chỉnh, bám sát kiến trúc microservices (auth, session, world, item, combat, boss, quest) với gRPC + NATS, mô tả Component Diagram và Deployment Diagram bằng Mermaid.[^1]

***

```markdown
# Kiến Trúc Triển Khai & Thành Phần – NRO Go Microservices

Tài liệu này mô tả hai khía cạnh chính của kiến trúc NRO Go khi tách thành microservices: **Component Diagram** (các service và quan hệ) và **Deployment Diagram** (cách triển khai trên hạ tầng thực tế).  
Hệ thống sử dụng **gRPC** cho các call đồng bộ, realtime và **NATS** cho pub/sub event bất đồng bộ. [attached_file:24]

---

## 1. Component Diagram

Sơ đồ sau mô tả các thành phần logic chính: API Gateway (TCP), các microservice (Auth, Session, World, Item, Combat, Boss, Quest) và các hệ thống bên ngoài (DB, NATS, Monitoring). [attached_file:24]

```

graph LR
%% ==== Clients \& Gateway ====
C[Game Client] -->|TCP Packet| GW[Game Gateway<br/>(TCP + Protocol)]

    %% ==== Core Services ====
    subgraph Core Services
        AUTH[Auth Service]
        SESSION[Session Service]
        WORLD[World Service]
        ITEM[Item/Inventory Service]
        COMBAT[Combat Service]
        BOSS[Boss Service]
        QUEST[Quest/NPC Service]
    end
    
    %% ==== Infrastructure ====
    subgraph Messaging
        NATS[(NATS<br/>Pub/Sub Bus)]
    end
    
    subgraph Databases
        DB_AUTH[(Auth DB)]
        DB_SESSION[(Session DB)]
        DB_WORLD[(World DB)]
        DB_ITEM[(Item DB)]
        DB_COMBAT[(Combat/Log DB)]
        DB_QUEST[(Quest DB)]
    end
    
    subgraph Observability
        LOG[Logging/Monitoring]
    end
    
    %% ==== Gateway → Services (gRPC) ====
    GW -->|gRPC: Login, SelectCharacter| AUTH
    GW -->|gRPC: CreateSession, UpdateSession| SESSION
    GW -->|gRPC: EnterMap, Move, Zone Query| WORLD
    GW -->|gRPC: GetInventory, Shop, Loot| ITEM
    GW -->|gRPC: UseSkill, Attack| COMBAT
    GW -->|gRPC: Boss Admin/State| BOSS
    GW -->|gRPC: InteractNPC, Quest APIs| QUEST
    
    %% ==== Service → DB ====
    AUTH --> DB_AUTH
    SESSION --> DB_SESSION
    WORLD --> DB_WORLD
    ITEM --> DB_ITEM
    COMBAT --> DB_COMBAT
    QUEST --> DB_QUEST
    
    %% ==== Events qua NATS ====
    SESSION -->|publish: SessionCreated, SessionEnded, PlayerMoved| NATS
    WORLD -->|publish: PlayerEnteredMap, MobSpawned, MobDespawned| NATS
    COMBAT -->|publish: PlayerDamaged, PlayerDied, MobKilled, PlayerGainedExp| NATS
    BOSS -->|publish: BossSpawned, BossPhaseChanged, BossKilled| NATS
    QUEST -->|publish: QuestUpdated, QuestCompleted| NATS
    
    %% ==== Services subscribe events ====
    QUEST -->|subscribe: MobKilled, PlayerGainedExp, BossKilled| NATS
    SESSION -->|subscribe: PlayerDamaged, PlayerDied| NATS
    WORLD -->|subscribe: BossSpawned, BossKilled| NATS
    ITEM -->|subscribe: MobKilled (DropLoot)| NATS
    
    %% ==== Logging/Monitoring ====
    GW --> LOG
    AUTH --> LOG
    SESSION --> LOG
    WORLD --> LOG
    ITEM --> LOG
    COMBAT --> LOG
    BOSS --> LOG
    QUEST --> LOG
    ```

### 1.1. Mô tả trách nhiệm từng component

- **Game Gateway**  
  - Nhận/gửi TCP packet với client; decode/encode protocol. [attached_file:24]  
  - Map `CMD_*` sang call gRPC tương ứng (Auth, Session, World, Item, Combat, Boss, Quest).

- **Auth Service**  
  - Đăng ký, đăng nhập, quản lý user, danh sách nhân vật.  
  - Ghi/đọc dữ liệu từ **Auth DB**.

- **Session Service**  
  - Quản lý session online (connection ↔ player, map, zone, pos, flags).  
  - Bắn event `SessionCreated`, `SessionEnded`, `PlayerMoved` lên NATS.

- **World Service**  
  - Quản lý map, zone, di chuyển, entity hiện diện trong map.  
  - Sub/publish các event liên quan tới việc vào/ra map, spawn/despawn mob.

- **Item/Inventory Service**  
  - Quản lý item, inventory, shop, drop loot.  
  - Sub `MobKilled` để sinh loot nếu cần, ghi vào **Item DB**.

- **Combat Service**  
  - Xử lý combat logic: damage, skill, effect, death.  
  - Publish các event combat (`PlayerDamaged`, `PlayerDied`, `MobKilled`…).

- **Boss Service**  
  - Quản lý AI/lifecycle boss, spawn/phase/skill pattern.  
  - Publish event boss cho World/Quest/Client.

- **Quest/NPC Service**  
  - NPC, quest, menu, logic PvE.  
  - Sub combat/world event để cập nhật tiến độ quest và publish `QuestUpdated`, `QuestCompleted`.

---

## 2. Deployment Diagram

Sơ đồ sau mô tả cách triển khai các component trên các node vật lý/container (Kubernetes / Docker / VM). Mỗi service có thể scale độc lập, giao tiếp qua gRPC (nội bộ cluster) và NATS (pub/sub). [attached_file:24]

```

graph TD
%% ==== Clients \& Edge ====
USERS[Internet Users] -->|TCP/HTTP(S)| LB[Edge / Load Balancer]

    %% ==== Gateway Layer ====
    LB --> GW_POD1[Gateway Pod 1]
    LB --> GW_POD2[Gateway Pod 2]
    LB --> GW_PODN[Gateway Pod N]
    
    subgraph Cluster "Kubernetes / Docker Cluster"
        %% ==== Service Pods ====
        subgraph Auth Tier
            AUTH1[Auth Service Pod 1]
            AUTH2[Auth Service Pod 2]
        end
    
        subgraph Session Tier
            SES1[Session Service Pod 1]
            SES2[Session Service Pod 2]
        end
    
        subgraph World Tier
            W1[World Service Pod 1]
            W2[World Service Pod 2]
        end
    
        subgraph Item Tier
            IT1[Item Service Pod 1]
            IT2[Item Service Pod 2]
        end
    
        subgraph Combat Tier
            CB1[Combat Service Pod 1]
            CB2[Combat Service Pod 2]
        end
    
        subgraph Boss Tier
            BS1[Boss Service Pod 1]
        end
    
        subgraph Quest Tier
            QT1[Quest Service Pod 1]
            QT2[Quest Service Pod 2]
        end
    
        %% ==== Messaging & Observability ====
        subgraph Infra
            NATSCL[NATS Cluster]
            LOG[Logging/Monitoring Stack]
        end
    
        %% ==== gRPC from Gateway to Services ====
        GW_POD1 -->|gRPC| AUTH1
        GW_POD1 -->|gRPC| SES1
        GW_POD1 -->|gRPC| W1
        GW_POD1 -->|gRPC| IT1
        GW_POD1 -->|gRPC| CB1
        GW_POD1 -->|gRPC| BS1
        GW_POD1 -->|gRPC| QT1
    
        GW_POD2 -->|gRPC| AUTH2
        GW_POD2 -->|gRPC| SES2
        GW_POD2 -->|gRPC| W2
        GW_POD2 -->|gRPC| IT2
        GW_POD2 -->|gRPC| CB2
        GW_POD2 -->|gRPC| QT2
    
        GW_PODN -->|gRPC (balanced)| AUTH1
        GW_PODN -->|gRPC (balanced)| SES1
        GW_PODN -->|gRPC (balanced)| W1
        GW_PODN -->|gRPC (balanced)| IT1
        GW_PODN -->|gRPC (balanced)| CB1
        GW_PODN -->|gRPC (balanced)| QT1
    
        %% ==== Services <-> NATS ====
        AUTH1 --> NATSCL
        AUTH2 --> NATSCL
        SES1 --> NATSCL
        SES2 --> NATSCL
        W1 --> NATSCL
        W2 --> NATSCL
        IT1 --> NATSCL
        IT2 --> NATSCL
        CB1 --> NATSCL
        CB2 --> NATSCL
        BS1 --> NATSCL
        QT1 --> NATSCL
        QT2 --> NATSCL
    
        %% ==== Logging / Monitoring ====
        GW_POD1 --> LOG
        GW_POD2 --> LOG
        GW_PODN --> LOG
        AUTH1 --> LOG
        AUTH2 --> LOG
        SES1 --> LOG
        SES2 --> LOG
        W1 --> LOG
        W2 --> LOG
        IT1 --> LOG
        IT2 --> LOG
        CB1 --> LOG
        CB2 --> LOG
        BS1 --> LOG
        QT1 --> LOG
        QT2 --> LOG
    end
    
    %% ==== Databases (managed outside cluster or as statefulset) ====
    subgraph Databases
        DBA[(Auth DB)]
        DBS[(Session DB)]
        DBW[(World DB)]
        DBI[(Item DB)]
        DBC[(Combat/Log DB)]
        DBQ[(Quest DB)]
    end
    
    AUTH1 --> DBA
    AUTH2 --> DBA
    
    SES1 --> DBS
    SES2 --> DBS
    
    W1 --> DBW
    W2 --> DBW
    
    IT1 --> DBI
    IT2 --> DBI
    
    CB1 --> DBC
    CB2 --> DBC
    
    QT1 --> DBQ
    QT2 --> DBQ
    ```

### 2.1. Các quyết định triển khai chính

- **Scale Gateway theo số connection**: Gateway giữ TCP connection và translate sang gRPC, vì vậy thường scale theo lượng người chơi online. [attached_file:24]  
- **Mỗi microservice scale độc lập**: Combat/World thường nặng CPU/RAM hơn, có thể scale nhiều pod hơn so với Auth hoặc Quest.  
- **NATS cluster tách riêng**: Có thể triển khai thành 3 node để đảm bảo HA cho pub/sub event.  
- **Database tách theo bounded context**: mỗi service sở hữu schema riêng (Auth DB, World DB, …) để giảm coupling, phù hợp microservices.

---

> Bạn có thể copy nguyên file này vào `docs/Component_Deployment.md` (hoặc tương đương) trong repo, sau đó chỉnh sửa thêm chi tiết cụ thể (tên namespace, domain, host, port) cho khớp môi trường của bạn.
```

