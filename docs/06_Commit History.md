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

#### `cbb530c` - feat(backend): Go Fiber 프로젝트 초기화 및 인증 시스템 구현
**Author**: haneumLee <2haneum@naver.com>  
**Branch**: feature/backend-init-auth  
**Type**: Feature Implementation

**주요 변경사항:**
- ✅ Go Fiber v2 프로젝트 구조 셋업 (Handler → Service → Repository)
- ✅ GORM PostgreSQL 연결 및 커넥션 풀 설정
- ✅ GORM 모델 구현: User, Subscription, Category, ShareGroup, ShareMember, SubscriptionShare
- ✅ OAuth 2.0 인증 플로우 구현 (Google/Apple/Naver/Kakao) + JWT 세션 관리
- ✅ JWT 미들웨어 추가 (Access 1h / Refresh 7d)
- ✅ 커스텀 에러 타입, 표준 API 응답 헬퍼, 입력값 검증 유틸 추가
- ✅ 헬스체크 엔드포인트 및 Graceful Shutdown 구현
- ✅ 모델, 유틸, 인증 서비스 단위 테스트 작성 (전체 통과)
- 📎 Refs: F-01, F-03, E1-1~E1-7, E4-1~E4-6, NFR-2.2

**Stats:**
- 26 files changed
- 3,236 insertions(+)

---

#### `4de11b0` - feat(backend): 구독 CRUD API 구현
**Author**: haneumLee <2haneum@naver.com>  
**Branch**: feature/backend-init-auth  
**Type**: Feature Implementation

**주요 변경사항:**
- ✅ SubscriptionRepository 구현 (필터링/정렬/페이지네이션 지원)
- ✅ SubscriptionService 구현 (생성/조회/수정/삭제/만족도 평가)
- ✅ SubscriptionHandler 구현 (6개 엔드포인트)
- ✅ 소유권 검증 로직 추가 (타 사용자 접근 시 403 반환)
- ✅ 금액 환산 필드 포함 응답 (monthlyAmount, annualAmount)
- ✅ 입력값 검증: 서비스명 1-100자, 금액 0~9,999,999원, 만족도 1-5점
- ✅ 구독 서비스 단위 테스트 28개 작성 (전체 통과)
- 📎 Refs: F-02, F-04, E1-1~E1-7, E2-1~E2-6, E3-1~E3-6

**Stats:**
- 5 files changed
- 1,298 insertions(+)

---

#### `737a7a2` - feat(backend): 대시보드 요약/해지 추천 API 및 시뮬레이션(해지/추가/적용) API 구현
**Author**: haneumLee <2haneum@naver.com>  
**Branch**: feature/dashboard-simulation  
**Type**: Feature Implementation

**주요 변경사항:**
- ✅ DashboardService: 월/연 총액, 활성/일시중지 카운트, 카테고리별 비중 계산
- ✅ DashboardService: 만족도 1-2점 및 고비용 저만족도 기반 해지 추천 로직
- ✅ SimulationService: 해지/추가 시뮬레이션 실시간 비용 변동 계산 (DB 미반영)
- ✅ SimulationService: 시뮬레이션 적용 시 소유권 검증 후 Soft Delete 처리
- ✅ DashboardHandler/SimulationHandler: 6개 엔드포인트 라우팅 추가
- ✅ 대시보드 서비스 단위 테스트 15개, 시뮬레이션 서비스 단위 테스트 14개 작성 (전체 통과)
- 📎 Refs: F-03, F-04, F-05, E1-1~E1-7

**Stats:**
- 7 files changed
- 1,374 insertions(+)

---

#### `c70fbbc` - feat(backend): 카테고리 CRUD 및 공유 그룹 CRUD API 구현
**Author**: haneumLee <2haneum@naver.com>  
**Branch**: feature/dashboard-simulation  
**Type**: Feature Implementation

**주요 변경사항:**
- ✅ CategoryRepository/Service/Handler: 시스템+사용자 카테고리 조회, 커스텀 카테고리 생성/수정/삭제
- ✅ 시스템 카테고리 수정/삭제 방지, 삭제 시 구독 항목 '기타'로 자동 재배치
- ✅ ShareGroupRepository/Service/Handler: 공유 그룹 CRUD 및 멤버 관리
- ✅ 그룹 생성 시 소유자 자동 추가(isOwner=true), 최소 2명 검증
- ✅ 소유권 검증 로직 적용(조회/수정/삭제)
- ✅ 그룹 삭제 시 SubscriptionShare 레코드 자동 제거
- ✅ 카테고리 서비스 단위 테스트 16개, 공유 그룹 서비스 단위 테스트 14개 작성 (전체 통과)
- 📎 Refs: F-07, F-10, E2-1~E2-6, E4-1~E4-6

**Stats:**
- 9 files changed
- 1,623 insertions(+)

---

#### `dc61fed` - feat(backend): main.go 의존성 주입, 구독 공유 분담 API, 시뮬레이션 Undo, 대시보드 개인 부담액 반영
**Author**: haneumLee <2haneum@naver.com>  
**Branch**: feature/backend-integration  
**Type**: Feature Implementation

