# SubKeep - Commit History

## 📋 프로젝트 커밋 이력

이 문서는 SubKeep 프로젝트의 주요 커밋 이력을 추적합니다.

---

## 2026년 2월

### 2026-02-13 (Thu)

#### `e205c09` - chore: 완전한 인프라를 갖춘 초기 프로젝트 설정
**Author**: haneumLee <2haneum@naver.com>  
**Branch**: main, dev  
**Type**: Infrastructure Setup

**주요 변경사항:**
- ✅ Git workflow 설정 (main, dev, feature 브랜치 전략)
- ✅ Husky를 통한 pre-commit hooks 구성
- ✅ Conventional Commits 검증 추가
- ✅ PostgreSQL 및 Redis용 Docker Compose 생성
- ✅ golang-migrate를 사용한 데이터베이스 마이그레이션 설정
- ✅ GitHub Actions CI/CD 워크플로우 구성
- ✅ Go 및 Next.js용 포괄적인 .gitignore 추가
- ✅ 모든 서비스에 대한 .env.example 템플릿 생성
- ✅ 배포를 위한 PM2 ecosystem 설정
- ✅ 린팅 설정 추가 (golangci-lint, ESLint, Prettier)
- ✅ 포괄적인 README 및 문서 작성

**생성된 파일 (35개):**
```
.commitlintrc.json
.githooks/commit-msg
.githooks/pre-commit
.github/workflows/SECRETS.md
.github/workflows/ci.yml
.github/workflows/deploy.yml
.gitignore
.husky/_/husky.sh
.husky/commit-msg
.husky/pre-commit
README.md
backend/.env.example
backend/.golangci.yml
backend/Dockerfile
backend/migrations/000001_create_users_table.down.sql
backend/migrations/000001_create_users_table.up.sql
backend/migrations/000001_init_schema.down.sql
backend/migrations/000001_init_schema.up.sql
backend/migrations/README.md
backend/scripts/migrate.sh
backend/seeds/dev/001_sample_data.sql
backend/seeds/prod/001_initial_data.sql
docker-compose.prod.yml
docker-compose.yml
docker/.env.example
docker/postgres/init/01-init.sql
docs/06_Commit History.md
docs/07_Trouble Shootings.md
docs/GIT_WORKFLOW.md
ecosystem.config.js
frontend/.env.example
frontend/.eslintrc.json
frontend/.prettierrc.json
frontend/Dockerfile
frontend/nginx.conf
```

**기술 스택:**
- Backend: Go 1.22 + Fiber + GORM + PostgreSQL 15
- Frontend: Next.js 14 + TypeScript + React
- Cache: Redis 7
- Deployment: Docker, PM2, OracleVM
- CI/CD: GitHub Actions

**Stats:**
- 35 files changed
- 3,604 insertions(+)

---

## Commit Convention

이 프로젝트는 [Conventional Commits](https://www.conventionalcommits.org/) 규칙을 따릅니다.

### Commit Types
- `feat`: 새로운 기능 추가
- `fix`: 버그 수정
- `docs`: 문서 변경
- `style`: 코드 포맷팅 (기능 변경 없음)
- `refactor`: 코드 리팩토링
- `test`: 테스트 추가/수정
- `chore`: 빌드, 설정 변경
- `perf`: 성능 개선
- `ci`: CI/CD 설정

### Commit Message Format
```
<type>(<scope>): <subject>

[optional body]

[optional footer]
```

**예시:**
```
feat(auth): JWT 토큰 갱신 로직 구현
fix(api): 구독 계산 버그 수정
docs: README 설치 가이드 업데이트
```

---

## Branch Strategy

### Main Branches
- `main`: 프로덕션 배포 브랜치 (protected)
- `dev`: 개발 통합 브랜치

### Supporting Branches
- `feature/*`: 기능 개발 브랜치
- `hotfix/*`: 긴급 수정 브랜치
- `release/*`: 릴리스 준비 브랜치 (선택적)

---

## Release History

### v0.1.0 (계획 중)
- 초기 MVP 릴리스
- 기본 CRUD 기능
- 사용자 인증
- 구독 관리 기능

---

## Statistics

### 전체 통계
- Total Commits: 1
- Contributors: 1
- Branches: 2 (main, dev)
- Tags: 0

### 브랜치별 커밋 수
- main: 1
- dev: 1

---

## Contributors

- **haneumLee** (2haneum@naver.com) - 1 commits

---

## 업데이트 방법

이 문서는 매 커밋 후 수동으로 업데이트하거나, Agent를 통해 자동으로 업데이트됩니다.

```bash
# 최근 10개 커밋 확인
git log --oneline -10

# 특정 브랜치의 커밋 확인
git log dev --oneline

# 통계 확인
git shortlog -sn --all
```

---

**Last Updated**: 2026-02-13  
**Next Review**: 매 주요 기능 커밋 후
