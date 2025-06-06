# 🔐 Keeper 
---

## 📁 Repository Structure

```

accountability-platform/
├── services/              # Microservices live here
│   ├── auth/              # Handles user/guardian authentication and roles
│   ├── session-manager/   # Creates, manages, and tracks Guardian-User sessions
│   ├── lock-manager/      # Enforces app locking logic on user devices
│   ├── notification/      # Sends push notifications to Users and Guardians
│   ├── usage-tracker/     # Logs usage data and lock/unlock attempts
│   └── analytics/         # Provides usage insights and visualizations
│
├── proto/                 # Shared gRPC (or OpenAPI) definitions for all services
│   ├── auth.proto
│   ├── session.proto
│   └── ...
│
├── shared/                # Shared libraries/utilities reused across services
│   ├── config/            # Config loading logic
│   ├── middleware/        # Tenant-aware and auth-aware middleware
│   ├── logging/           # Consistent logging setup
│   └── errors/            # Common error definitions
│
├── deployments/           # Kubernetes manifests and Helm charts
│   ├── dev/               # Development environment configs
│   ├── staging/
│   └── prod/
│
├── infra/                 # Infrastructure tools and local dev setup
│   ├── docker-compose.yml # Local service orchestration
│   ├── nginx/             
│   └── postgres/          
│
├── scripts/               
│   └── setup-local.sh
│
├── Makefile               
├── README.md             
└── go.work / package.json 

````

---

## 🚀 Key Features

- **Multitenant architecture** with session-based tenant isolation
- Fully decoupled **microservices** using gRPC or REST APIs
- Dockerized and Kubernetes-ready deployments
- **Tenant-aware middleware**, config, and logging
- Centralized **proto** definitions and **shared libraries**
- CI/CD-ready structure with support for environment-specific deployment

---

## 🧰 Tech Stack

| Component         | Technology                 |
|------------------|----------------------------|
| Backend Services | Go / Node.js (per service) |
| Transport        | gRPC or REST               |
| Messaging        | Redis Pub/Sub / NATS (optional) |
| Database         | PostgreSQL                 |
| Infrastructure   | Docker, Kubernetes         |
| Auth             | JWT + Role-Based Access    |
| Monitoring       | Prometheus + Grafana (optional) |

---

## 🛠️ Getting Started (Local Dev)

### 1. Prerequisites
- Docker + Docker Compose
- Go 1.21+ (if using Go services)
- `make` installed (or run commands manually)

### 2. Clone the Repository

```bash
git clone https://github.com/your-org/accountability-platform.git
cd accountability-platform
````

### 3. Start All Services Locally

```bash
make up
```

This will:

* Build all microservices
* Start Docker containers
* Connect services to Postgres and internal APIs

### 4. Compile Proto Files (if using gRPC)

```bash
make proto
```

---

## 🧪 Development Commands

| Command      | Description                        |
| ------------ | ---------------------------------- |
| `make up`    | Spin up all services with Docker   |
| `make down`  | Tear down all running services     |
| `make build` | Build all Docker images            |
| `make proto` | Compile `.proto` files to stubs    |
| `make test`  | Run unit tests across all services |

---

## 🧩 Contributing

1. Fork the repo
2. Create a feature branch
3. Make your changes (be tenant-safe!)
4. Submit a PR with a detailed description

---


## 🧩 MCV Structure 

                          ┌────────────────────┐
                          │   Authentication   │
                          └────────────────────┘
                                   ↓
                         ┌─────────────────────┐
                         │   API Gateway / BFF │
                         └─────────────────────┘
        ┌──────────────────────────┼──────────────────────────┐
        │                          │                          │
┌─────────────┐        ┌────────────────┐         ┌─────────────────┐
│ Session Mgr │◄──────►│  Lock Manager  │◄───────►│ Rules / Configs │
└─────────────┘        └────────────────┘         └─────────────────┘
        │                          │                          │
        ▼                          ▼                          ▼
┌─────────────┐        ┌────────────────┐         ┌─────────────────┐
│  Notif Svc  │◄──────►│ Usage Tracker  │────────►│ Analytics Engine│
└─────────────┘        └────────────────┘         └─────────────────┘
