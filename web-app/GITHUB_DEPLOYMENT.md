# 🌐 Deploy to Web - Access via URL

## ✅ Yes! You can host this and access via URL

Your web app can be deployed and accessed from anywhere using a public URL like:
- `https://ipo-pilot.railway.app`
- `https://ipo-pilot.onrender.com`
- `https://ipo-pilot.fly.dev`

---

## 🚀 Method 1: Railway (Recommended - Easiest!)

**Free tier available, takes 3 minutes**

### Step 1: Push to GitHub

```bash
# Initialize git (if not already)
cd /workspaces/IPO-PILOT/web-app
git init

# Add all files
git add .

# Commit
git commit -m "IPO Pilot Web Platform"

# Create GitHub repo and push
gh repo create ipo-pilot-web --public --source=. --push
# OR manually: create repo on GitHub.com, then:
git remote add origin https://github.com/dipudai/ipo-pilot-web.git
git branch -M main
git push -u origin main
```

### Step 2: Deploy on Railway

1. **Go to:** https://railway.app
2. **Click:** "Start a New Project"
3. **Select:** "Deploy from GitHub repo"
4. **Login with GitHub** (if needed)
5. **Select:** `dipudai/ipo-pilot-web`
6. **Railway auto-detects Go app!**
7. **Click:** "Deploy Now"

**That's it!** In 2-3 minutes:
- ✅ Live URL: `https://ipo-pilot-production.up.railway.app`
- ✅ Auto-deploys on every GitHub push
- ✅ Free $5 credit monthly

### Step 3: Access Your Site

Railway will show your URL. Click it or visit:
```
https://your-app-name.up.railway.app
```

**Admin login:**
- Email: `admin@ipopilot.com`
- Password: `admin123` (change immediately!)

---

## 🚀 Method 2: Render (100% Free Tier)

**Completely free, no credit card needed**

### Step 1: Push to GitHub (same as above)

### Step 2: Deploy on Render

1. **Go to:** https://render.com
2. **Sign up** with GitHub account
3. **Click:** "New +" → "Web Service"
4. **Connect GitHub repo:** `dipudai/ipo-pilot-web`
5. **Configure:**
   - **Name:** `ipo-pilot`
   - **Environment:** `Go`
   - **Build Command:** `go build -o app`
   - **Start Command:** `./app`
6. **Click:** "Create Web Service"

**Live in 5 minutes:**
- URL: `https://ipo-pilot.onrender.com`
- Free SSL certificate
- Auto-deploy on push

---

## 🚀 Method 3: Fly.io (Great for Go Apps)

**Free tier: 3 VMs, 3GB storage**

### Quick Deploy

```bash
# Install Fly CLI
curl -L https://fly.io/install.sh | sh

# Login
fly auth login

# Navigate to app
cd /workspaces/IPO-PILOT/web-app

# Deploy (one command!)
fly launch

# Answer prompts:
# App name: ipo-pilot
# Region: Choose closest to you
# Database: No (we use SQLite)
# Deploy now: Yes

# Live at: https://ipo-pilot.fly.dev
```

**Auto-deploy on GitHub push:**
```bash
# Set up GitHub Actions
fly deploy --remote-only

# Now every push to GitHub auto-deploys!
```

---

## 🚀 Method 4: Heroku (Classic Choice)

**Free tier available with credit card verification**

```bash
# Install Heroku CLI
curl https://cli-assets.heroku.com/install.sh | sh

# Login
heroku login

# Create app
cd /workspaces/IPO-PILOT/web-app
heroku create ipo-pilot

# Deploy
git push heroku main

# Open in browser
heroku open

# Live at: https://ipo-pilot.herokuapp.com
```

**Connect to GitHub for auto-deploy:**
1. Heroku Dashboard → Your App
2. Deploy tab → Deployment method → GitHub
3. Connect repository
4. Enable "Automatic Deploys"

---

## 🚀 Method 5: Vercel (Serverless)

**Free tier, unlimited deployments**

