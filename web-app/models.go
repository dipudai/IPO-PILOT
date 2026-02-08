package main

import (
	"time"

	"gorm.io/gorm"
)

// User represents a registered user
type User struct {
	gorm.Model
	Email           string         `gorm:"uniqueIndex;not null"`
	Password        string         `gorm:"not null"`
	Name            string         `gorm:"not null"`
	IsAdmin         bool           `gorm:"default:false"`
	IsActive        bool           `gorm:"default:true"`
	Subscriptions   []Subscription `gorm:"foreignKey:UserID"`
	Profiles        []Profile      `gorm:"foreignKey:UserID"`
	IPOApplications []IPOApplication `gorm:"foreignKey:UserID"`
}

// Subscription represents a user's subscription plan
type Subscription struct {
	gorm.Model
	UserID          uint      `gorm:"not null"`
	User            User      `gorm:"foreignKey:UserID"`
	PlanType        string    `gorm:"not null"` // basic, premium, enterprise
	Status          string    `gorm:"not null"` // active, expired, cancelled
	StartDate       time.Time `gorm:"not null"`
	EndDate         time.Time `gorm:"not null"`
	Price           float64   `gorm:"not null"`
	PaymentMethod   string    `gorm:"not null"`
	TransactionID   string    
	MaxProfiles     int       `gorm:"default:1"`
	MaxApplications int       `gorm:"default:100"`
}

// Profile represents a MeroShare account profile
type Profile struct {
	gorm.Model
	UserID          uint   `gorm:"not null"`
	User            User   `gorm:"foreignKey:UserID"`
	Name            string `gorm:"not null"`
	DPID            string `gorm:"not null"`
	BOID            string `gorm:"not null"`
	PasswordEnc     string `gorm:"not null"` // Encrypted
	CRNEnc          string `gorm:"not null"` // Encrypted
	TransactionPINEnc string `gorm:"not null"` // Encrypted
	DefaultBankID   int    `gorm:"not null"`
	DefaultKittas   int    `gorm:"default:10"`
	AskForKittas    bool   `gorm:"default:false"`
	IsActive        bool   `gorm:"default:true"`
	LastUsed        *time.Time
}

// IPOApplication represents an application made to an IPO
type IPOApplication struct {
	gorm.Model
	UserID         uint      `gorm:"not null"`
	User           User      `gorm:"foreignKey:UserID"`
	ProfileID      uint      `gorm:"not null"`
	Profile        Profile   `gorm:"foreignKey:ProfileID"`
	IPOSourceID    uint      `gorm:"not null"`
	IPOSource      IPOSource `gorm:"foreignKey:IPOSourceID"`
	CompanyName    string    `gorm:"not null"`
	CompanyShareID string    `gorm:"not null"`
	KittasApplied  int       `gorm:"not null"`
	BankID         int       `gorm:"not null"`
	Status         string    `gorm:"not null"` // pending, success, failed
	AppliedAt      time.Time `gorm:"not null"`
	ResponseMsg    string
}

// IPOSource represents an IPO data source
type IPOSource struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Type        string `gorm:"not null"` // meroshare, iporesult, cts, custom
	BaseURL     string `gorm:"not null"`
	APIKey      string
	IsActive    bool   `gorm:"default:true"`
	Priority    int    `gorm:"default:0"` // Higher priority checked first
	LastChecked *time.Time
	Description string
}

// MonitoringSession tracks active monitoring sessions
type MonitoringSession struct {
	gorm.Model
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID"`
	ProfileID uint      `gorm:"not null"`
	Profile   Profile   `gorm:"foreignKey:ProfileID"`
	IsActive  bool      `gorm:"default:true"`
	StartedAt time.Time `gorm:"not null"`
	StoppedAt *time.Time
	Interval  int       `gorm:"default:300"` // seconds
}

// IPOData represents IPO information from various sources
type IPOData struct {
	SourceID        uint      `json:"source_id"`
	CompanyName     string    `json:"company_name"`
	StockSymbol     string    `json:"stock_symbol"`
	CompanyShareID  string    `json:"company_share_id"`
	SectorName      string    `json:"sector_name"`
	StockPrice      string    `json:"stock_price"`
	MinUnits        string    `json:"min_units"`
	MaxUnits        string    `json:"max_units"`
	TotalUnits      string    `json:"total_units"`
	IssueOpenDate   string    `json:"issue_open_date"`
	IssueCloseDate  string    `json:"issue_close_date"`
	Status          string    `json:"status"`
	ShareType       string    `json:"share_type"`
	ShareGroup      string    `json:"share_group"`
	AppliedKittas   int       `json:"applied_kittas"`
	LastUpdated     time.Time `json:"last_updated"`
}

// DashboardStats for dashboard view
type DashboardStats struct {
	TotalProfiles       int
	ActiveProfiles      int
	TotalApplications   int
	SuccessfulApps      int
	PendingApps         int
	OpenIPOs            int
	SubscriptionStatus  string
	SubscriptionExpiry  time.Time
	RemainingDays       int
}
