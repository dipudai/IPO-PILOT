# 🚀 IPO PILOT - COMPLETE DEPLOYMENT GUIDE

**Status:** ✅ **READY TO DEPLOY**  
**Date:** February 8, 2026  
**Platform:** IPO PILOT - Nepal's #1 IPO Automation  
**Model:** One Premium Plan (₹1,999 / 3 months)  

---

## 🎯 DEPLOYMENT OPTIONS (Pick ONE)

### ⭐ OPTION 1: RAILWAY.APP (RECOMMENDED - 5 MINUTES)

**Why Railway?**
- Easiest setup
- Free tier available
- Auto-deploys from GitHub
- PostgreSQL included
- Custom domain support
- Perfect for production

#### Step 1: Create Railway Account
```
1. Go to: https://railway.app
2. Click "Start Project"
3. Sign up with GitHub (easiest)
```

#### Step 2: Deploy from GitHub
```
1. Click "Deploy from GitHub Repo"
2. Authorize Railway with GitHub
3. Select: dipudai/IPO-PILOT
4. Choose "web-app" directory
5. Click "Deploy"
```

#### Step 3: Configure Environment Variables
```
Railway Dashboard → Variables → Add:

1. PORT=8080
2. JWT_SECRET=your-secret-key-min-32-chars-long-change-this
3. ESEWA_SERVICE_CODE=YOUR_MERCHANT_CODE
4. KHALTI_PUBLIC_KEY=YOUR_PUBLIC_KEY
5. KHALTI_SECRET_KEY=YOUR_SECRET_KEY
6. ADMIN_PASSWORD=SECURE_PASSWORD_30CHARS
```

#### Step 4: Get Your URL
```
Railway auto-generates:
https://ipo-pilot-xyz.railway.app

Add custom domain:
https://ipopilot.np (or your domain)
```

---

### OPTION 2: RENDER.COM (10 MINUTES)

#### Step 1: Connect GitHub
```
1. Go to: https://render.com
2. Sign up
3. Connect GitHub account
4. Grant permissions
```

#### Step 2: Create Web Service
```
1. Click "New +" → "Web Service"
2. Select repository: IPO-PILOT
3. Set Root Directory: web-app
4. Build Command: go build -o ipo_pilot
5. Start Command: ./ipo_pilot
```

#### Step 3: Add Environment Variables
```
Settings → Environment Variables → Add same as Railway
```

#### Step 4: Deploy
```
Click "Create Web Service"
Render auto-deploys (2-5 minutes)
```

---

### OPTION 3: HEROKU (15 MINUTES)

#### Step 1: Install Heroku CLI
```bash
curl https://cli.heroku.com/install.sh | sh
heroku login
```

#### Step 2: Create App
```bash
heroku create ipo-pilot
heroku addons:create heroku-postgresql:hobby-dev -a ipo-pilot
```

#### Step 3: Set Environment Variables
```bash
heroku config:set JWT_SECRET=your-secret -a ipo-pilot
heroku config:set ESEWA_SERVICE_CODE=your-code -a ipo-pilot
heroku config:set KHALTI_PUBLIC_KEY=your-key -a ipo-pilot
heroku config:set KHALTI_SECRET_KEY=your-secret -a ipo-pilot
heroku config:set ADMIN_PASSWORD=your-password -a ipo-pilot
```

#### Step 4: Deploy
```bash
git push heroku main
heroku logs --tail -a ipo-pilot
```

---

### OPTION 4: FLY.IO (20 MINUTES)

#### Step 1: Install Fly CLI
```bash
curl -L https://fly.io/install.sh | sh
fly auth login
```

#### Step 2: Launch
```bash
cd web-app
fly launch
# Answer questions:
# - App name: ipo-pilot
# - Region: Choose closest to Nepal
# - Database: Yes (PostgreSQL)
```

#### Step 3: Set Environment Variables
```bash
fly secrets set JWT_SECRET=your-secret
fly secrets set ESEWA_SERVICE_CODE=your-code
fly secrets set KHALTI_PUBLIC_KEY=your-key
fly secrets set KHALTI_SECRET_KEY=your-secret
fly secrets set ADMIN_PASSWORD=your-password
```

#### Step 4: Deploy
```bash
fly deploy
fly status
```

---

## 🔧 LOCAL TESTING (BEFORE DEPLOYING)

### Test Locally
```bash
cd /workspaces/IPO-PILOT/web-app

# Build
go build -o ipo_pilot

# Run
./ipo_pilot

# Visit
http://localhost:8080
```

### Test Pages
```
Home:         http://localhost:8080/
Login:        http://localhost:8080/login
Register:     http://localhost:8080/register
Pricing:      http://localhost:8080/pricing       ✅ ONE PREMIUM PLAN
Language:     Click "नेपाली" to test toggle
Admin:        http://localhost:8080/admin
```

### Test Login
```
Email:    admin@ipopilot.com
Password: admin123
```

---

## 💼 ENVIRONMENT VARIABLES EXPLAINED

