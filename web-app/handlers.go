package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Home page
func homeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "IPO Pilot - Automated IPO Application Platform",
	})
}

// Login page
func loginPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

// Terms of Service page
func termsHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "terms.html", gin.H{
		"title": "Terms of Service - IPO Pilot",
	})
}

// Privacy Policy page
func privacyHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "privacy.html", gin.H{
		"title": "Privacy Policy - IPO Pilot",
	})
}

// Login handler
func loginHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user User
	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account is not active"})
		return
	}

	if !checkPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := generateJWT(user.ID, user.IsAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Get subscription/trial info
	var subscription Subscription
	trialInfo := gin.H{}
	if err := db.Where("user_id = ? AND status = ?", user.ID, "active").First(&subscription).Error; err == nil {
		if subscription.IsTrial {
			trialDaysRemaining := int(time.Until(*subscription.TrialEndDate).Hours() / 24)
			if trialDaysRemaining < 0 {
				trialDaysRemaining = 0
			}
			trialInfo = gin.H{
				"status":          "active",
				"is_trial":        true,
				"days_remaining":  trialDaysRemaining,
				"expires_at":      subscription.TrialEndDate.Format("2006-01-02"),
				"subscription_id": subscription.ID,
			}
		} else {
			daysRemaining := int(time.Until(subscription.EndDate).Hours() / 24)
			if daysRemaining < 0 {
				daysRemaining = 0
			}
			trialInfo = gin.H{
				"status":          "paid",
				"is_trial":        false,
				"plan_type":       subscription.PlanType,
				"days_remaining":  daysRemaining,
				"expires_at":      subscription.EndDate.Format("2006-01-02"),
				"subscription_id": subscription.ID,
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":      user.ID,
			"email":   user.Email,
			"name":    user.Name,
			"isAdmin": user.IsAdmin,
			"trial":   trialInfo,
		},
	})
}

// Register page
func registerPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{})
}

// Register handler
func registerHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Name     string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Check if user exists
	var existingUser User
	if err := db.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Create user
	user := User{
		Email:    input.Email,
		Password: hashedPassword,
		Name:     input.Name,
		IsActive: true,
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Auto-create 7-day free trial subscription
	trialEndDate := time.Now().AddDate(0, 0, 7) // 7 days from now
	subscription := Subscription{
		UserID:          user.ID,
		PlanType:        "trial",
		Status:          "active",
		IsTrial:         true,
		TrialEndDate:    &trialEndDate,
		StartDate:       time.Now(),
		EndDate:         trialEndDate,
		Price:           0,
		PaymentMethod:   "free_trial",
		MaxProfiles:     3,          // Generous limit for trial
		MaxApplications: 999,        // Unlimited IPO applications for trial
	}

	if err := db.Create(&subscription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create trial subscription"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful! You have 7 days free trial access.",
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
			"trial": gin.H{
				"status":      "active",
				"days_remaining": 7,
				"expires_at":   trialEndDate.Format("2006-01-02"),
			},
		},
	})
}

// Pricing page - Dynamic Pricing with Future Roadmap
func pricingHandler(c *gin.Context) {
	// 2026: PREMIUM-ONLY TIER - ₹1,999 for 3 months (50% discount from original ₹3,999)
	premiumPrice := 1999
	announcementBanner := "🎉 2026 LAUNCH YEAR SPECIAL! One Premium Plan - ₹1,999 for 3 months!"
	
	c.HTML(http.StatusOK, "pricing.html", gin.H{
		"year": 2026,
		"currency": "NPR",
		"announcement": announcementBanner,
		"paymentMethods": []string{"eSewa", "Khalti", "Bank Transfer"},
		"plans": []gin.H{
			{
				"name":     "Premium",
				"nameNepali": "प्रिमियम",
				"price":    premiumPrice,
				"priceUSD": fmt.Sprintf("$%d", premiumPrice/75),
				"duration": "3 months",
				"durationNepali": "3 महीना",
				"popular":  true,
				"badge":    "Only Plan",
				"badgeNepali": "एकमात्र योजना",
				"description": "Everything you need for IPO automation",
				"descriptionNepali": "IPO स्वचालन के लिए आपको सभी कुछ",
				"features": []string{
					"✓ Unlimited MeroShare Accounts",
					"✓ Unlimited IPO Applications",
					"✓ Real-time IPO Notifications",
					"✓ 24/7 Priority Email & Chat Support",
					"✓ 2-minute Smart Monitoring",
					"✓ Multi-Source IPO Tracking (All Exchanges)",
					"✓ SMS Alerts for New IPOs",
					"✓ Secure Credential Encryption",
					"✓ Mobile-Friendly Dashboard",
				},
				"featuresNepali": []string{
					"✓ असीमित MeroShare खाताहरू",
					"✓ असीमित IPO आवेदनहरू",
					"✓ रिअल-टाइम IPO सूचनाहरू",
					"✓ 24/7 प्राथमिकता समर्थन",
					"✓ 2-मिनेट स्मार्ट निगरानी",
					"✓ बहु-स्रोत IPO ट्र्यैकिङ",
					"✓ SMS अलर्ट",
					"✓ सुरक्षित एन्क्रिप्शन",
					"✓ मोबाइल-अनुकूल डैशबोर्ड",
				},
			},
		},
	})
}

