# SubKeep

**Smart Subscription Management System**

SubKeep은 단순한 구독 관리를 넘어, **구독 유지/해지 판단을 돕는** 지능형 시스템입니다.

## 📋 Table of Contents

- [Features](#-features)
- [Tech Stack](#-tech-stack)
- [Architecture](#-architecture)
- [Getting Started](#-getting-started)
- [Development](#-development)
- [Deployment](#-deployment)
- [API Documentation](#-api-documentation)
- [Contributing](#-contributing)
- [License](#-license)

## ✨ Features

### Core Capabilities

- **🧮 정확한 월 환산 계산**: 연간/분기별 구독도 월 단위로 정확히 환산
- **📊 이용률 기반 판단**: 실제 사용 패턴을 분석하여 유지/해지 추천
- **👥 공유 최적화**: 가족/친구와의 공유 시 1인당 실비용 계산
- **📈 트렌드 분석**: 구독 비용 추이 및 카테고리별 지출 분석
- **🔔 스마트 알림**: 갱신일, 무료 체험 종료, 가격 인상 알림
- **🤝 협업 자동화**: YAML 보고서 기반 팀 협업

### Technical Excellence

- **🔐 JWT 인증**: Stateless 인증 시스템
- **📝 자동 문서화**: OpenAPI/Swagger 자동 생성
- **✅ 품질 보증**: Pre-commit hooks, 자동 테스트, Linting
- **🐳 컨테이너화**: Docker/Docker Compose 완전 지원
- **🚀 CI/CD**: GitHub Actions 자동 배포
- **📊 모니터링**: 구조화된 로깅 및 에러 추적

## 🛠 Tech Stack

### Backend
- **Language**: Go 1.22
- **Framework**: Fiber v2 (HTTP)
- **ORM**: GORM
- **Database**: PostgreSQL 15
- **Auth**: golang-jwt/jwt v5
- **Validation**: go-playground/validator v10
- **Migration**: golang-migrate
- **Config**: godotenv

### Frontend
- **Language**: TypeScript 5.7
- **Framework**: Next.js 14 (App Router) + React 18
- **Styling**: Tailwind CSS 3.4
- **State/Data Fetching**: TanStack React Query 5
- **Form**: React Hook Form + Zod (validation)
- **HTTP Client**: Axios
- **Utilities**: date-fns, clsx, tailwind-merge
- **Testing**: Jest 29, Testing Library

### Infrastructure
- **Containerization**: Docker, Docker Compose
- **Process Manager**: PM2
- **Web Server**: Nginx
- **CI/CD**: GitHub Actions
- **Version Control**: Git (Git Flow)

### Development Tools
- **Backend Linting**: golangci-lint
- **Frontend Linting**: ESLint, Prettier
- **Git Hooks**: Husky, commitlint
- **API Docs**: Swagger/OpenAPI

## 🏗 Architecture

```
┌─────────────┐
│   Client    │
│  (Browser)  │
└──────┬──────┘
       │
       │ HTTPS
       ▼
┌─────────────┐
│   Nginx     │  Reverse Proxy
│ (Port 80)   │
└──────┬──────┘
       │
       ├─────────────┐
       │             │
       ▼             ▼
┌──────────┐   ┌──────────┐
│ Frontend │   │ Backend  │
│  React   │   │   Go     │
│ (Port    │   │ (Port    │
│  3000)   │   │  8080)   │
└──────────┘   └────┬─────┘
                    │
       ┌────────────┼────────────┐
       │            │            │
       ▼            ▼            ▼
┌──────────┐ ┌──────────┐ ┌──────────┐
│PostgreSQL│ │  Redis   │ │  Files   │
│(Port     │ │ (Port    │ │          │
│ 5432)    │ │  6379)   │ │          │
└──────────┘ └──────────┘ └──────────┘
```

### Directory Structure

```
subkeep/
├── backend/              # Go backend
│   ├── cmd/              # Application entrypoints
│   ├── internal/         # Private application code
│   ├── pkg/              # Public libraries
│   ├── migrations/       # Database migrations
│   ├── seeds/            # Seed data
│   ├── scripts/          # Utility scripts
│   └── tests/            # Tests
├── frontend/             # React frontend
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── hooks/
│   │   ├── services/
│   │   └── utils/
│   └── public/
├── docker/               # Docker configurations
├── docs/                 # Documentation
├── .github/              # GitHub Actions workflows
├── .claude/              # Claude AI agents
└── ecosystem.config.js   # PM2 configuration
```

## 🚀 Getting Started

### Prerequisites

- **Go** 1.22+
- **Node.js** 20+
- **PostgreSQL** 15+
- **Redis** 7+
- **Docker** & **Docker Compose** (optional but recommended)

### Quick Start with Docker

```bash
# 1. Clone repository
git clone <repository-url>
cd subkeep

# 2. Copy environment files
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
cp docker/.env.example docker/.env

# 3. Edit environment variables
nano docker/.env

# 4. Start all services
docker compose up -d

# 5. View logs
docker compose logs -f

# 6. Access application
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
# Swagger Docs: http://localhost:8080/swagger
```

### Manual Installation

#### Backend Setup

```bash
cd backend

# Install dependencies
go mod download

# Copy environment file
cp .env.example .env
nano .env

# Run migrations
./scripts/migrate.sh up

# Load seed data (development)
./scripts/migrate.sh seed dev

# Run backend
go run cmd/server/main.go

# Or build and run
go build -o bin/server cmd/server/main.go
./bin/server
```

#### Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Copy environment file
cp .env.example .env
nano .env

# Run development server
npm run dev

# Build for production
npm run build
```

## 💻 Development

### Setup Git Hooks

```bash
# Configure git hooks path
git config core.hooksPath .githooks

# Verify hooks are executable
ls -la .githooks/
```

### Branch Strategy

We follow **Git Flow**:

- `main`: Production-ready code
- `dev`: Development integration branch
- `feature/*`: New features
- `fix/*`: Bug fixes
- `refactor/*`: Code refactoring

See [docs/GIT_WORKFLOW.md](docs/GIT_WORKFLOW.md) for details.

### Commit Convention

We use **Conventional Commits**:

```bash
git commit -m "feat(api): add user authentication endpoint"
git commit -m "fix(ui): resolve button alignment issue"
git commit -m "docs: update installation guide"
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`, `build`, `ci`, `revert`

### Running Tests

#### Backend
```bash
cd backend
go test -v -race ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

#### Frontend
```bash
cd frontend
npm test
npm test -- --coverage
```

### Linting

#### Backend
```bash
cd backend
gofmt -l .
go vet ./...
golangci-lint run
```

#### Frontend
```bash
cd frontend
npm run lint
npm run format
npx tsc --noEmit
```

### Database Migrations

```bash
cd backend

# Create new migration
./scripts/migrate.sh create add_subscription_table

# Apply migrations
./scripts/migrate.sh up

# Rollback last migration
./scripts/migrate.sh down 1

# Check current version
./scripts/migrate.sh version
```

## 🚢 Deployment

### Production Deployment to OracleVM

#### 1. Setup Server

```bash
# On OracleVM server
sudo apt update
sudo apt install -y git docker.io docker-compose postgresql-client redis-tools

# Install PM2
npm install -g pm2

# Create deployment directory
mkdir -p /home/deploy/subkeep
```

#### 2. Configure SSH Access

```bash
# Generate SSH key for GitHub Actions
ssh-keygen -t ed25519 -C "github-actions"

# Add public key to server
cat ~/.ssh/id_ed25519.pub >> ~/.ssh/authorized_keys

# Add private key to GitHub Secrets
# Settings → Secrets → DEPLOY_SSH_KEY
```

#### 3. Setup Environment Variables

```bash
# On OracleVM
cd /home/deploy/subkeep
cp backend/.env.example backend/.env

# Edit production values
nano backend/.env
```

#### 4. Deploy

```bash
# Using PM2 (recommended)
pm2 start ecosystem.config.js --env production
pm2 save
pm2 startup

# Or using Docker
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

#### 5. Setup CI/CD

See [.github/workflows/SECRETS.md](.github/workflows/SECRETS.md) for GitHub Actions setup.

### Monitoring

```bash
# PM2 monitoring
pm2 monit
pm2 logs subkeep-backend

# Docker monitoring
docker compose logs -f
docker stats
```

## 📚 API Documentation

### Swagger UI

Visit: `http://localhost:8080/swagger`

### API Endpoints

#### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh token

#### Subscriptions
- `GET /api/v1/subscriptions` - List subscriptions
- `POST /api/v1/subscriptions` - Create subscription
- `GET /api/v1/subscriptions/:id` - Get subscription
- `PUT /api/v1/subscriptions/:id` - Update subscription
- `DELETE /api/v1/subscriptions/:id` - Delete subscription

See full API spec: [docs/09_API Specification.md](docs/09_API%20Specification.md)

## 🤝 Contributing

### Workflow

1. Fork the repository
2. Create feature branch from `dev`
3. Make your changes
4. Write/update tests
5. Ensure all tests pass
6. Submit Pull Request to `dev`

### Code Review

- At least 1 approval required
- All CI checks must pass
- No merge conflicts
- Follows coding standards

### Development Guidelines

- **Code Style**: Follow language conventions
- **Testing**: Maintain >80% coverage
- **Documentation**: Update docs for new features
- **Commits**: Use Conventional Commits
- **Security**: Never commit secrets

## 📝 Documentation

- [Business Requirements](docs/01_Business%20Requirement%20Document.md)
- [Functional Requirements](docs/02_Functional%20Requirement%20Specification.md)
- [User Stories](docs/03_User%20Stories.md)
- [System Architecture](docs/08_System%20Architecture%20Design.md)
- [API Specification](docs/09_API%20Specification.md)
- [Database Design](docs/10_Entity%20Relationship%20Diagram%20and%20Database%20Design%20Document.md)
- [Git Workflow](docs/GIT_WORKFLOW.md)

## 📄 License

[MIT License](LICENSE)

## 👥 Team

- **Developer**: [Your Name]
- **Deployment**: OracleVM Server

## 🔗 Links

- **Repository**: <repository-url>
- **Production**: <production-url>
- **Documentation**: [docs/](docs/)

## ⚡ Quick Commands

```bash
# Development
npm run dev          # Start frontend dev server
go run cmd/server    # Start backend dev server
docker compose up    # Start all services

# Testing
go test ./...        # Run backend tests
npm test             # Run frontend tests

# Deployment
pm2 start ecosystem.config.js --env production  # Deploy with PM2
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d  # Deploy with Docker

# Database
./backend/scripts/migrate.sh up    # Run migrations
./backend/scripts/migrate.sh down  # Rollback migration

# Git
git checkout -b feature/my-feature  # Create feature branch
git commit -m "feat(scope): message"  # Commit with convention
```

---

**Built with ❤️ by SubKeep Team**