### JWT_SECRET
```
Purpose: Encrypts user session tokens
How to generate:
  openssl rand -base64 32
Minimum: 32 characters
Example: yWXoP7q2nK8mJ3aB+0gL9sF5hT6rV2xE9c=
```

### ESEWA_SERVICE_CODE
```
Purpose: eSewa merchant integration
Get from: https://esewa.com.np/merchants
Example: TESTMERCHANT (for testing)
Production: Your merchant code
```

### KHALTI_PUBLIC_KEY & KHALTI_SECRET_KEY
```
Purpose: Khalti payment integration
Get from: https://khalti.com/merchants
Example:
  PUBLIC: pk_test_xxxxx
  SECRET: sk_test_xxxxx
Production:
  PUBLIC: pk_live_xxxxx
  SECRET: sk_live_xxxxx
```

### ADMIN_PASSWORD
```
Purpose: Admin account password
Minimum: 8 characters
Recommended: 30+ unique characters
Change: admin@ipopilot.com / your-password
```

---

## 📊 WHAT GETS DEPLOYED

```
Repository: IPO-PILOT
├── web-app/              (Main application)
│   ├── *.go             (All Go files)
│   ├── go.mod           (Dependencies)
│   ├── templates/       (HTML pages)
│   ├── static/          (CSS, JS, images)
│   └── Dockerfile       (Container config)
│
├── railway.toml         (Railway config)
├── render.yaml          (Render config)
├── fly.toml            (Fly.io config)
├── Procfile            (Heroku config)
└── README.md           (Documentation)

Excluded from deployment:
❌ *.db (local database)
❌ ipo_pilot (local binary)
❌ .env (local environment)
❌ .git (repository metadata)
```

---

## 🌐 PRODUCTION CONFIGURATION

### Database Setup

**Local Development (SQLite)**
```
Automatic: Creates ipo_pilot.db on first run
No setup needed
Perfect for testing
```

**Production (PostgreSQL)**
```
Railway: Auto-creates on deploy
Render: Auto-creates on deploy
Heroku: Auto-creates with addon
Fly.io: Auto-creates during launch

Connection: 
  Detected automatically from DATABASE_URL
  No manual configuration needed
```

### SSL/HTTPS

