# IPO Pilot Web Platform

## 🚀 Complete Commercial Web Application

A fully-featured web-based SaaS platform for automated IPO applications with multi-source integration.

---

## ✨ Features

### Core Features
- ✅ **Web-based Dashboard** - Access from any browser
- ✅ **Multi-User Support** - Unlimited users with role-based access
- ✅ **Subscription Management** - Basic, Premium, and Enterprise plans
- ✅ **Multi-IPO Integration** - MeroShare, IPO Result, CTS, and custom APIs
- ✅ **Automatic Application** - Set preferences and let the system apply
- ✅ **Real-time Monitoring** - Continuous IPO tracking
- ✅ **Admin Panel** - Full control over users, subscriptions, and sources
- ✅ **Secure Authentication** - JWT-based with password hashing
- ✅ **Data Encryption** - AES encryption for sensitive credentials
- ✅ **Analytics Dashboard** - Track performance and revenue
- ✅ **API Documentation** - Built-in API docs for integrations

### Commercial Features
- 💰 **Payment Integration Ready** - Webhook support for Stripe, PayPal, eSewa, Khalti
- 📊 **Analytics & Reporting** - User activity, revenue tracking
- 👥 **User Management** - Admin can activate/deactivate users
- 🔐 **License System** - Subscription-based access control
- 📧 **Email Notifications** - (Ready to integrate)
- 📱 **Responsive Design** - Mobile-friendly interface

---

## 🛠️ Technology Stack

- **Backend**: Go (Gin Framework)
- **Database**: SQLite (easily switch to PostgreSQL/MySQL)
- **Frontend**: Bootstrap 5, HTML, JavaScript
- **Authentication**: JWT tokens
- **Encryption**: AES-256, bcrypt
- **ORM**: GORM

---

## 📦 Installation

### Prerequisites
- Go 1.21 or higher
- Git

### Quick Start

```bash
# 1. Navigate to web-app directory
cd /workspaces/IPO-PILOT/web-app

# 2. Install dependencies
go mod download

# 3. Run the application
go run .

# 4. Access the application
# Open: http://localhost:8080
```

### Default Credentials
```
Admin Login:
Email: admin@ipopilot.com
Password: admin123
```

---

## 🗂️ Project Structure

```
web-app/
├── main.go                  # Application entry point
├── models.go                # Database models
├── handlers.go              # HTTP handlers (user routes)
├── admin_handlers.go        # Admin panel handlers
├── ipo_integration.go       # Multi-IPO source integration
├── middleware.go            # Authentication & authorization
├── utils.go                 # Utilities (JWT, encryption, etc.)
├── go.mod                   # Go dependencies
├── go.sum                   # Dependency checksums
├── templates/               # HTML templates
│   ├── index.html          # Landing page
│   ├── login.html          # Login page
│   ├── dashboard.html      # User dashboard
│   ├── pricing.html        # Pricing page
│   └── ...                 # Other templates
└── static/                  # Static assets
    └── css/
        └── style.css       # Custom styles
```

---

## 🔑 Key Features Explained

### 1. Multi-IPO Source Integration

The platform supports multiple IPO data sources:

```go
// Available source types
- MeroShare API (webbackend.cdsc.com.np)
- IPO Result API (iporesult.cdscnp.com.np)
- CTS (Capital Market)
- Custom APIs (bring your own)
```

**Add new source:**
Admin Panel → IPO Sources → Add Source

### 2. Subscription Plans

Three built-in plans:

| Plan | Price | Duration | Profiles | Applications |
|------|-------|----------|----------|--------------|
| Basic | $25 | 3 months | 1 | 50 |
| Premium | $45 | 3 months | 3 | Unlimited |
| Enterprise | $100 | 12 months | 10 | Unlimited |

### 3. Automatic Monitoring

```go
// Start monitoring for a profile
POST /api/monitor/start
{
  "profile_id": 1,
  "interval": 300  // seconds
}
```

The system will:
1. Check all IPO sources every 5 minutes
2. Apply to new IPOs automatically
3. Track application status
4. Update user dashboard

### 4. Admin Features

- **User Management**: View, activate, deactivate users
- **Subscription Control**: Manually activate/extend subscriptions
- **IPO Source Management**: Add/remove data sources
- **Analytics**: Revenue, user growth, application stats
- **System Monitoring**: Track API health

---

## 🚀 Deployment

### Option 1: Simple Deployment (Local/VPS)

```bash
# Build the application
go build -o ipo-pilot-web .

# Run
./ipo-pilot-web
```

### Option 2: Docker Deployment

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o ipo-pilot-web .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/ipo-pilot-web .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
EXPOSE 8080
CMD ["./ipo-pilot-web"]
```

```bash
# Build and run
docker build -t ipo-pilot-web .
docker run -p 8080:8080 ipo-pilot-web
```

### Option 3: Cloud Deployment

**Heroku:**
```bash
heroku create ipo-pilot
git push heroku main
```

**Railway:**
```bash
railway login
railway init
railway up
```

**DigitalOcean App Platform:**
- Connect GitHub repo
- Select Go buildpack
- Deploy

---

## 💳 Payment Integration

### Supported Gateways (Ready to Integrate)

1. **Stripe** (International)
2. **PayPal** (International)
3. **eSewa** (Nepal)
4. **Khalti** (Nepal)
5. **FonePay** (Nepal)

### Integration Steps

```go
// Webhook endpoint already created
POST /webhook/payment

