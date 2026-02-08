package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Command line flags
	count := flag.Int("count", 1, "Number of licenses to generate")
	validate := flag.String("validate", "", "Validate a license key")
	flag.Parse()

	fmt.Println("\n╔════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    🔑 IPO Pilot License Generator                         ║")
	fmt.Println("║                                                                            ║")
	fmt.Println("║   Generate secure license keys with built-in MD5 checksum validation      ║")
	fmt.Println("║   Valid for both Windows Desktop and Android Mobile versions              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════════════╝\n")

	// If validating a key
	if *validate != "" {
		fmt.Printf("Validating: %s\n", *validate)
		if validateLicenseKey(*validate) {
			fmt.Println("✓ License key is VALID")
			fmt.Println("\nKey Details:")
			fmt.Printf("  Base: %s\n", (*validate)[:15])
			fmt.Printf("  Checksum: %s\n", (*validate)[15:])
			fmt.Printf("  Format: Valid\n")
			fmt.Printf("  Validity: 3 months from activation\n")
			fmt.Printf("  Binding: Machine/Device ID\n")
		} else {
			fmt.Println("✗ License key is INVALID")
			fmt.Println("\nPossible reasons:")
			fmt.Println("  • Wrong format (must be XXXX-XXXX-XXXX-XXXX)")
			fmt.Println("  • Incorrect checksum")
			fmt.Println("  • Must be uppercase letters and numbers")
		}
		return
	}

	// Generate licenses
	fmt.Printf("Generating %d license key(s)...\n\n", *count)

	fmt.Println("╔════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║ License Key                 │ Valid │ Activation Expires                 ║")
	fmt.Println("╠════════════════════════════════════════════════════════════════════════════╣")

	for i := 1; i <= *count; i++ {
		key := generateLicenseKey()
		expiryDate := time.Now().AddDate(0, 3, 0).Format("2006-01-02")
		fmt.Printf("║ %-27s │   ✓   │ From today to %s         ║\n", key, expiryDate)
	}

	fmt.Println("╚════════════════════════════════════════════════════════════════════════════╝\n")

	fmt.Println("📋 Instructions for customers:\n")
	fmt.Println("  Windows Version:")
	fmt.Println("    1. Run: go run ipo_master_customer.go")
	fmt.Println("    2. Select: License Information")
	fmt.Println("    3. Choose: Activate License")
	fmt.Println("    4. Enter: License key (UPPERCASE)")
	fmt.Println("\n  Android Version:")
	fmt.Println("    1. Open IPO Pilot app")
	fmt.Println("    2. Go to: License Screen")
	fmt.Println("    3. Enter: License key (UPPERCASE)")
	fmt.Println("    4. Tap: Activate License\n")

	fmt.Println("✨ Features:")
	fmt.Println("  • MD5 checksum validation (prevents typos)")
	fmt.Println("  • Hardware machine ID binding (prevents sharing)")
	fmt.Println("  • 3-month validity period")
	fmt.Println("  • Works on both Windows and Android")
	fmt.Println("  • Same license key for all platforms\n")
}

// generateLicenseKey creates a new license key with MD5 checksum
func generateLicenseKey() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 19) // 4-4-4-4 format with dashes

	// Generate random part (16 characters without dashes)
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

// validateLicenseKey checks if a license key is valid
func validateLicenseKey(key string) bool {
	// Check length
	if len(key) != 19 {
		return false
	}

	// Check format (XXXX-XXXX-XXXX-XXXX)
	parts := make([]string, 4)
	start := 0
	for i := 0; i < 4; i++ {
		end := start + 4
		if i < 3 && (end+1 > len(key) || key[end] != '-') {
			return false
		}
		if i == 3 {
			parts[i] = key[start:]
		} else {
			parts[i] = key[start : start+4]
		}

		// Validate each part contains only uppercase letters and numbers
		for _, ch := range parts[i] {
			if !((ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')) {
				return false
			}
		}
		start = end + 1
	}

	// Validate checksum using MD5
	keyWithoutChecksum := key[:15] // First 15 chars (XXXX-XXXX-XXXX)
	expectedChecksum := fmt.Sprintf("%X", md5.Sum([]byte(keyWithoutChecksum)))[:4]
	actualChecksum := key[15:] // Last 4 chars

	return expectedChecksum == actualChecksum
}