// Dashboard handler
func dashboardHandler(c *gin.Context) {
	userID := c.GetUint("userID")

	var stats DashboardStats
	
	// Get profile counts
	var totalProfiles, activeProfiles, totalApps, successApps, pendingApps int64
	db.Model(&Profile{}).Where("user_id = ?", userID).Count(&totalProfiles)
	db.Model(&Profile{}).Where("user_id = ? AND is_active = ?", userID, true).Count(&activeProfiles)
	
	// Get application counts
	db.Model(&IPOApplication{}).Where("user_id = ?", userID).Count(&totalApps)
	db.Model(&IPOApplication{}).Where("user_id = ? AND status = ?", userID, "success").Count(&successApps)
	db.Model(&IPOApplication{}).Where("user_id = ? AND status = ?", userID, "pending").Count(&pendingApps)
	
	stats.TotalProfiles = int(totalProfiles)
	stats.ActiveProfiles = int(activeProfiles)
	stats.TotalApplications = int(totalApps)
	stats.SuccessfulApps = int(successApps)
	stats.PendingApps = int(pendingApps)
	
	// Get subscription info
	var subscription Subscription
	if err := db.Where("user_id = ? AND status = ?", userID, "active").Order("end_date DESC").First(&subscription).Error; err == nil {
		stats.SubscriptionStatus = "Active"
		stats.SubscriptionExpiry = subscription.EndDate
		stats.RemainingDays = int(time.Until(subscription.EndDate).Hours() / 24)
		
		// Add trial info to response if it's a trial
		if subscription.IsTrial {
			c.Set("isTrial", true)
			c.Set("trialRemainingDays", stats.RemainingDays)
		}
	} else {
		stats.SubscriptionStatus = "Inactive"
	}

	// Get open IPOs count (from all sources)
	openIPOs, _ := getOpenIPOsFromAllSources()
	stats.OpenIPOs = len(openIPOs)

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"stats": stats,
	})
}

// Profiles handler
func profilesHandler(c *gin.Context) {
	userID := c.GetUint("userID")

	var profiles []Profile
	db.Where("user_id = ?", userID).Find(&profiles)

	c.HTML(http.StatusOK, "profiles.html", gin.H{
		"profiles": profiles,
	})
}

