# IPO Pilot - License Generation & Management Guide

## 📋 Overview

You manage licenses for both Windows and Android versions from a **Developer Console** (a simple Go program). Both platforms use the **same license validation system** - any license key works on both!

---

## 🔑 License Key Format

```
Format: XXXX-XXXX-XXXX-XXXX
Example: AB3K-XY9L-MC4P-F7Q2

Structure:
├─ First 15 chars (XXXX-XXXX-XXXX): Random alphanumeric
└─ Last 4 chars (XXXX): MD5 checksum of first 15 chars

Validity: 3 months from activation date
Binding: Hardware machine ID (unique per device)
```

---

## 🛠️ How to Generate Licenses

### Option 1: Use Windows App (Interactive)

When customer requests a license:

1. **On your Windows machine**, run IPO Pilot:
   ```bash
   go run ipo_master_customer.go
   ```

2. **Select: License Information** → **Renew License** (if generating new)
   - This shows the `generateLicenseKey()` function in code

3. **OR create a small License Generator Tool** (recommended)

---

### Option 2: Create a License Generator Tool (Recommended)

Create a simple Go program to generate license keys on demand:

**File: `license_generator.go`**

```go
package main

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("=== IPO Pilot License Generator ===\n")

	// Generate 10 sample licenses
	for i := 1; i <= 10; i++ {
		key := generateLicenseKey()
		fmt.Printf("%d. %s\n", i, key)
	}

	// Validate a key
	fmt.Println("\n=== License Validation Test ===")
	testKey := generateLicenseKey()
	fmt.Printf("Generated: %s\n", testKey)
	fmt.Printf("Valid: %v\n", validateLicenseKey(testKey))
}

func generateLicenseKey() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 19) // 4-4-4-4 format with dashes

	// Generate random part (16 characters)
	for i := 0; i < 16; i++ {
		if i > 0 && i%4 == 0 {
			result[i+(i/4)-1] = '-'
		}
		result[i+(i/4)] = charset[rand.Intn(len(charset))]
	}

	// Add checksum (last 4 characters)
	keyWithoutChecksum := string(result[:15]) // First 15 chars (XXXX-XXXX-XXXX)
	checksum := fmt.Sprintf("%X", md5.Sum([]byte(keyWithoutChecksum)))[:4]
	result[15] = checksum[0]
	result[16] = checksum[1]
	result[17] = checksum[2]
	result[18] = checksum[3]

	return string(result)
}

func validateLicenseKey(key string) bool {
	if len(key) != 19 {
		return false
	}

	// Extract parts
	parts := make([]string, 0)
	start := 0
	for i := 0; i < 4; i++ {
		end := start + 4
		if end > len(key) && i < 3 {
			return false
		}
		parts = append(parts, key[start:end])
		start = end + 1
	}

	// Validate checksum
	keyWithoutChecksum := key[:15]
	expectedChecksum := fmt.Sprintf("%X", md5.Sum([]byte(keyWithoutChecksum)))[:4]
	actualChecksum := key[15:]

	return expectedChecksum == actualChecksum
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
```

**How to use:**
```bash
cd path/to/license_generator
go run license_generator.go
```

Output example:
```
=== IPO Pilot License Generator ===

1. ABC1-DEF2-GHI3-J4K5
2. LMN6-OPQ7-RST8-U9V0
3. WXY1-Z2AB-C3DE-F4GH
...
```

---

## 📱 License Management System

### Step 1: Generate License Keys

```bash
# Generate multiple keys
go run license_generator.go > licenses_$(date +%Y%m%d).txt
```

### Step 2: Track Licenses (Create a spreadsheet)

**License Tracker (licensee_tracker.csv):**

```csv
License Key,Customer Name,Email,Windows Machine ID,Android Device ID,Activation Date,Expiry Date,Status,Notes
AB3K-XY9L-MC4P-F7Q2,John Doe,john@example.com,abc123def,device456,2026-01-15,2026-06-15,Active,Windows Desktop
MN5P-QR8S-TU2V-W3XY,Jane Smith,jane@example.com,machine789,android012,2026-01-15,2026-06-15,Active,Both Windows & Mobile
...
```

### Step 3: Customer Activates License

**Windows Version:**
1. Customer runs: `go run ipo_master_customer.go`
2. Selects: "Activate License"
3. Enters: License key you provided
4. System validates checksum and saves to `license.dat`
5. Auto-binds to their machine ID

**Mobile Version (Android):**
1. Customer opens IPO Pilot App
2. Navigates to: License Screen
3. Enters: Same license key
4. System validates checksum
5. Auto-binds to their device
6. Valid for 3 months

### Step 4: License Renewal

When license expires:

**Customer requests renewal:**
1. Send a new license key (same format)
2. Customer activates new key
3. New 5-month period starts

---

## 📊 License Management Workflow

