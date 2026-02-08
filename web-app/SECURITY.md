# Security Guidelines for IPO Pilot

## 🔒 Keeping Secrets Safe

This document outlines best practices for securing sensitive information in the IPO Pilot public repository.

### ⚠️ NEVER Commit

The following should **NEVER** be committed to version control:

- `.env` files (use `.env.example` as template)
- API keys and secrets
- Database credentials
- Private keys (`.key`, `.pem`)
- Password hashes or tokens
- SSL certificates
- OAuth tokens
- Payment gateway credentials

### ✅ Do This

1. **Use `.env.example`**
   - Copy `.env.example` to `.env` locally
   - Fill in your actual values in `.env`
   - Never commit `.env` (already in `.gitignore`)

2. **Use Environment Variables**
   - For production, use Railway's environment variable UI
   - All sensitive config should be env vars, not hardcoded

3. **For Local Development**
   ```bash
   # Copy the example template
   cp web-app/.env.example web-app/.env
   
   # Edit with your local values
   nano web-app/.env
   
   # Load environment before running
   source web-app/.env
   go run ./web-app
   ```

4. **Verify Before Commit**
   ```bash
   # Check for accidentally staged secrets
   git diff --cached | grep -i "password\|secret\|key\|token\|api"
   ```

### 🚨 If You Accidentally Committed Secrets

DO NOT just delete and recommit - the secret is in git history:

1. **Rotate the credentials immediately**
2. **Use git-filter-branch or BFG Repo-Cleaner to remove from history**
3. **Notify admin if it was a shared secret**

### 🔐 Production Environment Variables (Railway.app)

Set these in Railway's Environment tab:

- `DB_TYPE=postgres`
- `DB_HOST=<railway-postgres-host>`
- `DB_PORT=5432`
- `DB_USER=postgres`
- `DB_PASSWORD=<railway-generated>`
- `DB_NAME=railway`
- `PORT=8080`
- `GIN_MODE=release`
- `JWT_SECRET=<generate-strong-random>`
- `ESEWA_MERCHANT_CODE=<your-code>`
- `KHALTI_API_KEY=<your-key>`
- etc.

### 🔑 Generating Strong Secrets

```bash
# Generate JWT secret (32 bytes = 256 bits)
openssl rand -base64 32

# Generate API key-like secret
python3 -c "import secrets; print(secrets.token_urlsafe(32))"
```

### 🛡️ Security Checklist

- [ ] `.env` file is in `.gitignore`
- [ ] `.env.example` shows structure without real values
- [ ] No secrets in code comments
- [ ] No secrets in documentation
- [ ] Environment variables used for all config
- [ ] Production passwords use strong random values
- [ ] Regular security audit of git history
- [ ] Team members know to use `.env.example`

### 📚 Related Files

- [`.env.example`](.env.example) - Template with placeholders
- [`.gitignore`](../../.gitignore) - Excludes sensitive files
- [DEVELOPER_COMPLETE_GUIDE.md](../IPO%20Pilot%20-%20Admin/DEVELOPER_COMPLETE_GUIDE.md) - Development setup

### 🆘 Questions?

If you have security concerns or accidentally commit a secret:
1. Contact the admin immediately
2. Do not attempt to cover it up
3. We'll rotate credentials and clean history

**Remember: Security is everyone's responsibility!**
