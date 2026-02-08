# 🔐 IPO Pilot Login Troubleshooting Guide

## The Problem You're Experiencing

After registering and trying to login, you might see:
- ❌ Redirects back to login page
- ❌ Error message appears
- ❌ Page goes blank
- ❌ Dashboard doesn't load

## How Authentication Works (Now Fixed!)

### Step 1: User Submits Login Form
```
POST /login
Content-Type: application/json
{
  "email": "user@example.com",
  "password": "password123"
}
```

### Step 2: Backend Validates & Returns Token
```
Response:
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "User Name",
    "isAdmin": false
  }
}
```

### Step 3: Frontend Stores Token (FIXED!)
```javascript
// Now does BOTH:
localStorage.setItem('auth_token', data.token);  // For reference
document.cookie = `auth_token=${data.token}; path=/; max-age=86400`;  // For requests
```

### Step 4: Browser Redirects to Dashboard
```
GET /dashboard
Cookie: auth_token=eyJhbGciOiJIUzI1NiIs...
```

### Step 5: Middleware Validates Cookie
```go
// authMiddleware() in middleware.go:
1. Looks for token in Authorization header
2. Falls back to "auth_token" cookie  ✅ (This is what we fixed)
3. Validates JWT signature
4. Sets userID in request context
5. Allows dashboardHandler to access userID
```

### Step 6: Dashboard Renders
```
✅ Dashboard page loads with user data
```

---

## 🧪 Testing Steps

### In Your Browser Developer Console (F12 → Console Tab):

#### Test 1: Check if Token is Stored
```javascript
// Should show your token
console.log(localStorage.getItem('auth_token'));

// Should show your user info
console.log(localStorage.getItem('user'));

// Should show auth_token cookie exists
console.log(document.cookie);
```

#### Test 2: Monitor Lock/Register Flow
```javascript
// Open Console tab (F12) and watch the logs
// During registration:
// - You should see: "Attempting registration with email: test@example.com"
// - Then: "Response: 201 {...}"
// - Then redirect to login

// During login:
// - You should see: "Attempting login with email: test@example.com"
// - Then: "Response: 200 {...}" (with token!)
// - Then: "Login successful, redirecting..."
// - Then browser redirects to /dashboard
```

#### Test 3: Check Cookies are Sent
```javascript
// In Network tab (F12 → Network):
// 1. Refresh the page
// 2. Click on the first request (usually "localhost" or "dashboard")
// 3. Look for Request Headers → Cookie
// 4. Should see: auth_token=eyJhbGc...
```

---

## 🔧 If You Still Have Issues

### Issue 1: Token Not Appearing in Console
**Problem:** `console.log(localStorage.getItem('auth_token'))` shows `null`

**Solutions:**
- ❌ Refresh page (clears localStorage if no persistence)
- ✅ Check if localStorage is enabled (Settings → Storage)
- ✅ Check browser privacy settings
- ✅ Try in Incognito/Private mode

### Issue 2: Cookie Not Showing
**Problem:** `document.cookie` doesn't include `auth_token`

**Solutions:**
```javascript
// Check if cookie setter worked
document.cookie = `auth_token=test-value; path=/; max-age=86400`;
console.log(document.cookie);  // Should now show "auth_token=test-value"

// Some browsers block cookies in certain conditions:
// - Incognito/Private mode
// - Third-party cookie restrictions
// - HTTPS only (for non-localhost)
```

### Issue 3: Still Redirected to Login
**Problem:** After login, browser redirects back to login

**Solutions:**
1. **Check Network Tab** (F12 → Network):
   - Look for redirect responses (301, 302, 307)
   - Click the request and check Headers
   - Verify credentials are in the request

2. **Test with curl** (no JavaScript issues):
   ```bash
   # 1. Register
   curl -X POST http://localhost:8080/register \
     -H "Content-Type: application/json" \
     -d '{"name":"Test","email":"test@ex.com","password":"pass123"}'
   
   # 2. Login
   curl -X POST http://localhost:8080/login \
     -H "Content-Type: application/json" \
     -d '{"email":"test@ex.com","password":"pass123"}'
   # → Copy the "token" value
   
   # 3. Access dashboard with token cookie
   curl -b "auth_token=YOUR_TOKEN_HERE" http://localhost:8080/dashboard
   # → Should show dashboard HTML, not redirect to login
   ```

### Issue 4: Error Message Displayed
**Problem:** You see "Invalid credentials" or "Login failed"

**Possible Causes:**
- ❌ Email/password mismatch - verify exact spelling
- ❌ User doesn't exist - check you registered first
- ❌ Password too short - minimum 6 characters required
- ❌ User account deactivated

**Solution:**
- Check browser console for exact error message
- Verify registration was successful first
- Try the test credentials:
  - **Email:** `admin@ipopilot.com`
  - **Password:** `admin123`

---

## 🚀 Quick Fixes

### 1. Clear Browser Cache and Try Again
```javascript
// In console:
localStorage.clear();
// Then refresh page and try login again
```

### 2. Check if You're on Correct URL
- ✅ Correct: https://ipo-pilot-production.up.railway.app/login
- ❌ Wrong: localhost:8080/login (if on production Railway)

### 3. Wait 2-3 Minutes After Registration
- Some databases have replication delay
- Try login again after a few minutes

### 4. Check Browser Privacy Settings
- Cookies might be disabled
- Private/Incognito mode might block cookies
- Third-party cookie restrictions might apply

---

## 📝 Complete Login Flow Test

```javascript
// Run this in browser console (F12 → Console) after fixing login:

async function testLoginFlow() {
  try {
    // 1. Register
    console.log("1️⃣ Registering...");
    const regResp = await fetch('/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name: 'Flow Test',
        email: `test-${Date.now()}@example.com`,
        password: 'TestPass123'
      })
    });
    console.log("Registration:", regResp.status, await regResp.json());
    
    // Wait before login
    await new Promise(r => setTimeout(r, 1000));
    
    // 2. Login
    console.log("2️⃣ Logging in...");
    const email = JSON.parse(localStorage.getItem('user') || '{}').email || 'admin@ipopilot.com';
    const loginResp = await fetch('/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        email: email,
        password: 'TestPass123'
      })
    });
    const loginData = await loginResp.json();
    console.log("Login:", loginResp.status, loginData);
    
    // 3. Check token
    if (loginResp.ok) {
      console.log("✅ Login successful!");
      console.log("Token stored:", !!localStorage.getItem('auth_token'));
      console.log("Cookie set:", document.cookie.includes('auth_token'));
    }
  } catch (error) {
    console.error("❌ Test failed:", error);
  }
}

testLoginFlow();
```

---

## 🆘 Still Need Help?

If you're still having issues after trying these:

1. **Take a screenshot** of the error
2. **Check browser console** for JavaScript errors
3. **Check Network tab** to see what requests are failing
4. **Share with us:**
   - The exact error message you see
   - Console output (F12 → Console tab)
   - Network request details (F12 → Network tab)

---

## ✅ After Login Works

Once you can login successfully:
- ✅ Dashboard loads with your profile stats
- ✅ Can navigate to other pages (Profiles, IPOs, etc.)
- ✅ Logout button works
- ✅ Can create/manage MeroShare profiles

**If dashboard loads but features don't work:** Check the Console tab for JavaScript errors

---

**Last Updated:** February 8, 2026  
**Platform:** IPO Pilot v1.0  
**Auth Method:** JWT + HTTP Cookies
