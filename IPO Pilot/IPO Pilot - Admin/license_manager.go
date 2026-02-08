package main

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type LicenseRecord struct {
	Key        string    `json:"key"`
	IssuedTo   string    `json:"issued_to"`
	IssuedDate time.Time `json:"issued_date"`
	Status     string    `json:"status"` // "active", "revoked", "expired"
	Notes      string    `json:"notes"`
}

type LicenseManager struct {
	Licenses []LicenseRecord `json:"licenses"`
}

const LICENSE_DB = "license_manager.json"

func generateLicenseKey() string {
	// Generate a license key with checksum for validation
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 19) // 4-4-4-4 format with dashes

	// Generate random part
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

	parts := strings.Split(key, "-")
	if len(parts) != 4 {
		return false
	}

	for _, part := range parts {
		if len(part) != 4 {
			return false
		}
	}

	// Validate checksum
	keyWithoutChecksum := key[:15] // First 15 chars (XXXX-XXXX-XXXX)
	expectedChecksum := fmt.Sprintf("%X", md5.Sum([]byte(keyWithoutChecksum)))[:4]
	actualChecksum := key[15:] // Last 4 chars

	return strings.ToUpper(expectedChecksum) == strings.ToUpper(actualChecksum)
}

func loadLicenseManager() (*LicenseManager, error) {
	if _, err := os.Stat(LICENSE_DB); os.IsNotExist(err) {
		return &LicenseManager{Licenses: []LicenseRecord{}}, nil
	}

	data, err := os.ReadFile(LICENSE_DB)
	if err != nil {
		return nil, err
	}

	var lm LicenseManager
	err = json.Unmarshal(data, &lm)
	return &lm, err
}

func (lm *LicenseManager) save() error {
	data, err := json.MarshalIndent(lm, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(LICENSE_DB, data, 0644)
}

func (lm *LicenseManager) generateNewLicense(issuedTo, notes string) string {
	key := generateLicenseKey()

	// Ensure unique key
	for lm.isKeyExists(key) {
		key = generateLicenseKey()
	}

	record := LicenseRecord{
		Key:        key,
		IssuedTo:   issuedTo,
		IssuedDate: time.Now(),
		Status:     "active",
		Notes:      notes,
	}

	lm.Licenses = append(lm.Licenses, record)
	return key
}

func (lm *LicenseManager) isKeyExists(key string) bool {
	for _, license := range lm.Licenses {
		if license.Key == key {
			return true
		}
	}
	return false
}

func (lm *LicenseManager) revokeLicense(key string) bool {
	for i, license := range lm.Licenses {
		if license.Key == key {
			lm.Licenses[i].Status = "revoked"
			return true
		}
	}
	return false
}

func (lm *LicenseManager) listLicenses() {
	fmt.Println("\n=== LICENSE MANAGER ===")
	fmt.Println("Total licenses issued:", len(lm.Licenses))

	if len(lm.Licenses) == 0 {
		fmt.Println("No licenses found.")
		return
	}

	fmt.Printf("\n%-25s %-20s %-12s %-10s %s\n", "LICENSE KEY", "ISSUED TO", "DATE", "STATUS", "NOTES")
	fmt.Println(strings.Repeat("-", 100))

	for _, license := range lm.Licenses {
		fmt.Printf("%-25s %-20s %-12s %-10s %s\n",
			license.Key,
			truncateString(license.IssuedTo, 20),
			license.IssuedDate.Format("2006-01-02"),
			license.Status,
			license.Notes)
	}
}

func truncateString(str string, maxLen int) string {
	if len(str) > maxLen {
		return str[:maxLen-3] + "..."
	}
	return str
}

func showMenu() {
	fmt.Println("\n=== IPO Pilot License Manager ===")
	fmt.Println("1. Generate new license")
	fmt.Println("2. List all licenses")
	fmt.Println("3. Revoke license")
	fmt.Println("4. Validate license key")
	fmt.Println("5. Exit")
	fmt.Print("Choose an option (1-5): ")
}

func main() {
	rand.Seed(time.Now().UnixNano())

	lm, err := loadLicenseManager()
	if err != nil {
		fmt.Println("Error loading license database:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		showMenu()
		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			fmt.Print("Enter customer name/email: ")
			issuedTo, _ := reader.ReadString('\n')
			issuedTo = strings.TrimSpace(issuedTo)

			fmt.Print("Enter notes (optional): ")
			notes, _ := reader.ReadString('\n')
			notes = strings.TrimSpace(notes)

			key := lm.generateNewLicense(issuedTo, notes)
			err := lm.save()
			if err != nil {
				fmt.Println("Error saving license:", err)
				continue
			}

			fmt.Println("\n✅ LICENSE GENERATED SUCCESSFULLY!")
			fmt.Println("License Key:", key)
			fmt.Println("Valid for: 3 months from activation")
			fmt.Println("Customer:", issuedTo)
			fmt.Println("\n⚠️  IMPORTANT: Save this key securely!")

		case "2":
			lm.listLicenses()

		case "3":
			fmt.Print("Enter license key to revoke: ")
			key, _ := reader.ReadString('\n')
			key = strings.TrimSpace(key)

			if lm.revokeLicense(key) {
				err := lm.save()
				if err != nil {
					fmt.Println("Error saving changes:", err)
				} else {
					fmt.Println("License revoked successfully!")
				}
			} else {
				fmt.Println("License key not found!")
			}

		case "4":
			fmt.Print("Enter license key to validate: ")
			key, _ := reader.ReadString('\n')
			key = strings.TrimSpace(key)

			if validateLicenseKey(key) {
				// Check if it's in our database
				found := false
				for _, license := range lm.Licenses {
					if license.Key == key {
						fmt.Printf("✅ Valid license key\nStatus: %s\nIssued to: %s\nIssued: %s\n",
							license.Status, license.IssuedTo, license.IssuedDate.Format("2006-01-02"))
						found = true
						break
					}
				}
				if !found {
					fmt.Println("✅ Valid license key format (not in database)")
				}
			} else {
				fmt.Println("❌ Invalid license key format")
			}

		case "5":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