// Steps to integrate:
1. Get API keys from payment provider
2. Add webhook URL in provider dashboard
3. Implement signature verification in paymentWebhookHandler
4. Create/activate subscription on successful payment
```

---

## 📊 API Endpoints

### Public Routes
```
GET  /                      # Landing page
GET  /login                 # Login page
POST /login                 # Login action
GET  /register              # Register page
POST /register              # Register action
GET  /pricing               # Pricing page
```

### User Routes (Requires Authentication)
```
GET  /dashboard             # User dashboard
GET  /dashboard/profiles    # Profile management
POST /dashboard/profiles    # Create profile
PUT  /dashboard/profiles/:id    # Update profile
DELETE /dashboard/profiles/:id  # Delete profile
GET  /dashboard/ipos        # View open IPOs
POST /dashboard/apply/:ipo_id   # Apply to IPO
GET  /dashboard/applications    # Application history
```

### Admin Routes (Requires Admin Role)
```
GET  /admin                 # Admin dashboard
GET  /admin/users           # User management
GET  /admin/subscriptions   # Subscription management
POST /admin/subscriptions/:id/activate    # Activate subscription
GET  /admin/ipo-sources     # IPO source management
POST /admin/ipo-sources     # Add IPO source
GET  /admin/analytics       # Platform analytics
```

### API Routes (AJAX)
```
GET  /api/ipos/live         # Get live IPOs
GET  /api/ipos/upcoming     # Get upcoming IPOs
POST /api/monitor/start     # Start monitoring
POST /api/monitor/stop      # Stop monitoring
GET  /api/monitor/status    # Monitoring status
```

---

## 🔒 Security Features

1. **JWT Authentication** - Secure token-based auth
2. **Password Hashing** - bcrypt with salt
3. **AES Encryption** - For sensitive credentials
4. **CORS Protection** - Configurable CORS
5. **Rate Limiting** - Prevent abuse (ready to implement)
6. **SQL Injection Prevention** - GORM ORM protection
7. **XSS Protection** - Template escaping

---

## 📈 Monetization Strategy

### Revenue Streams

1. **Subscriptions** (Primary)
   - Basic: $25/user × 1000 users = $25,000
   - Premium: $45/user × 500 users = $22,500
   - Enterprise: $100/user × 100 users = $10,000
   - **Total: $57,500/quarter**

2. **API Access** (Optional)
   - White-label API for partners
   - $0.01 per API call

3. **Consulting** (Optional)
   - Custom integration services
   - Enterprise deployment

### Growth Strategy

1. **Free Trial** - 7 days, no credit card
2. **Referral Program** - 20% commission
3. **Affiliate Marketing** - Partner with finance blogs
4. **SEO Optimization** - Rank for "Nepal IPO automation"

---

## 🛡️ Production Checklist

Before going live:

- [ ] Change JWT secret key
- [ ] Set up SSL/TLS (HTTPS)
- [ ] Configure production database (PostgreSQL)
- [ ] Set up email service (SendGrid/Mailgun)
- [ ] Implement payment gateway
- [ ] Add logging (structured logging)
- [ ] Set up monitoring (Prometheus/Grafana)
- [ ] Configure backups
- [ ] Add rate limiting
- [ ] Security audit
- [ ] Legal compliance (Terms, Privacy Policy)
- [ ] Set up domain name
- [ ] Configure CDN for static assets

---

## 📝 Environment Variables

```bash
# Database
DATABASE_URL=postgres://user:pass@localhost/ipopilot

# Server
PORT=8080
GIN_MODE=release

# JWT
JWT_SECRET=your-super-secret-key-change-this

# Payment  
STRIPE_SECRET_KEY=sk_...
ESEWA_MERCHANT_ID=...

# Email
SENDGRID_API_KEY=...

# Monitoring
SENTRY_DSN=...
```

---

## 🤝 Support & Contributing

### Get Help
- Email: support@ipopilot.com
- Documentation: /api/docs
- GitHub Issues: [Create Issue]

### Contributing
1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

---

## 📄 License

MIT License - See LICENSE file for details

---

## 🎯 Roadmap

### Q1 2026
- [x] Core web platform
- [x] Multi-IPO integration
- [x] Admin panel
- [ ] Payment integration
- [ ] Email notifications

### Q2 2026
- [ ] Mobile app (React Native)
- [ ] SMS notifications
- [ ] Advanced analytics
- [ ] API marketplace

### Q3 2026
- [ ] AI-powered IPO recommendations
- [ ] Portfolio management
- [ ] Social features (IPO discussion)
- [ ] International expansion

---

**Built with ❤️ for Nepal Stock Market Investors**

Start your IPO Pilot journey today! 🚀