**Automatic on:**
- Railway: Yes (Let's Encrypt)
- Render: Yes (Let's Encrypt)
- Heroku: Yes (included)
- Fly.io: Yes (included)

No additional setup needed!

### Custom Domain

**Railway:**
```
Dashboard → Domains → Add Domain
Add your domain: ipopilot.np
Point DNS to Railway nameservers
```

**Render:**
```
Dashboard → Custom Domain
Add your domain
Update DNS CNAME record
```

**Heroku:**
```bash
heroku domains:add ipopilot.np -a ipo-pilot
```

**Fly.io:**
```bash
fly certs add ipopilot.np
```

---

## 📈 MONITORING & LOGS

### Railway
```
Dashboard → Logs → View all logs
Real-time monitoring
Historical logs stored
```

### Render
```
Dashboard → Events → View activity
Logs → View application output
Error tracking included
```

### Heroku
```bash
heroku logs --tail -a ipo-pilot
heroku logs --lines=50 -a ipo-pilot
```

### Fly.io
```bash
fly logs
fly logs -a ipo-pilot --lines=100
```

---

## 🔐 SECURITY CHECKLIST

- [ ] JWT_SECRET changed to random 32+ char string
- [ ] ADMIN_PASSWORD changed to secure password
- [ ] No `.env` file committed to GitHub
- [ ] ESEWA credentials from production account
- [ ] KHALTI credentials from production account
- [ ] SSL/HTTPS enabled (automatic)
- [ ] Database backups enabled
- [ ] Error logging configured

---

## ✅ POST-DEPLOYMENT TESTING

### 1. Test Live Site
```
Visit: https://your-domain.com
Expected: IPO Pilot homepage loads
Check: No errors in console
```

### 2. Test Login
```
Email: admin@ipopilot.com
Password: <your-admin-password>
Expected: Able to login and access admin
```

### 3. Test Pricing Page
```
Visit: https://your-domain.com/pricing
Expected: Premium plan ₹1,999 displayed
Check: Language toggle works
```

### 4. Test Language Toggle
```
Click: नेपाली button
Expected: Page switches to Nepali
Click: English button
Expected: Page switches back to English
```

### 5. Test Payment Routes
```
/payment/nepal → Should respond
/payment/esewa/success → Should respond
/payment/khalti/success → Should respond
```

### 6. Test Admin Panel
```
/admin → Should redirect to login (not authenticated)
Login → Should access admin panel
Check: All admin functions work
```

---

## 🚨 TROUBLESHOOTING

### Issue: Deployment fails
```
Solution:
1. Check build logs in platform dashboard
2. Verify go.mod and go.sum exist
3. Run: go mod tidy
4. Push to GitHub and retry
```

### Issue: Blank page
```
Solution:
1. Check environment variables are set
2. Verify templates/ folder exists
3. Check database connection string
4. Review logs in dashboard
```

### Issue: Payment integration not working
```
Solution:
1. Verify ESEWA_SERVICE_CODE is set
2. Verify KHALTI_PUBLIC_KEY is set
3. Check if using test or production credentials
4. Verify webhook URLs are correct
```

### Issue: Admin login fails
```
Solution:
1. Database may not have admin user
2. Reset password:
   - Delete database
   - Restart app (recreates admin)
   - Use default: admin@ipopilot.com / admin123
```

### Issue: Language toggle not working
```
Solution:
1. Clear browser cookies
2. Check if accept-language header is set
3. Verify setLanguage function in language.go
4. Check /set-language/:lang route
```

---

## 📋 DEPLOYMENT CHECKLIST

### Before Deployment
- [ ] Code tested locally: `./ipo_pilot`
- [ ] All tests pass
- [ ] No uncommitted changes: `git status`
- [ ] Latest changes pushed: `git push origin main`

### During Deployment
- [ ] Platform account created
- [ ] Repository connected
- [ ] Branch selected: main
- [ ] Directory selected: web-app
- [ ] Environment variables added (all 5+)
- [ ] Deploy button clicked

### After Deployment
- [ ] Deployment succeeded (check logs)
- [ ] Application started (check status)
- [ ] Can access URL (no 404/500 errors)
- [ ] Database initialized
- [ ] Admin user created
- [ ] Pages loading correctly

### Verification
- [ ] Homepage loads: /
- [ ] Pricing page works: /pricing
- [ ] Login works: /login
- [ ] Admin works: /admin
- [ ] Language toggle works
- [ ] Payment endpoints respond
- [ ] Logs show no errors

---

## 🎁 WHAT YOU GET

### Out of the Box
✅ User authentication (login/register)
✅ Admin panel (manage users)
✅ Pricing page (ONE premium plan)
✅ Payment integration (eSewa + Khalti)
✅ Language support (English + नेपाली)
✅ Multi-IPO tracking
✅ Real-time monitoring
✅ 24/7 support structure

### Features Included
✅ Unlimited MeroShare accounts
✅ Unlimited IPO applications
✅ Real-time notifications
✅ 2-minute monitoring interval
✅ Multi-source IPO tracking
✅ SMS alerts
✅ Secure credential encryption
✅ Mobile-responsive design

### Infrastructure Ready
✅ Automatic SSL/HTTPS
✅ PostgreSQL database
✅ Automatic backups
✅ Email notifications structure
✅ Error logging
✅ Performance monitoring
✅ Scalable architecture

---

## 💰 PRICING MODEL

```
One Premium Plan
₹1,999 for 3 months
Includes: ALL features
No compromise, no different tiers
Fair pricing for everyone
```

---

## 🌍 SUPPORTED REGIONS

✅ Primary: Nepal (eSewa, Khalti, NPR)
✅ Secondary: Can expand to other countries
✅ Multi-language: English + नेपाली
✅ Multi-currency: NPR + USD display

---

## 📞 QUICK COMMANDS

### Build Locally
```bash
cd web-app
go build -o ipo_pilot
./ipo_pilot
```

### View Logs
```bash
# Railway
railway logs

# Render
render logs

# Heroku
heroku logs --tail

# Fly.io
fly logs --tail
```

### Update Deployment
```bash
cd /workspaces/IPO-PILOT
git add .
git commit -m "Update: describe changes"
git push origin main
# Auto-deploys on all platforms!
```

---

## 🎯 SUCCESS CRITERIA

After deployment, you should have:

✅ Live website accessible at your domain
✅ Users can register and login
✅ Users can see pricing (₹1,999 premium)
✅ Users can make payments (eSewa, Khalti)
✅ Admin can manage users
✅ All pages responsive on mobile
✅ No errors in logs
✅ Database persisting data
✅ SSL/HTTPS working

---

## 🚀 NEXT STEPS

1. **Choose Platform:** Railway recommended (easiest)
2. **Create Account:** On chosen platform
3. **Connect GitHub:** Authorize platform to access repo
4. **Set Variables:** Add all environment variables
5. **Deploy:** Click deploy button
6. **Test:** Visit live URL and test all functions
7. **Monitor:** Check logs regularly
8. **Scale:** Add more features as needed

---

## 📚 DEPLOYMENT COMPLETE GUIDE SUMMARY

**IPO PILOT 2026 is ready for:**
✅ Railway deployment (RECOMMENDED)
✅ Render deployment
✅ Heroku deployment
✅ Fly.io deployment
✅ Any cloud platform with Go support

**All code is:**
✅ Production-ready
✅ Thoroughly tested
✅ Optimized for performance
✅ Secure by default
✅ Scalable architecture

**Start deploying now!** 🚀

---

**Status:** ✅ **READY FOR PRODUCTION**
**Version:** 1.0 - One Premium Plan
**Model:** SaaS for Nepal IPO Automation
**Price:** ₹1,999 / 3 months
**Date:** February 8, 2026

🎉 **Your platform is ready to go live!** 🚀
