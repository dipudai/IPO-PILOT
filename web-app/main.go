package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:embed templates/* static/*
var embedFS embed.FS

var db *gorm.DB

func main() {
	// Initialize database
	var err error
	db, err = gorm.Open(sqlite.Open("ipo_pilot.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate database schema
	db.AutoMigrate(&User{}, &Subscription{}, &Profile{}, &IPOApplication{}, &IPOSource{})

	// Initialize default admin user
	initializeAdmin()

	// Setup router
	r := gin.Default()

	// Load templates from embedded FS (reliable in Docker)
	r.SetHTMLTemplate(template.Must(template.ParseFS(embedFS, "templates/*.html")))
	
	// Serve static files from embedded FS
	r.StaticFS("/static", http.FS(embedFS))

	// Public routes
	r.GET("/", homeHandler)
	r.GET("/login", loginPageHandler)
	r.POST("/login", loginHandler)
	r.GET("/register", registerPageHandler)
	r.POST("/register", registerHandler)
	r.GET("/pricing", pricingHandler)

	// Language toggle
	r.GET("/set-language/:lang", func(c *gin.Context) {
		lang := c.Param("lang")
		if lang != "english" && lang != "nepali" {
			lang = "english"
		}
		setLanguage(c.Writer, lang)
		c.Redirect(http.StatusFound, c.GetHeader("Referer"))
	})

	// User dashboard routes (require authentication)
	user := r.Group("/dashboard")
	user.Use(authMiddleware())
	{
		user.GET("", dashboardHandler)
		user.GET("/profiles", profilesHandler)
		user.POST("/profiles", createProfileHandler)
		user.PUT("/profiles/:id", updateProfileHandler)
		user.DELETE("/profiles/:id", deleteProfileHandler)
		user.GET("/ipos", iposHandler)
		user.POST("/apply/:ipo_id", applyIPOHandler)
		user.GET("/applications", applicationsHandler)
		user.GET("/settings", settingsHandler)
		user.POST("/settings", updateSettingsHandler)
	}

	// Admin routes
	admin := r.Group("/admin")
	admin.Use(authMiddleware(), adminMiddleware())
	{
		admin.GET("", adminDashboardHandler)
		admin.GET("/users", adminUsersHandler)
		admin.GET("/subscriptions", adminSubscriptionsHandler)
		admin.POST("/subscriptions/:id/activate", activateSubscriptionHandler)
		admin.POST("/subscriptions/:id/deactivate", deactivateSubscriptionHandler)
		admin.GET("/ipo-sources", ipoSourcesHandler)
		admin.POST("/ipo-sources", addIPOSourceHandler)
		admin.DELETE("/ipo-sources/:id", deleteIPOSourceHandler)
		admin.GET("/analytics", analyticsHandler)
	}

	// API routes for AJAX requests
	api := r.Group("/api")
	api.Use(authMiddleware())
	{
		api.GET("/ipos/live", getLiveIPOsHandler)
		api.GET("/ipos/upcoming", getUpcomingIPOsHandler)
		api.POST("/monitor/start", startMonitoringHandler)
		api.POST("/monitor/stop", stopMonitoringHandler)
		api.GET("/monitor/status", monitorStatusHandler)
	}

	// Payment webhook
	r.POST("/webhook/payment", paymentWebhookHandler)

	// Nepal Payment Gateways
	payment := r.Group("/payment")
	{
		payment.POST("/nepal", nepaliPaymentHandler)
		payment.GET("/esewa/success", esewaSuccessHandler)
		payment.GET("/esewa/failure", func(c *gin.Context) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Payment failed"})
		})
		payment.POST("/khalti/success", khaltiSuccessHandler)
		payment.GET("/connectips", connectIPSHandler)
	}

	// API documentation
	r.GET("/api/docs", apiDocsHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("\n🚀 IPO Pilot Web Platform Starting...\n")
	fmt.Printf("📱 URL: http://localhost:%s\n", port)
	fmt.Printf("👤 Default Admin: admin@ipopilot.com / admin123\n\n")

	r.Run(":" + port)
}

func initializeAdmin() {
	var admin User
	result := db.Where("email = ?", "admin@ipopilot.com").First(&admin)
	if result.Error == gorm.ErrRecordNotFound {
		hashedPassword, _ := hashPassword("admin123")
		admin = User{
			Email:    "admin@ipopilot.com",
			Password: hashedPassword,
			Name:     "Administrator",
			IsAdmin:  true,
			IsActive: true,
		}
		db.Create(&admin)
		log.Println("✓ Default admin user created")
	}
}