```bash
# Install Vercel CLI
npm i -g vercel

# Navigate to app
cd /workspaces/IPO-PILOT/web-app

# Deploy
vercel

# Follow prompts:
# Set up and deploy: Y
# Link to existing project: N
# Project name: ipo-pilot
# Directory: ./

# Live at: https://ipo-pilot.vercel.app
```

---

## 📝 Required: Create Procfile (for some platforms)

Create file: `Procfile` (no extension)

```bash
cd /workspaces/IPO-PILOT/web-app
cat > Procfile << 'EOF'
web: ./app
EOF
```

---

## ⚙️ Environment Variables (Production)

For production deployment, set these on your hosting platform:

```bash
# Database (use PostgreSQL in production)
DATABASE_URL=postgres://user:pass@host:5432/dbname

# JWT Secret (change this!)
JWT_SECRET=your-super-secret-key-change-this-in-production

# Server Port (usually auto-set by platform)
PORT=8080

# App Environment
GIN_MODE=release
```

**Setting on Railway:**
```
Dashboard → Variables → Add:
JWT_SECRET=your-random-secret-key
GIN_MODE=release
```

**Setting on Render:**
```
Dashboard → Environment → Add:
JWT_SECRET=your-random-secret-key
GIN_MODE=release
```

**Setting on Heroku:**
```bash
heroku config:set JWT_SECRET=your-random-secret-key
heroku config:set GIN_MODE=release
```

---

## 🗄️ Database Options (Production)

**SQLite (Default) - Good for <100 users:**
- ✅ No setup needed
- ✅ Works immediately
- ⚠️ Single file, limited scale

**PostgreSQL (Recommended for production):**

### Railway:
1. Dashboard → New → Database → PostgreSQL
2. Copy `DATABASE_URL`
3. Add to environment variables

### Render:
1. New → PostgreSQL
2. Copy `Internal Database URL`
3. Add to your web service environment

### Update code to use PostgreSQL:
Edit [main.go](main.go):
```go
// Change from:
db, err := gorm.Open(sqlite.Open("ipopilot.db"), &gorm.Config{})

// To:
databaseURL := os.Getenv("DATABASE_URL")
if databaseURL == "" {
    databaseURL = "sqlite://ipopilot.db" // fallback
}
db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
```

Add PostgreSQL driver to [go.mod](go.mod):
```bash
go get gorm.io/driver/postgres
```

---

## 🔒 Production Checklist

Before going live:

```bash
# 1. Change JWT secret
# Edit utils.go or set JWT_SECRET env var

# 2. Change admin password
# Login and update via dashboard

# 3. Set production mode
export GIN_MODE=release

# 4. Use PostgreSQL
# Switch from SQLite

# 5. Enable HTTPS
# Automatic on Railway/Render/Vercel

# 6. Set up domain (optional)
# Add custom domain in platform dashboard
```

---

## 🌍 Custom Domain (Optional)

**On Railway:**
1. Settings → Domains → Add Domain
2. Enter: `ipo-pilot.com`
3. Update DNS: Add CNAME record
4. SSL auto-enabled

**On Render:**
1. Settings → Custom Domain → Add
2. Enter: `ipo-pilot.com`
3. Update DNS records as shown
4. SSL auto-enabled

**On Vercel:**
1. Settings → Domains → Add
2. Enter your domain
3. Follow DNS instructions
4. SSL auto-enabled

---

## 📊 Monitoring Your Live Site

**Railway Dashboard:**
- Real-time logs
- CPU/Memory usage
- Request metrics
- Automatic restarts

**Render Dashboard:**
- Build logs
- Deploy history
- Performance metrics
- Health checks

**Heroku Dashboard:**
- Dyno metrics
- Log streams
- Error tracking
- Add-ons available

---

## 🎯 Quick Comparison

| Platform | Free Tier | Setup Time | Auto-Deploy | Best For |
|----------|-----------|------------|-------------|----------|
| **Railway** | $5/month credit | 3 min | ✅ Yes | Easiest start |
| **Render** | 750 hrs/month | 5 min | ✅ Yes | 100% free forever |
| **Fly.io** | 3 VMs free | 5 min | ✅ Yes | Go apps |
| **Heroku** | 550 hrs/month | 10 min | ✅ Yes | Classic, reliable |
| **Vercel** | Unlimited | 2 min | ✅ Yes | Fastest deploys |

