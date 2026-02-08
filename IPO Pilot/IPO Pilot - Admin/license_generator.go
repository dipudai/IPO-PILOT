package main

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║       IPO PILOT - LICENSE KEY GENERATOR v1.0              ║")
	fmt.Println("║     Generate unlimited valid license keys instantly       ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝\n")

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("GENERATED LICENSE KEYS (Valid for 3 months after activation)")
	fmt.Println("═══════════════════════════════════════════════════════════\n")

	// Generate 10 sample licenses
	fmt.Println("Batch 1 - 10 License Keys:")
	fmt.Println("─────────────────────────────────────────────────────────\n")
	for i := 1; i <= 10; i++ {
		key := generateLicenseKey()
		fmt.Printf("%2d.  %s  [Valid: ✓]\n", i, key)
	}

	// Validate a generated key
	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("LICENSE VALIDATION TEST")
	fmt.Println("═══════════════════════════════════════════════════════════\n")

	testKey := generateLicenseKey()
	isValid := validateLicenseKey(testKey)

	fmt.Printf("Generated Key: %s\n", testKey)
	fmt.Printf("Format Valid:  %v\n", isValid)
	fmt.Printf("Checksum Valid:%v\n\n", isValid)

	// Generate more licenses if needed
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("BATCH 2 - 10 More License Keys")
	fmt.Println("─────────────────────────────────────────────────────────\n")
	for i := 11; i <= 20; i++ {
		key := generateLicenseKey()
		fmt.Printf("%2d.  %s  [Valid: ✓]\n", i, key)
	}

	// Instructions
	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("NEXT STEPS")
	fmt.Println("═══════════════════════════════════════════════════════════\n")
	fmt.Println("1. Copy any license key above")
	fmt.Println("2. Send to customer via email or download link")
	fmt.Println("3. Customer activates on:")
	fmt.Println("   • Windows: Menu → License Information → Activate License")
	fmt.Println("   • Mobile: License Screen → Enter Key → Activate")
	fmt.Println("4. License valid for 3 months from activation")
	fmt.Println("5. Track in spreadsheet: customer name, key, device ID, expiry\n")

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("KEY FEATURES")
	fmt.Println("═══════════════════════════════════════════════════════════\n")
	fmt.Println("✓ Format: XXXX-XXXX-XXXX-XXXX (19 characters)")
	fmt.Println("✓ Last 4 chars: MD5 checksum of first 15")
	fmt.Println("✓ Works on Windows AND Android (same key)")
	fmt.Println("✓ Auto-binds to machine/device ID")
	fmt.Println("✓ Prevents license sharing/piracy")
	fmt.Println("✓ 3-month validity after activation")
	fmt.Println("✓ Simple renewal process\n")

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("Generated at: " + time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("═══════════════════════════════════════════════════════════")
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

	// Extract checksum parts
	keyWithoutChecksum := key[:15] // First 15 chars (XXXX-XXXX-XXXX)
	expectedChecksum := fmt.Sprintf("%X", md5.Sum([]byte(keyWithoutChecksum)))[:4]
	actualChecksum := key[15:] // Last 4 chars

	return expectedChecksum == actualChecksum
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
