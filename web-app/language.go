package main

import (
	"net/http"
	"os"
)

// Language management for multi-language support
type LanguageString struct {
	English string
	Nepali  string
}

// Get user's preferred language from cookie or header
func getUserLanguage(r *http.Request) string {
	// Check cookie first
	if lang, err := r.Cookie("language"); err == nil {
		if lang.Value == "nepali" {
			return "nepali"
		}
	}

	// Check Accept-Language header
	acceptLang := r.Header.Get("Accept-Language")
	if len(acceptLang) > 0 {
		// If user's browser is set to Nepali, suggest Nepali
		if acceptLang[:2] == "ne" {
			return "nepali"
		}
	}

	// Default to English
	return "english"
}

// Set language preference
func setLanguage(w http.ResponseWriter, lang string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "language",
		Value:    lang,
		MaxAge:   86400 * 365, // 1 year
		Path:     "/",
		HttpOnly: true,
	})
}

// Get environment type (for pricing roadmap)
func getEnvironmentMode() string {
	mode := os.Getenv("APP_MODE")
	if mode == "" {
		mode = "production"
	}
	return mode
}

// Pricing tiers based on year and environment
type PricingTier struct {
	Name            string
	NameNepali      string
	PriceNPR        int
	PriceUSD        int
	Duration        string
	DurationNepali  string
	IsPopular       bool
	Description     string
	DescriptionNepali string
	Features        []string
	FeaturesNepali  []string
}

// Get pricing based on year - PREMIUM ONLY
func getPricingForYear(year int) (premiumPrice int) {
	// 2026 LAUNCH: PREMIUM-ONLY - 50% Discount!
	return 1999  // Single tier: ₹1,999 (was ₹3,999)
}

// Get pricing message for current year
func getPricingMessage(year int) string {
	return "🎉 2026 LAUNCH YEAR SPECIAL! One Premium Plan - ₹1,999 for 3 months!"
}