---

## ✅ Complete Example: GitHub → Railway

**Full workflow from code to live URL:**

```bash
# 1. Prepare code
cd /workspaces/IPO-PILOT/web-app

# 2. Make sure everything works locally
go run .
# Test at http://localhost:8080

# 3. Create .gitignore
cat > .gitignore << 'EOF'
*.db
*.db-shm
*.db-wal
.env
.DS_Store
app
tmp/
EOF

# 4. Initialize git
git init
git add .
git commit -m "Initial commit - IPO Pilot Web Platform"

# 5. Push to GitHub
# Option A: Using GitHub CLI
gh repo create ipo-pilot-web --public --source=. --push

# Option B: Manual
# - Go to github.com/new
# - Create repo: ipo-pilot-web
# - Run:
git remote add origin https://github.com/dipudai/ipo-pilot-web.git
git branch -M main
git push -u origin main

# 6. Deploy on Railway
# - Visit: railway.app
# - Click "Start a New Project"
# - Click "Deploy from GitHub repo"
# - Select "dipudai/ipo-pilot-web"
# - Click "Deploy"

# 7. Wait 2-3 minutes...

# 8. Done! Your site is live at:
# https://ipo-pilot-production.up.railway.app

# 9. Visit and login:
# Email: admin@ipopilot.com
# Password: admin123
```

**From now on, every time you push to GitHub, it auto-deploys!**

```bash
# Make changes to code
nano main.go

# Commit and push
git add .
git commit -m "Updated feature"
git push

# Railway automatically deploys new version!
# Live in 1-2 minutes
```

---

## 🆘 Troubleshooting

### Build fails on platform

**Error:** `cannot find package`
```bash
# Make sure go.mod is committed
git add go.mod go.sum
git commit -m "Add dependencies"
git push
```

### App crashes on startup

**Check logs:**
- Railway: Deployments → View Logs
- Render: Logs tab
- Heroku: `heroku logs --tail`

**Common fix:** Set PORT env var
```bash
# Platform sets PORT automatically, use it:
port := os.Getenv("PORT")
if port == "" {
    port = "8080"
}
r.Run(":" + port)
```

### Database issues

**SQLite permission error:**
```bash
# Switch to PostgreSQL for production
# OR ensure writable directory
```

### Can't access site

**Check:**
1. Is deployment successful? (green checkmark)
2. Is health check passing?
3. Is PORT configured correctly?
4. Check platform status page

---

## 🎉 You're Live!

Once deployed, share your URL:
- **Railway:** `https://ipo-pilot-production.up.railway.app`
- **Render:** `https://ipo-pilot.onrender.com`
- **Fly.io:** `https://ipo-pilot.fly.dev`
- **Heroku:** `https://ipo-pilot.herokuapp.com`
- **Vercel:** `https://ipo-pilot.vercel.app`

Users can:
1. Visit your URL
2. Register account
3. Choose subscription
4. Start using IPO automation

You can:
1. Access admin panel: `/admin`
2. Monitor users
3. Manage subscriptions
4. View analytics
5. Earn revenue!

---

## 📱 Next Steps

1. **Test Everything**
   - Create test user account
   - Try all features
   - Test on mobile

2. **Set Up Payments**
   - Integrate payment gateway
   - Test subscription flow
   - Verify webhooks work

3. **Custom Domain**
   - Buy domain (Namecheap, GoDaddy)
   - Configure DNS
   - Enable in platform

4. **Marketing**
   - Share URL on social media
   - Create landing page content
   - SEO optimization

5. **Monitor & Scale**
   - Watch error logs
   - Track user growth
   - Upgrade plan as needed

---

**Your web app is ready for the world!** 🚀

Any platform you choose will give you a live URL accessible from anywhere. Railway is recommended for beginners due to its simplicity.

**Questions? Check the platform-specific docs or deployment logs for detailed errors.**
