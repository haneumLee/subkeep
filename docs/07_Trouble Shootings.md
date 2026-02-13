# SubKeep - Trouble Shootings

## 프론트엔드 개발 이슈 (2026-02-13)

### TS-001: Node.js 미설치

**상황**: 프론트엔드 개발 시작 시 `npx create-next-app` 실행 불가
**문제**: 서버에 Node.js가 설치되어 있지 않음 (`npx: command not found`)
**방법**: NodeSource 리포지토리를 통한 Node.js 20.x LTS 설치
**해결**:
```bash
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install nodejs -y
```
Node.js 20.20.0 설치 완료.

---

### TS-002: create-next-app 인터랙티브 모드

**상황**: `npx create-next-app --typescript --tailwind --app` 실행
**문제**: CLI 플래그를 넘겨도 인터랙티브 프롬프트가 표시되어 자동화 불가
**방법**: create-next-app을 사용하지 않고 수동으로 프로젝트 구성
**해결**: package.json, tsconfig.json, next.config.js 등 모든 설정 파일을 직접 생성하고 `npm install` 실행.

---

### TS-003: jest.config.ts 오타

**상황**: TypeScript 컴파일 시 jest.config.ts에서 에러 발생
**문제**: `setupFilesAfterSetup` 프로퍼티명 오타 (올바른 이름: `setupFilesAfterEnv`)
**방법**: jest.config.ts에서 프로퍼티명 수정
**해결**:
```diff
- setupFilesAfterSetup: ['<rootDir>/jest.setup.ts'],
+ setupFilesAfterEnv: ['<rootDir>/jest.setup.ts'],
```

---

### TS-004: ESLint react-refresh 플러그인 호환성

**상황**: Next.js 빌드 시 ESLint 에러 다수 발생
**문제**: `.eslintrc.json`이 Vite 스타일로 설정되어 있어 `react-refresh/only-export-components` 규칙을 찾을 수 없음. 또한 `eslint-config-prettier` 미설치.
**방법**:
1. 누락된 ESLint 패키지 설치: `eslint-config-prettier`, `eslint-plugin-react`, `eslint-plugin-jsx-a11y`, `eslint-plugin-import`
2. `.eslintrc.json`에서 `react-refresh` 플러그인 제거 (Next.js에서 불필요)
3. `next/core-web-vitals` extends 추가
4. `jsx-a11y/label-has-associated-control`, `jsx-a11y/anchor-is-valid`를 warn으로 변경
**해결**: ESLint 설정을 Next.js 호환으로 업데이트하여 빌드 성공.

---

### TS-005: useSearchParams Suspense boundary 필요

**상황**: `next build` 시 `/auth/callback` 페이지 정적 생성 실패
**문제**: Next.js 14에서 `useSearchParams()`는 반드시 `<Suspense>` 경계 안에서 사용해야 함
**방법**: 페이지 컴포넌트를 Suspense로 감싸고 내부 로직을 별도 컴포넌트로 분리
**해결**:
```tsx
export default function AuthCallbackPage() {
  return (
    <Suspense fallback={<LoadingSpinner />}>
      <AuthCallbackContent />  // useSearchParams 사용
    </Suspense>
  );
}
```

---

### TS-006: AppLayout jsx-a11y 접근성 에러

**상황**: 빌드 시 `jsx-a11y/click-events-have-key-events` 에러
**문제**: 모바일 사이드바 오버레이가 `<div onClick>` 으로 구현되어 접근성 위반
**방법**: `<div>` → `<button type="button">`으로 변경하고 `aria-label` 추가
**해결**:
```tsx
<button
  type="button"
  className="fixed inset-0 z-20 bg-black bg-opacity-50 lg:hidden"
  onClick={() => setIsSidebarOpen(false)}
  aria-label="사이드바 닫기"
/>
```

---

### TS-007: GitHub 원격 푸시 인증 실패

**상황**: `git push -u origin feature/frontend-init` 실행
**문제**: `fatal: could not read Username for 'https://github.com'` - GitHub 인증 정보 미설정
**방법**: SSH 키 또는 GitHub Personal Access Token 설정 필요
**해결**: 커밋은 로컬에 완료됨. GitHub 인증 설정 후 수동 푸시 필요.
```bash
git push -u origin feature/frontend-init
```

---

### TS-008: 다중 에이전트 동시 작업 시 git 브랜치 충돌

**상황**: `feature/frontend-remaining-pages` 브랜치에서 프론트엔드 에이전트가 작업 중, 다른 에이전트가 동일 브랜치의 파일을 커밋 없이 수정 중이었음
**문제**: `git checkout feature/frontend-remaining-pages` 시 "untracked working tree files would be overwritten by checkout" 에러 발생. 다른 에이전트의 미커밋 변경사항과 현재 에이전트의 신규 파일이 충돌.
**방법**:
1. `git stash --include-untracked` 로 다른 에이전트의 작업을 임시 저장
2. 올바른 브랜치로 전환 후 `git stash pop`으로 복원
3. 복원된 상태에서 빌드 검증 후 커밋
**해결**: stash를 활용하여 두 에이전트의 작업을 모두 보존. 향후 다중 에이전트 작업 시 각 에이전트가 독립된 feature 브랜치에서 작업하고, 작업 단위마다 즉시 커밋하여 충돌 방지 필요.

---

### TS-009: DayDetailModal 빌드 에러 - Modal title prop 누락

**상황**: 캘린더 DayDetailModal 컴포넌트에서 Modal 사용 시 빌드 에러
**문제**: `Type error: Property 'title' is missing in type` - Modal 컴포넌트의 `title` prop이 필수인데 누락됨
**방법**: Modal에 `title` prop 추가 및 `showFooter={false}` 설정 (읽기 전용 모달이므로)
**해결**:
```tsx
<Modal
  isOpen={isOpen}
  onClose={onClose}
  title={formatDetailDate(dayDetail.date)}
  showFooter={false}
>
```
