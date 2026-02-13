# Git Branch Strategy for SubKeep

## Overview

SubKeep follows a **Git Flow** branching strategy with three main branches and feature-based development.

## Branch Structure

```
main (production)
  │
  └─── dev (development)
         │
         ├─── feature/user-auth
         ├─── feature/subscription-crud
         ├─── fix/timezone-bug
         └─── refactor/api-optimization
```

## Branches

### 1. `main` (Production)
- **Purpose**: Production-ready code
- **Protection**: Protected, requires PR approval
- **Deployment**: Auto-deploys to OracleVM production server
- **Merge from**: `dev` only
- **Versioning**: Tagged with semantic versions (v1.0.0, v1.1.0, etc.)

### 2. `dev` (Development)
- **Purpose**: Integration branch for features
- **Protection**: Protected, requires PR approval
- **Testing**: All features tested here before merging to main
- **Merge from**: feature/*, fix/*, refactor/*
- **Merge to**: `main` for releases

### 3. Feature Branches
- **Naming**: `feature/<description>`, `fix/<description>`, `refactor/<description>`
- **Created from**: `dev`
- **Merged to**: `dev`
- **Lifecycle**: Short-lived, deleted after merge

## Branch Naming Convention

```
feature/<short-description>    # New features
fix/<bug-description>          # Bug fixes
refactor/<what-refactored>     # Code refactoring
docs/<what-documented>         # Documentation updates
test/<what-tested>             # Test additions
chore/<task-description>       # Maintenance tasks
```

**Examples:**
```
feature/user-authentication
feature/subscription-tracking
fix/monthly-calculation-error
fix/timezone-mismatch
refactor/database-queries
docs/api-specification
test/subscription-service
chore/update-dependencies
```

## Workflow

### Starting New Work

```bash
# Update dev branch
git checkout dev
git pull origin dev

# Create feature branch
git checkout -b feature/user-authentication

# Work on feature
git add .
git commit -m "feat(auth): implement JWT authentication"
git push origin feature/user-authentication

# Create Pull Request: feature/user-authentication → dev
```

### Merging to Dev

```bash
# PR Review required
# Tests must pass
# No merge conflicts
# At least 1 approval (if team exists)

# After merge, delete feature branch
git branch -d feature/user-authentication
git push origin --delete feature/user-authentication
```

### Releasing to Production

```bash
# Create release PR: dev → main
# Full testing & QA
# Update version and CHANGELOG
# Merge to main
git checkout main
git pull origin main

# Tag release
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

## Commit Message Convention

We follow **Conventional Commits** format:

```
<type>(<scope>): <subject>

[optional body]

[optional footer]
```

### Types
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style (formatting, semicolons, etc.)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Maintenance tasks
- `build`: Build system changes
- `ci`: CI/CD changes
- `revert`: Revert previous commit

### Examples
```bash
git commit -m "feat(api): add user registration endpoint"
git commit -m "fix(ui): resolve button alignment on mobile"
git commit -m "docs: update installation instructions"
git commit -m "refactor(db): optimize subscription queries"
git commit -m "test(auth): add JWT validation tests"
```

## Pull Request Guidelines

### Title Format
```
[Type] Brief description

Examples:
[Feature] Add user authentication
[Fix] Resolve timezone calculation bug
[Refactor] Optimize database queries
```

### Description Template
```markdown
## What
Brief description of changes

## Why
Reason for this change

## How
Technical approach

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests pass
- [ ] Manual testing completed

## Screenshots (if UI changes)

## Related Issues
Closes #123
```

## Branch Protection Rules

### `main` branch
- ✅ Require pull request before merging
- ✅ Require approvals: 1+ (if team)
- ✅ Require status checks to pass
- ✅ Require up-to-date before merge
- ✅ No force pushes
- ✅ No deletions

### `dev` branch
- ✅ Require pull request before merging
- ✅ Require status checks to pass
- ✅ No force pushes
- ✅ No deletions

## Local Setup

```bash
# Clone repository
git clone <repository-url>
cd subkeep

# Setup git hooks
git config core.hooksPath .githooks

# Verify hooks are working
.githooks/pre-commit
```

## Deployment Mapping

| Branch | Environment | Server | Auto-Deploy |
|--------|-------------|--------|-------------|
| `main` | Production  | OracleVM | Yes (PM2) |
| `dev`  | Development | Local    | No |
| `feature/*` | Local | Local | No |

## Emergency Hotfix

```bash
# Create hotfix from main
git checkout main
git checkout -b fix/critical-security-issue

# Fix the issue
git commit -m "fix(security): patch XSS vulnerability"

# Merge to main immediately
# PR review expedited
git checkout main
git merge fix/critical-security-issue
git tag -a v1.0.1 -m "Hotfix: Security patch"
git push origin main --tags

# Backport to dev
git checkout dev
git merge main
git push origin dev

# Delete hotfix branch
git branch -d fix/critical-security-issue
```

## Best Practices

1. **Always pull before creating new branch**
2. **Keep feature branches small and focused**
3. **Commit early, commit often**
4. **Write descriptive commit messages**
5. **Test before pushing**
6. **Keep `dev` and `main` clean**
7. **Delete merged branches**
8. **Tag all production releases**
9. **Never commit sensitive data**
10. **Use `.env` for environment variables**

## Git Hooks

Automated checks run before commits:

- ✅ Go code formatting (gofmt)
- ✅ Go linting (golangci-lint)
- ✅ TypeScript linting (ESLint)
- ✅ Code formatting (Prettier)
- ✅ Unit tests
- ✅ Commit message validation

## Troubleshooting

### Merge Conflicts
```bash
# Update your branch with latest dev
git checkout feature/my-feature
git fetch origin
git merge origin/dev

# Resolve conflicts
# Mark as resolved
git add .
git commit
```

### Failed Pre-commit Hook
```bash
# Fix the issues reported
# Bypass only in emergencies
git commit --no-verify
```

### Accidentally Committed to Wrong Branch
```bash
# Move commit to correct branch
git log  # Find commit hash
git checkout correct-branch
git cherry-pick <commit-hash>
git checkout wrong-branch
git reset HEAD~1 --hard
```
