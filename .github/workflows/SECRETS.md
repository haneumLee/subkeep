# GitHub Actions Secrets

Configure the following secrets in your GitHub repository:  
**Settings → Secrets and variables → Actions → New repository secret**

## Required Secrets

### Deployment
```
DEPLOY_SSH_KEY          # Private SSH key for deployment server
DEPLOY_USER             # SSH username (e.g., deploy)
PRODUCTION_HOST         # OracleVM IP or domain (e.g., 192.168.1.100)
DEPLOY_PATH             # Deployment directory (e.g., /home/deploy/subkeep)
```

### Database
```
DATABASE_URL            # PostgreSQL connection string
                        # Example: postgres://user:pass@host:5432/dbname?sslmode=disable
```

### Optional (if using external services)
```
CODECOV_TOKEN          # For code coverage reports
SENTRY_DSN             # For error monitoring
GITHUB_TOKEN           # Automatically provided by GitHub Actions
```

## Setup Instructions

### 1. Generate SSH Key Pair

On your local machine:
```bash
ssh-keygen -t ed25519 -C "github-actions-deploy" -f ~/.ssh/subkeep_deploy
```

### 2. Add Public Key to OracleVM Server

```bash
# Copy public key
cat ~/.ssh/subkeep_deploy.pub

# On OracleVM server
mkdir -p ~/.ssh
echo "<PUBLIC_KEY>" >> ~/.ssh/authorized_keys
chmod 700 ~/.ssh
chmod 600 ~/.ssh/authorized_keys
```

### 3. Add Private Key to GitHub Secrets

```bash
# Get private key content
cat ~/.ssh/subkeep_deploy

# Add to GitHub:
# Settings → Secrets → New secret
# Name: DEPLOY_SSH_KEY
# Value: <paste private key>
```

### 4. Test SSH Connection

```bash
ssh -i ~/.ssh/subkeep_deploy deploy@<PRODUCTION_HOST>
```

### 5. Setup Deployment Directory on OracleVM

```bash
# On OracleVM server
mkdir -p /home/deploy/subkeep/{backend/bin,logs}
```

## Environment Variables Setup

Create `.env` file on production server:
```bash
# On OracleVM
cd /home/deploy/subkeep
cp backend/.env.example backend/.env

# Edit with production values
nano backend/.env
```

## Verify Deployment Setup

```bash
# Check PM2 is installed
pm2 -v

# Check directory permissions
ls -la /home/deploy/subkeep

# Check SSH access
ssh -i ~/.ssh/subkeep_deploy deploy@<HOST> "echo 'SSH OK'"
```

## Security Notes

- ✅ Never commit private keys
- ✅ Use separate SSH keys for CI/CD
- ✅ Limit SSH key permissions (read-only where possible)
- ✅ Rotate keys regularly
- ✅ Use strong database passwords
- ✅ Enable firewall on OracleVM