**주요 변경사항:**
- ✅ main.go: Repository → Service → Handler 의존성 주입 및 routes.SetupRoutes() 연결
- ✅ main.go: DB AutoMigrate 및 기본 카테고리 시딩 추가
- ✅ SubscriptionShareRepository/Service/Handler: 구독-공유그룹 연결/수정/해제/조회 API (4개 엔드포인트)
- ✅ SimulationService: Undo 기능 추가 (인메모리 30초 TTL, POST /simulation/undo)
- ✅ SubscriptionRepository: Restore 메서드 추가 (soft delete 복원)
- ✅ DashboardService: 공유 분담(equal/custom_amount/custom_ratio) 개인 부담액 기준 합산
- ✅ SimulationService: 시뮬레이션 계산에 공유 분담 개인 부담액 반영
- ✅ 단위 테스트 218개 전체 통과 (SubscriptionShare 13개, Undo 4개, Dashboard 공유 4개 추가)
- 📎 Refs: F-03, F-05, F-10, E1-1~E1-7

**Stats:**
- 13 files changed
- 1,578 insertions(+), 61 deletions(-)

---

#### `e985dcd` - feat: F-09 리포트/차트 API 구현 (GET /api/v1/reports/overview)
**Author**: haneumLee <2haneum@naver.com>  
**Branch**: feature/backend-integration  
**Type**: Feature Implementation

**주요 변경사항:**
- ✅ ReportService 구현 (카테고리별 분류, 월별 추이, 평균 비용, 요약 통계)
- ✅ ReportHandler 구현 (GET /api/v1/reports/overview)
- ✅ 공유 분담 개인 부담액 기반 리포트 계산
- 📎 Refs: F-09

**Stats:**
- 6 files changed
- 618 insertions(+)

---

#### `fcbe252` - feat(backend): F-08 결제일 캘린더 API 보강 및 F-09 리포트/차트 단위 테스트 추가
**Author**: haneumLee <2haneum@naver.com>
**Branch**: feature/backend-integration
**Type**: Feature Implementation + Testing

**주요 변경사항:**
- ✅ CalendarService: GetDayDetail(일별 결제 상세), GetUpcomingPayments(향후 N일 결제 예정) 메서드 추가
- ✅ CalendarHandler: GET /api/v1/calendar/daily, GET /api/v1/calendar/upcoming 핸들러 추가
- ✅ routes.go: calendar 라우트 3개 등록 (monthly, daily, upcoming)
- ✅ CalendarService 단위 테스트 22개 작성 (전체 통과)
- ✅ ReportService 단위 테스트 17개 작성 (전체 통과)
- 📎 Refs: F-08, F-09, E1-1~E1-7

**Stats:**
- 5 files changed
- 1,621 insertions(+)

---

#### `51ebec1` - feat(frontend): 프론트엔드 초기 설정 및 전체 페이지/컴포넌트 구현
**Author**: haneumLee <2haneum@naver.com>
**Branch**: feature/frontend-init
**Type**: Feature Implementation

**주요 변경사항:**
- ✅ Next.js 14 App Router 프로젝트 수동 구성 (package.json, tsconfig.json, tailwind.config.ts)
- ✅ AppLayout: 반응형 사이드바 네비게이션, 모바일 햄버거 메뉴, 로그아웃 모달
- ✅ AuthContext: JWT 기반 인증 상태 관리 (login/logout/refresh)
- ✅ UI 컴포넌트: Button, Modal, Toast, LoadingSpinner
- ✅ 대시보드, 구독 관리, 시뮬레이션, 공유 그룹 페이지 구현
- ✅ React Query 기반 커스텀 훅 (useSubscriptions, useDashboard, useSimulation, useShareGroups 등)
- ✅ 단위 테스트 52개 작성 (전체 통과)
- 📎 Refs: F-01~F-05, F-10

**Stats:**
- 60 files changed
- 6,839 insertions(+)

---

#### `8d71bc8` - feat(frontend): F-08 캘린더, F-09 리포트, F-11 설정 페이지 구현
**Author**: haneumLee <2haneum@naver.com>
**Branch**: feature/frontend-remaining-pages
**Type**: Feature Implementation + Testing

**주요 변경사항:**
- ✅ F-08 결제일 캘린더: 월별 캘린더 그리드, 일별 상세 모달, 다가오는 결제 목록
- ✅ F-09 리포트/차트: CSS-only 카테고리 도넛 차트, 월별 추이 바 차트, 비용 요약, 구독 요약
- ✅ F-11 설정 페이지: 프로필 정보 표시, 카테고리 CRUD 관리(추가/수정/삭제)
- ✅ AppLayout 네비게이션에 캘린더/리포트 메뉴 추가
- ✅ img 태그 → Next.js Image 컴포넌트 변환
- ✅ next.config.js에 OAuth 프로바이더 아바타 이미지 도메인 설정
- ✅ types/index.ts에 Calendar/Report 타입 추가
- ✅ useCalendar, useReports, useCategories(CRUD) 훅 추가
- ✅ 전체 69개 단위 테스트 통과 (15 suites)
- 📎 Refs: F-07, F-08, F-09, F-11

**Stats:**
- 33 files changed
- 1,585 insertions(+), 5 deletions(-)

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
- Total Commits: 11 (e205c09, dad6813, cbb530c, 4de11b0, 737a7a2, c70fbbc, dc61fed, e985dcd, fcbe252, 51ebec1, 8d71bc8)
- Contributors: 1
- Branches: 7 (main, dev, feature/backend-init-auth, feature/dashboard-simulation, feature/backend-integration, feature/frontend-init, feature/frontend-remaining-pages)
- Tags: 0

### 브랜치별 커밋 수
- main: 1
- dev: 1
- feature/backend-init-auth: 2
- feature/dashboard-simulation: 2
- feature/backend-integration: 3
- feature/frontend-init: 1
- feature/frontend-remaining-pages: 1

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