// Create profile handler
func createProfileHandler(c *gin.Context) {
	userID := c.GetUint("userID")

	var input struct {
		Name            string `json:"name" binding:"required"`
		DPID            string `json:"dpid" binding:"required"`
		BOID            string `json:"boid" binding:"required"`
		Password        string `json:"password" binding:"required"`
		CRN             string `json:"crn" binding:"required"`
		TransactionPIN  string `json:"transaction_pin" binding:"required"`
		DefaultBankID   int    `json:"default_bank_id" binding:"required"`
		DefaultKittas   int    `json:"default_kittas"`
		AskForKittas    bool   `json:"ask_for_kittas"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Check subscription limits
	var profileCount int64
	db.Model(&Profile{}).Where("user_id = ?", userID).Count(&profileCount)
	
	var subscription Subscription
	if err := db.Where("user_id = ? AND status = ?", userID, "active").First(&subscription).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "No active subscription found"})
		return
	}

	if int(profileCount) >= subscription.MaxProfiles {
		c.JSON(http.StatusForbidden, gin.H{"error": "Profile limit reached. Upgrade your plan."})
		return
	}

	// Encrypt sensitive data
	key := generateEncryptionKey(input.Name)
	passwordEnc, _ := encryptAES(key, input.Password)
	crnEnc, _ := encryptAES(key, input.CRN)
	pinEnc, _ := encryptAES(key, input.TransactionPIN)

	profile := Profile{
		UserID:            userID,
		Name:              input.Name,
		DPID:              input.DPID,
		BOID:              input.BOID,
		PasswordEnc:       passwordEnc,
		CRNEnc:            crnEnc,
		TransactionPINEnc: pinEnc,
		DefaultBankID:     input.DefaultBankID,
		DefaultKittas:     input.DefaultKittas,
		AskForKittas:      input.AskForKittas,
		IsActive:          true,
	}

	if err := db.Create(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Profile created successfully",
		"profile": profile,
	})
}

// Update profile handler
func updateProfileHandler(c *gin.Context) {
	userID := c.GetUint("userID")
	profileID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var profile Profile
	if err := db.Where("id = ? AND user_id = ?", profileID, userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Update allowed fields
	if err := db.Model(&profile).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"profile": profile,
	})
}

// Delete profile handler
func deleteProfileHandler(c *gin.Context) {
	userID := c.GetUint("userID")
	profileID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	result := db.Where("id = ? AND user_id = ?", profileID, userID).Delete(&Profile{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile deleted successfully"})
}

// IPOs handler
func iposHandler(c *gin.Context) {
	ipos, err := getOpenIPOsFromAllSources()
	if err != nil {
		c.HTML(http.StatusOK, "ipos.html", gin.H{
			"error": "Failed to fetch IPOs",
			"ipos":  []IPOData{},
		})
		return
	}

	c.HTML(http.StatusOK, "ipos.html", gin.H{
		"ipos": ipos,
	})
}

// Apply to IPO handler
func applyIPOHandler(c *gin.Context) {
	userID := c.GetUint("userID")
	ipoID := c.Param("ipo_id")

	var input struct {
		ProfileID uint `json:"profile_id" binding:"required"`
		Kittas    int  `json:"kittas" binding:"required,min=10"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Verify profile ownership
	var profile Profile
	if err := db.Where("id = ? AND user_id = ?", input.ProfileID, userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	// Check application limit
	var appCount int64
	db.Model(&IPOApplication{}).Where("user_id = ?", userID).Count(&appCount)
	
	var subscription Subscription
	db.Where("user_id = ? AND status = ?", userID, "active").First(&subscription)
	if int(appCount) >= subscription.MaxApplications {
		c.JSON(http.StatusForbidden, gin.H{"error": "Application limit reached"})
		return
	}

	// Apply to IPO (simulated - integrate with actual MeroShare API)
	application := IPOApplication{
		UserID:         userID,
		ProfileID:      input.ProfileID,
		CompanyShareID: ipoID,
		KittasApplied:  input.Kittas,
		BankID:         profile.DefaultBankID,
		Status:         "pending",
		AppliedAt:      time.Now(),
	}

	if err := db.Create(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create application"})
		return
	}

	// TODO: Actual IPO application logic here
	go processIPOApplication(&application, &profile)

	c.JSON(http.StatusOK, gin.H{
		"message":     "Application submitted successfully",
		"application": application,
	})
}

// Applications handler
func applicationsHandler(c *gin.Context) {
	userID := c.GetUint("userID")

	var applications []IPOApplication
	db.Preload("Profile").Where("user_id = ?", userID).Order("created_at DESC").Find(&applications)

	c.HTML(http.StatusOK, "applications.html", gin.H{
		"applications": applications,
	})
}

// Settings handler
func settingsHandler(c *gin.Context) {
	userID := c.GetUint("userID")

	var user User
	db.First(&user, userID)

	c.HTML(http.StatusOK, "settings.html", gin.H{
		"user": user,
	})
}

// Update settings handler
func updateSettingsHandler(c *gin.Context) {
	userID := c.GetUint("userID")

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := db.Model(&User{}).Where("id = ?", userID).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
}

// Payment webhook handler
func paymentWebhookHandler(c *gin.Context) {
	// TODO: Implement payment gateway webhook
	// This would be called by Stripe, PayPal, Khalti, eSewa, etc.
	
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Verify payment signature
	// Create/update subscription
	// Send confirmation email

	c.JSON(http.StatusOK, gin.H{"received": true})
}

// API docs handler
func apiDocsHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "api_docs.html", gin.H{})
}

// Helper function to process IPO application
func processIPOApplication(app *IPOApplication, profile *Profile) {
	// Decrypt credentials
	key := generateEncryptionKey(profile.Name)
	password, _ := decryptAES(key, profile.PasswordEnc)
	
	// Call MeroShare API
	success, message := applyToMeroShareIPO(profile, password, app.CompanyShareID, app.KittasApplied)
	
	if success {
		app.Status = "success"
		app.ResponseMsg = message
	} else {
		app.Status = "failed"
		app.ResponseMsg = message
	}
	
	db.Save(app)
}

// Placeholder for actual MeroShare API call
func applyToMeroShareIPO(profile *Profile, password string, shareID string, kittas int) (bool, string) {
	// TODO: Implement actual MeroShare API integration
	fmt.Printf("Applying to IPO: %s with %d kittas for profile %s\n", shareID, kittas, profile.Name)
	return true, "Application successful"
}