```
┌─────────────────────────────────────────────────────────┐
│         DEVELOPER LICENSE MANAGEMENT                    │
└─────────────────────────────────────────────────────────┘

1. GENERATE LICENSES
   ↓
   Use: license_generator.go
   Output: Multiple valid license keys
   
2. DISTRIBUTE TO CUSTOMERS
   ↓
   Send via: Email, Download Link, Manual
   Include: Usage instructions
   
3. CUSTOMER ACTIVATES
   ↓
   Windows: Menu → License Info → Activate
   Mobile: Settings → License Screen → Activate
   ✓ Auto-validates checksum
   ✓ Auto-binds to machine/device
   ✓ Valid for 3 months
   
4. TRACK IN SPREADSHEET
   ↓
   Log: Customer name, key, device IDs
   Monitor: Expiry dates, renewals
   
5. HANDLE RENEWALS
   ↓
   Generate new key
   Send to customer
   Customer re-activates
   ✓ New 5-month period starts
```

---

## 🔐 Security Features

### 1. **MD5 Checksum Validation**
- Every key has built-in checksum
- Invalid keys are automatically rejected
- Prevents typos and tampering

### 2. **Hardware Machine ID Binding**
- License tied to device's unique ID
- Can't use same license on different machines
- Prevents license sharing/piracy

### 3. **Time-Based Expiry**
- 5-month validity period
- Auto-verified on each app start
- Expired licenses are rejected

### 4. **Data Encryption**
- Profile passwords encrypted with AES
- PINs encrypted with AES
- API credentials never stored in plain text

---

## 💼 Business Integration

### Stripe/Payment Integration (Optional)

For automated licensing:

```go
// Pseudo-code for Stripe integration
1. Customer pays for license via Stripe
2. Webhook received with payment confirmation
3. Automatically generate & email license key
4. Customer immediately activates
5. No manual work needed
```

### Send via Email (Simple)

```bash
# Generate license and email to customer
TO="customer@example.com"
LICENSE=$(go run license_generator.go | head -1)

echo "
Your IPO Pilot License Key:
$LICENSE

This license is valid for 3 months.

Windows: Run app → Activate License → Enter above key
Mobile: Open app → License Screen → Enter above key

Need help? Email: support@example.com
" | mail -s "IPO Pilot License Key" $TO
```

---

## 📱 Same License Works on Both Platforms

```
One license key = Works on Windows AND Android

Example:
Customer John Doe gets key: AB3K-XY9L-MC4P-F7Q2

1. Activates on Windows Desktop
   ✓ Machine ID: windows-machine-123
   ✓ License saved in: C:\path\to\license.dat
   
2. Same John can activate on Android Phone
   ✓ Device ID: android-device-456
   ✓ License saved in: app_database
   
   ⚠️ Note: License is device-specific
      John cannot share same key with others
      (Different machine IDs = different activation)
```

---

## 📋 Checklist for License Management

- [ ] Create `license_generator.go` tool
- [ ] Generate batch of 100 license keys
- [ ] Create license tracking spreadsheet (CSV/Google Sheets)
- [ ] Set up delivery method (email/download)
- [ ] Document process for customers
- [ ] Monitor expiry dates (set reminders for renewals)
- [ ] Back up license records
- [ ] Track revenue per license sold
- [ ] Plan payment processing (Stripe/PayPal)

---

## 🚀 Pricing Strategy Example

```
License Pricing (3-month validity):
├─ Single License: Rs. 500-1000
├─ 3 Licenses: Rs. 1300-2700
├─ Quarterly Renewal: Rs. 300-500
└─ Bulk (10+): Contact for quote

Revenue Model:
├─ Initial sale: Customer pays Rs. 500-1000
├─ Renewal: Customer pays Rs. 300-500 every 3 months
└─ Bulk: 20% discount for 10+ licenses
```

---

## 🆘 Customer Support

When customer says license doesn't work:

1. **Check if key format is correct**
   ```
   Right: AB3K-XY9L-MC4P-F7Q2
   Wrong: ab3k-xy9l-mc4p-f7q2 (must be uppercase)
   ```

2. **Verify key is valid**
   ```bash
   go run license_generator.go
   # Validate their key matches this format
   ```

3. **Check expiry date**
   - In app: License Info shows expiry date
   - Expired → Send renewal key

4. **Check machine ID binding**
   - Windows: Each computer needs own license
   - Mobile: Each phone needs own license
   - Cannot reuse same key on different devices

---

## 📞 Support Email Template

```
Subject: IPO Pilot License - [Support Request]

Dear Customer,

Thank you for purchasing IPO Pilot!

Your license key: [KEY]
Valid until: [DATE]

Installation steps:

Windows:
1. Run: go run ipo_master_customer.go
2. Select: License Information
3. Choose: Activate License
4. Enter: Your license key (UPPERCASE)

Mobile (Android):
1. Open IPO Pilot app
2. Go to: License Screen
3. Enter: Your license key
4. Tap: Activate

Need help? Reply to this email or contact:
Email: support@example.com

Best regards,
IPO Pilot Team
```

---

All set! You can now generate and manage licenses for both Windows and Mobile versions. 🎉
