package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Admin dashboard handler
func adminDashboardHandler(c *gin.Context) {
	var stats struct {
		TotalUsers          int64
		ActiveUsers         int64
		TotalSubscriptions  int64
		ActiveSubscriptions int64
		TotalApplications   int64
		TotalRevenue        float64
		IPOSources          int64
	}

	db.Model(&User{}).Count(&stats.TotalUsers)
	db.Model(&User{}).Where("is_active = ?", true).Count(&stats.ActiveUsers)
	db.Model(&Subscription{}).Count(&stats.TotalSubscriptions)
	db.Model(&Subscription{}).Where("status = ?", "active").Count(&stats.ActiveSubscriptions)
	db.Model(&IPOApplication{}).Count(&stats.TotalApplications)
	db.Model(&Subscription{}).Select("COALESCE(SUM(price), 0)").Row().Scan(&stats.TotalRevenue)
	db.Model(&IPOSource{}).Where("is_active = ?", true).Count(&stats.IPOSources)

	c.HTML(http.StatusOK, "admin_dashboard.html", gin.H{
		"stats": stats,
	})
}

// Admin users handler
func adminUsersHandler(c *gin.Context) {
	var users []User
	db.Preload("Subscriptions").Find(&users)

	c.HTML(http.StatusOK, "admin_users.html", gin.H{
		"users": users,
	})
}

// Admin subscriptions handler
func adminSubscriptionsHandler(c *gin.Context) {
	var subscriptions []Subscription
	db.Preload("User").Order("created_at DESC").Find(&subscriptions)

	c.HTML(http.StatusOK, "admin_subscriptions.html", gin.H{
		"subscriptions": subscriptions,
	})
}

// Activate subscription handler
func activateSubscriptionHandler(c *gin.Context) {
	subscriptionID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var subscription Subscription
	if err := db.First(&subscription, subscriptionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	subscription.Status = "active"
	subscription.StartDate = time.Now()
	
	// Set end date based on plan
	months := 3
	if subscription.PlanType == "enterprise" {
		months = 12
	}
	subscription.EndDate = time.Now().AddDate(0, months, 0)

	if err := db.Save(&subscription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Subscription activated",
		"subscription": subscription,
	})
}

// Deactivate subscription handler
func deactivateSubscriptionHandler(c *gin.Context) {
	subscriptionID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var subscription Subscription
	if err := db.First(&subscription, subscriptionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	subscription.Status = "cancelled"
	if err := db.Save(&subscription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription deactivated"})
}

// IPO sources handler
func ipoSourcesHandler(c *gin.Context) {
	var sources []IPOSource
	db.Order("priority DESC").Find(&sources)

	c.HTML(http.StatusOK, "admin_ipo_sources.html", gin.H{
		"sources": sources,
	})
}

// Add IPO source handler
func addIPOSourceHandler(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Type        string `json:"type" binding:"required"`
		BaseURL     string `json:"base_url" binding:"required"`
		APIKey      string `json:"api_key"`
		Priority    int    `json:"priority"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	source := IPOSource{
		Name:        input.Name,
		Type:        input.Type,
		BaseURL:     input.BaseURL,
		APIKey:      input.APIKey,
		Priority:    input.Priority,
		IsActive:    true,
		Description: input.Description,
	}

	if err := db.Create(&source).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add IPO source"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "IPO source added successfully",
		"source":  source,
	})
}

// Delete IPO source handler
func deleteIPOSourceHandler(c *gin.Context) {
	sourceID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := db.Delete(&IPOSource{}, sourceID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete IPO source"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "IPO source deleted successfully"})
}

// Analytics handler
func analyticsHandler(c *gin.Context) {
	// Get daily/weekly/monthly stats
	var analyticsData struct {
		DailyApplications  []map[string]interface{}
		WeeklyRevenue      []map[string]interface{}
		PopularIPOs        []map[string]interface{}
		UserGrowth         []map[string]interface{}
		SubscriptionBreakdown map[string]int64
	}

	// Daily applications (last 30 days)
	db.Model(&IPOApplication{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ?", time.Now().AddDate(0, 0, -30)).
		Group("DATE(created_at)").
		Order("date").
		Scan(&analyticsData.DailyApplications)

	// Subscription breakdown - PREMIUM ONLY
	analyticsData.SubscriptionBreakdown = make(map[string]int64)
	var premiumCount int64
	db.Model(&Subscription{}).Where("plan_type = ? AND status = ?", "premium", "active").Count(&premiumCount)
	analyticsData.SubscriptionBreakdown["premium"] = premiumCount

	c.HTML(http.StatusOK, "admin_analytics.html", gin.H{
		"analytics": analyticsData,
	})
}
