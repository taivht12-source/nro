```
# NRO Go - 4+1 Architecture Views

Tài liệu này trình bày 5 view kiến trúc dựa trên 4+1 View Model để mô tả tổng thể hệ thống.

---

## 1. Logical View (View Logic)

Mô tả sự phân rã chức năng của hệ thống thành các module chính và mối quan hệ giữa chúng, dựa trên Clean Architecture.

```

graph TD
subgraph Core Domain (internal/core)
D[Domain Models (Entities)]
P[Ports (Interfaces/Repositories)]
end

    subgraph Application Layer (internal/app)
        S[Services (Business Logic)]
        C[Commands (CQRS Handlers)]
    end
    
    subgraph Infrastructure Layer (internal/infrastructure)
        N[Network (TCP Controller)]
        I[Persistence (DB Implementation)]
        CF[Config]
    end
    
    subgraph Entry Points (cmd/)
        E[Server Entry Point]
    end
    
    subgraph Shared Packages (pkg/)
        PR[Protocol Definitions]
    end
    
    subgraph External
        DB[(MySQL Database)]
        CL[Client/Test Client]
    end
    
    CL --> N
    N --> C
    N --> S
    
    C --> P
    S --> P
    
    P --> D
    D --> P
    
    I --> P
    I --> DB
    
    E --> CF
    E --> N
    E --> I
    E --> C
    E --> S
    
    N --> PR
    CL --> PR
    
    style D fill:#e8f5e9,stroke:#4caf50
    style P fill:#e8f5e9,stroke:#4caf50
    style S fill:#f3e5f5,stroke:#9c27b0
    style C fill:#f3e5f5,stroke:#9c27b0
    style N fill:#fff3e0,stroke:#ff9800
    style I fill:#fce4ec,stroke:#e91e63
    style E fill:#e1f5ff,stroke:#2196f3
    style DB fill:#fce4ec,stroke:#e91e63
    style CL fill:#e1f5ff,stroke:#2196f3
    style PR fill:#ffecb3,stroke:#ffc107
    ```

---

## 2. Process View (View Tiến Trình)

Mô tả các tiến trình (processes) của hệ thống, cách chúng giao tiếp và luồng dữ liệu chính.

```

graph LR
subgraph Client Interaction
A[Client] -->|TCP Packet| B(Network Controller)
end

    subgraph Server Processes
        B --> C{Command/Service Router}
        C -->|Login Command| D[Login Handler (CQRS)]
        C -->|Game Action| E[Game Service (e.g., Combat, Map)]
        D --> F[User Repository]
        E --> F
        F --> G[(Database)]
        G --> F
        F --> D
        F --> E
        E -->|State Update| B
        D -->|Login Success| B
    end
    
    B -->|Send Response| A
    
    style A fill:#e1f5ff
    style B fill:#fff3e0
    style C fill:#ffecb3
    style D fill:#f3e5f5
    style E fill:#f3e5f5
    style F fill:#e8f5e9
    style G fill:#fce4ec
    ```

---

## 3. Development View (View Phát Triển)

Mô tả cấu trúc module và tổ chức mã nguồn, dựa trên cấu trúc thư mục.

```

graph TD
A[go-src/src/]
A --> B(cmd/)
A --> C(internal/)
A --> D(pkg/)
A --> E(proto/)

    B --> B1[server/]
    B --> B2[test_client/]
    
    C --> C1[app/]
    C --> C2[core/]
    C --> C3[infrastructure/]
    
    C1 --> C1a[commands/]
    C1 --> C1b[services/]
    
    C2 --> C2a[domain/]
    C2 --> C2b[ports/]
    
    C3 --> C3a[config/]
    C3 --> C3b[network/]
    C3 --> C3c[persistence/]
    C3 --> C3d[session/]
    
    D --> D1[protocol/]
    
    E --> E1[auth/]
    E --> E2[item/]
    E --> E3[world/]
    
    style A fill:#bbdefb
    style B fill:#e1f5ff
    style C fill:#e1f5ff
    style D fill:#e1f5ff
    style E fill:#e1f5ff
    style C1 fill:#f3e5f5
    style C2 fill:#e8f5e9
    style C3 fill:#fce4ec
    ```

---

## 4. Physical View (View Vật Lý)

Mô tả sự triển khai hệ thống trên các thành phần phần cứng (hoặc máy ảo/container).

```

graph LR
subgraph Production Environment
LB[Load Balancer/Proxy] -->|TCP/HTTP| S1(NRO Go Server Instance 1)
LB -->|TCP/HTTP| S2(NRO Go Server Instance 2)
LB -->|TCP/HTTP| S3(NRO Go Server Instance N)
end

    S1 -->|MySQL Protocol| DB[(Database Server)]
    S2 -->|MySQL Protocol| DB
    S3 -->|MySQL Protocol| DB
    
    subgraph Monitoring
        M[Monitoring/Logging System]
    end
    
    S1 -->|Logs/Metrics| M
    S2 -->|Logs/Metrics| M
    S3 -->|Logs/Metrics| M
    
    style LB fill:#ffecb3
    style S1 fill:#e1f5ff
    style S2 fill:#e1f5ff
    style S3 fill:#e1f5ff
    style DB fill:#fce4ec
    style M fill:#e8f5e9
    ```

---

## 5. Scenarios View (View Kịch Bản/Use Cases)

Mô tả luồng tương tác chính của người dùng, minh họa cách các thành phần trong các view khác hoạt động cùng nhau.

### Kịch Bản: Đăng Nhập Người Chơi

```

sequenceDiagram
actor Client
participant NetworkController
participant LoginHandler
participant UserRepository
participant Database

    Client->>NetworkController: Gửi gói tin Login (Username, Password)
    NetworkController->>LoginHandler: Gọi Handle(LoginCommand)
    LoginHandler->>UserRepository: Gọi GetUserByUsername(Username)
    UserRepository->>Database: SELECT * FROM users WHERE username = ?
    Database-->>UserRepository: Trả về User Record
    UserRepository-->>LoginHandler: Trả về Domain User Entity
    LoginHandler->>LoginHandler: Xác thực Password
    alt Xác thực thành công
        LoginHandler->>NetworkController: Trả về User Entity
        NetworkController->>Client: Gửi gói tin Login Success
    else Xác thực thất bại
        LoginHandler->>NetworkController: Trả về Error
        NetworkController->>Client: Gửi gói tin Login Failed
    end

```
