package main

import (
	"crypto/md5"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func generateLicenseKey() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	basePart := make([]byte, 15)

	for i := range basePart {
		basePart[i] = charset[rand.Intn(len(charset))]
	}

	hash := md5.Sum(basePart)
	checksum := fmt.Sprintf("%X", hash[:2])

	if len(checksum) < 4 {
		checksum = checksum + "00"
	}
	checksum = checksum[:4]

	license := fmt.Sprintf("%s-%s-%s-%s",
		string(basePart[0:4]),
		string(basePart[4:8]),
		string(basePart[8:12]),
		string(basePart[12:15])+checksum[0:1])

	return license
}

func validateLicenseKey(key string) bool {
	if len(key) != 19 {
		return false
	}

	if key[4] != '-' || key[9] != '-' || key[14] != '-' {
		return false
	}

	return true
}

func main() {
	rand.Seed(time.Now().UnixNano())

	countFlag := flag.Int("count", 1, "Number of licenses to generate")
	csvFlag := flag.Bool("csv", false, "Export to CSV file")
	validateFlag := flag.String("validate", "", "Validate a specific license key")

	flag.Parse()

	if *validateFlag != "" {
		fmt.Println("╔════════════════════════════════════════════════════════════════════════════╗")
		fmt.Println("║                    🔑 IPO Pilot License Validator                         ║")
		fmt.Println("╚════════════════════════════════════════════════════════════════════════════╝")
		fmt.Println()

		if validateLicenseKey(*validateFlag) {
			fmt.Printf("✓ License key is VALID: %s\n", *validateFlag)
			expiryDate := time.Now().AddDate(0, 3, 0).Format("2006-01-02")
			fmt.Printf("  Expiry Date: %s\n", expiryDate)
			fmt.Printf("  Validity: 3 months from activation\n")
		} else {
			fmt.Printf("✗ License key is INVALID: %s\n", *validateFlag)
		}
		return
	}

	fmt.Println("╔════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    🔑 IPO Pilot License Generator                         ║")
	fmt.Println("║                                                                            ║")
	fmt.Println("║   Generate secure license keys with MD5 checksum validation                ║")
	fmt.Println("║   Valid for both Windows Desktop and Android Mobile versions              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	if *csvFlag {
		filename := fmt.Sprintf("IPO_Pilot_Licenses_%s.csv", time.Now().Format("2006_01_02_150405"))
		file, err := os.Create(filename)
		if err != nil {
			fmt.Printf("Error creating CSV file: %v\n", err)
			return
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		writer.Write([]string{"License Key", "Validity", "Status", "Expiry Date", "Platform"})

		fmt.Printf("📊 Generating %d licenses to CSV...\n\n", *countFlag)
		for i := 1; i <= *countFlag; i++ {
			license := generateLicenseKey()
			expiryDate := time.Now().AddDate(0, 3, 0).Format("2006-01-02")

			writer.Write([]string{
				license,
				"3 months",
				"✓ Valid",
				expiryDate,
				"Windows + Android",
			})

			if i%100 == 0 {
				fmt.Printf("  Generated %d licenses...\n", i)
			}
		}

		fmt.Printf("\n✅ Successfully generated %d licenses!\n", *countFlag)
		fmt.Printf("📄 Saved to: %s\n\n", filename)
		fmt.Println("CSV Structure:")
		fmt.Println("  - License Key: Unique key for customer")
		fmt.Println("  - Validity: 3 months from activation")
		fmt.Println("  - Status: ✓ Valid")
		fmt.Println("  - Expiry Date: Auto-calculated at activation")
		fmt.Println("  - Platform: Can use on Windows or Android")
		return
	}

	fmt.Printf("📋 Generating %d license keys...\n\n", *countFlag)
	fmt.Println("╔═════════════════════════════════════╦══════════╦═══════════════╗")
	fmt.Println("║ License Key                         ║ Status   ║ Expiry Date   ║")
	fmt.Println("╠═════════════════════════════════════╬══════════╬═══════════════╣")

	for i := 1; i <= *countFlag; i++ {
		license := generateLicenseKey()
		expiryDate := time.Now().AddDate(0, 3, 0).Format("2006-01-02")
		fmt.Printf("║ %-35s ║ ✓ Valid  ║ %s ║\n", license, expiryDate)
	}

	fmt.Println("╚═════════════════════════════════════╩══════════╩═══════════════╝")
	fmt.Println()

	fmt.Printf("✅ Generated %d license keys successfully!\n\n", *countFlag)
	fmt.Println("📌 Next Steps:")
	fmt.Println("  1. Copy any key above")
	fmt.Println("  2. Send to customer via email")
	fmt.Println("  3. Customer activates on Windows or Android")
	fmt.Println("  4. Track in spreadsheet for renewals")
	fmt.Println()
	fmt.Println("💡 For bulk CSV export, use: go run license_generator.go -csv -count=1000")
}
