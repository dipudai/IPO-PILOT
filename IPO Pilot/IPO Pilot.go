package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	crand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/fatih/color"
)

var (
	red    = color.New(color.FgRed).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	blue   = color.New(color.FgBlue).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
)

type Config struct {
	Name           string
	DPID           string
	BOID           string
	Password       string
	CRN            string
	TransactionPIN string
	DefaultBankID  int
	DefaultKittas  int
	AskForKittas   bool
}

type BankBrief struct {
	Code string
	Id   int
	Name string
}

type License struct {
	Key         string
	MachineID   string
	ActivatedAt time.Time
	ExpiryDate  time.Time
}

type IPOData struct {
	CompanyName     string      `json:"companyName"`
	StockSymbol     string      `json:"stockSymbol"`
	CompanyImageURL interface{} `json:"companyImageUrl"`
	SectorName      string      `json:"sectorName"`
	StockPrice      string      `json:"stockPrice"`
	TotalCapital    string      `json:"totalCapital"`
	MinUnits        string      `json:"minUnits"`
	MaxUnits        string      `json:"maxUnits"`
	TotalUnits      string      `json:"totalUnits"`
	IssueOpenDate   string      `json:"issueOpenDate"`
	IssueCloseDate  string      `json:"issueCloseDate"`
	Status          int         `json:"status"`
	ShareType       string      `json:"shareType"`
	ShareGroup      string      `json:"shareGroup"`
	AppliedKittas   int         `json:"appliedKittas"`
}

type ApplicationResponse struct {
	Object struct {
		Message string `json:"message"`
	} `json:"object"`
}

func Log(message string) {
	fmt.Println(cyan("[") + time.Now().Format("15:04:05") + cyan("]") + " " + message)
}

func LogLevel(level, message string) {
	var colorFunc func(...interface{}) string
	switch level {
	case "ERROR":
		colorFunc = red
	case "SUCCESS":
		colorFunc = green
	case "WARNING":
		colorFunc = yellow
	case "INFO":
		colorFunc = blue
	default:
		colorFunc = fmt.Sprint
	}
	fmt.Println(cyan("[") + time.Now().Format("15:04:05") + cyan("]") + " " + colorFunc("["+level+"]") + " " + message)
}

func encrypt(text, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(crand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func decrypt(cryptoText, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	cipherText, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}
	if len(cipherText) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)
	return string(cipherText), nil
}

func generateKey(name string) string {
	hash := md5.Sum([]byte(name + "IPO~Master"))
	return fmt.Sprintf("%x", hash)[:16]
}

func validateLicenseKey(key string) bool {
	// Validate license key format and checksum
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

	// Validate checksum using MD5
	keyWithoutChecksum := key[:15] // First 15 chars (XXXX-XXXX-XXXX)
	expectedChecksum := fmt.Sprintf("%X", md5.Sum([]byte(keyWithoutChecksum)))[:4]
	actualChecksum := key[15:] // Last 4 chars

	return strings.ToUpper(expectedChecksum) == strings.ToUpper(actualChecksum)
}

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

func checkLicense() bool {
	licenseFile := "license.dat"

	// Check if license file exists
	if _, err := os.Stat(licenseFile); os.IsNotExist(err) {
		return false
	}

	// Read license file
	licenseData, err := os.ReadFile(licenseFile)
	if err != nil {
		return false
	}

	var license License
	err = json.Unmarshal(licenseData, &license)
	if err != nil {
		return false
	}

	// Validate license key format
	if !validateLicenseKey(license.Key) {
		return false
	}

	// Get current machine ID
	currentMachineID, err := machineid.ID()
	if err != nil {
		return false
	}

	// Check if machine ID matches
	if license.MachineID != currentMachineID {
		return false
	}

	// Check if license is expired
	if time.Now().After(license.ExpiryDate) {
		return false
	}

	return true
}

func activateLicense() bool {
	fmt.Println(cyan("\n=== License Activation ==="))
	fmt.Println("Please enter your license key to activate IPO~Master")
	fmt.Println("Contact the developer to purchase a license key.")
	fmt.Println("Email: planearn01@gmail.com")
	fmt.Print("License Key: ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	key := strings.TrimSpace(input)

	if key == "" {
		fmt.Println(red("License key cannot be empty."))
		return false
	}

	// Validate license key format
	if !validateLicenseKey(key) {
		fmt.Println(red("Invalid license key format. Please check your key and try again."))
		return false
	}

	// Get machine ID
	machineID, err := machineid.ID()
	if err != nil {
		fmt.Println(red("Failed to get machine ID. Please try again."))
		return false
	}

	// Create license
	license := License{
		Key:         key,
		MachineID:   machineID,
		ActivatedAt: time.Now(),
		ExpiryDate:  time.Now().AddDate(0, 5, 0), // 5 months from now
	}

	// Save license
	licenseData, err := json.MarshalIndent(license, "", "  ")
	if err != nil {
		fmt.Println(red("Failed to create license. Please try again."))
		return false
	}

	err = os.WriteFile("license.dat", licenseData, 0644)
	if err != nil {
		fmt.Println(red("Failed to save license. Please try again."))
		return false
	}

	fmt.Println(green("✅ License activated successfully!"))
	fmt.Printf("Valid until: %s\n", license.ExpiryDate.Format("2006-01-02"))
	fmt.Println("You can now use IPO~Master for the next 5 months.")
	return true
}

func showIntroMsg() {
	Log("----------------------------------------------------------------------------")
	Log("IPO~Master BY dallefx - Monitor and automatically apply to open IPOs in MeroShare")
	Log("----------------------------------------------------------------------------")
}

func createNewProfile() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a name for this profile: ")
	input, _ := reader.ReadString('\n')
	name := strings.TrimSpace(input)

	if name == "" {
		fmt.Println(red("Profile name cannot be empty."))
		return
	}

	// Sanitize filename
	sanitizedName := strings.ReplaceAll(name, " ", "_")
	sanitizedName = strings.ReplaceAll(sanitizedName, "/", "_")
	sanitizedName = strings.ReplaceAll(sanitizedName, "\\", "_")

	config := Config{Name: name}

	fmt.Print("Enter DP ID: ")
	input, _ = reader.ReadString('\n')
	config.DPID = strings.TrimSpace(input)

	fmt.Print("Enter BO ID: ")
	input, _ = reader.ReadString('\n')
	config.BOID = strings.TrimSpace(input)

	fmt.Print("Enter Password: ")
	input, _ = reader.ReadString('\n')
	config.Password = strings.TrimSpace(input)

	// Login to get auth token
	clientIdDict := GetClientIds()
	clientIdStr := clientIdDict[config.DPID]
	if clientIdStr == "" {
		fmt.Println(red("Invalid DP ID."))
		return
	}
	clientId := 0
	fmt.Sscan(clientIdStr, &clientId)
	authRequestBody := map[string]interface{}{"clientId": clientId, "username": config.BOID, "password": config.Password}
	req, _ := json.Marshal(authRequestBody)
	resp, err := retryHTTPRequest("POST", "https://webbackend.cdsc.com.np/api/meroShare/auth/", "application/json", bytes.NewBuffer(req), 3)
	if err != nil {
		fmt.Println(red("Authentication failed. Please check your credentials."))
		return
	}
	authToken := resp.Header.Get("Authorization")
	if authToken == "" {
		fmt.Println(red("Failed to get auth token."))
		return
	}

	// Fetch banks
	request, _ := http.NewRequest("GET", "https://webbackend.cdsc.com.np/api/meroShare/bank/", bytes.NewBufferString(""))
	request.Header.Add("Authorization", authToken)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(red("Error getting banks."))
		return
	}
	defer response.Body.Close()

	var bankBriefs []BankBrief
	err = json.NewDecoder(response.Body).Decode(&bankBriefs)
	if err != nil {
		fmt.Println(red("Error parsing bank data."))
		return
	}

	fmt.Println("Please select the default bank:")
	for i, bankBrief := range bankBriefs {
		fmt.Printf("%d. %s (ID: %d)\n", i+1, bankBrief.Name, bankBrief.Id)
	}
	fmt.Print("Enter the number: ")
	bankSelectionStr, _ := reader.ReadString('\n')
	bankSelection := 0
	fmt.Sscan(strings.TrimSpace(bankSelectionStr), &bankSelection)
	if bankSelection < 1 || bankSelection > len(bankBriefs) {
		fmt.Println(red("Invalid bank selection."))
		return
	}
	selectedBank := bankBriefs[bankSelection-1]
	config.DefaultBankID = selectedBank.Id
	fmt.Printf("Selected: %s\n", selectedBank.Name)

	fmt.Print("Enter CRN for this bank: ")
	input, _ = reader.ReadString('\n')
	config.CRN = strings.TrimSpace(input)

	fmt.Print("Enter Default Kittas (press Enter for 10): ")
	input, _ = reader.ReadString('\n')
	kittasStr := strings.TrimSpace(input)
	if kittasStr == "" {
		config.DefaultKittas = 10
	} else {
		if kittas, err := strconv.Atoi(kittasStr); err == nil {
			config.DefaultKittas = kittas
		} else {
			config.DefaultKittas = 10
		}
	}

	fmt.Print("Ask for kittas each time? (y/n, default: n): ")
	input, _ = reader.ReadString('\n')
	askStr := strings.TrimSpace(strings.ToLower(input))
	config.AskForKittas = askStr == "y" || askStr == "yes"

	// Generate encryption key
	key := generateKey(sanitizedName)

	// Encrypt sensitive data
	encryptedPassword, err := encrypt(config.Password, key)
	if err != nil {
		fmt.Println(red("Error encrypting password."))
		return
	}
	config.Password = encryptedPassword

	encryptedPIN, err := encrypt(config.TransactionPIN, key)
	if err != nil {
		fmt.Println(red("Error encrypting PIN."))
		return
	}
	config.TransactionPIN = encryptedPIN

	// Save config
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println(red("Error saving profile."))
		return
	}

	filename := "profile_" + sanitizedName + ".json"
	err = os.WriteFile(filename, configData, 0644)
	if err != nil {
		fmt.Println(red("Error saving profile."))
		return
	}

	// Save encryption key
	keyFilename := "key_" + sanitizedName + ".dat"
	err = os.WriteFile(keyFilename, []byte(key), 0644)
	if err != nil {
		fmt.Println(red("Error saving encryption key."))
		return
	}

	fmt.Println(green("Profile created successfully!"))
}

func loadProfile(filename string) (Config, error) {
	configData, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(configData, &config)
	if err != nil {
		return Config{}, err
	}

	// Load encryption key
	sanitizedName := strings.TrimPrefix(strings.TrimSuffix(filename, ".json"), "profile_")
	keyFilename := "key_" + sanitizedName + ".dat"
	keyData, err := os.ReadFile(keyFilename)
	if err != nil {
		return Config{}, err
	}
	key := string(keyData)

	// Decrypt sensitive data
	decryptedPassword, err := decrypt(config.Password, key)
	if err != nil {
		return Config{}, err
	}
	config.Password = decryptedPassword

	decryptedPIN, err := decrypt(config.TransactionPIN, key)
	if err != nil {
		return Config{}, err
	}
	config.TransactionPIN = decryptedPIN

	return config, nil
}

func getProfileName(filename string) string {
	configData, err := os.ReadFile(filename)
	if err != nil {
		return filename
	}

	var config Config
	err = json.Unmarshal(configData, &config)
	if err != nil {
		return filename
	}

	return config.Name
}

func runSelectedProfilesContinuously() {
	filesA, errA := filepath.Glob("./profile_*.json")
	filesB, errB := filepath.Glob("./config*.json")
	if errA != nil && errB != nil {
		LogLevel("ERROR", "Error retrieving config file list!")
		return
	}
	files := append(filesA, filesB...)
	if len(files) == 0 {
		fmt.Println(yellow("No profiles found. Please create a profile first."))
		return
	}

	fmt.Println(cyan("\nAvailable profiles:"))
	for i, file := range files {
		fmt.Printf("%d. %s\n", i+1, getProfileName(file))
	}

	fmt.Print("Enter profile numbers to run (comma-separated, e.g., 1,3,5): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		fmt.Println(red("No profiles selected."))
		return
	}

	selectedIndices := strings.Split(input, ",")
	var selectedProfiles []string

	for _, idxStr := range selectedIndices {
		idxStr = strings.TrimSpace(idxStr)
		idx, err := strconv.Atoi(idxStr)
		if err != nil || idx < 1 || idx > len(files) {
			fmt.Printf(red("Invalid profile number: %s\n"), idxStr)
			continue
		}
		selectedProfiles = append(selectedProfiles, files[idx-1])
	}

	if len(selectedProfiles) == 0 {
		fmt.Println(red("No valid profiles selected."))
		return
	}

	fmt.Println(green("Starting continuous monitoring for selected profiles..."))
	fmt.Println("Press Ctrl+C to stop.")

	var wg sync.WaitGroup
	stopChan := make(chan bool)

	for _, profileFile := range selectedProfiles {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			monitorIPOs(file, stopChan)
		}(profileFile)
	}

	// Wait for interrupt
	// c := make(chan os.Signal, 1) // Commented out for Windows compatibility

	go func() {
		// <-c
		time.Sleep(24 * time.Hour) // Run for 24 hours or until interrupted
		close(stopChan)
	}()

	wg.Wait()
	fmt.Println(green("Monitoring stopped."))
}

func runAllProfilesOnce() {
	filesA, errA := filepath.Glob("./profile_*.json")
	filesB, errB := filepath.Glob("./config*.json")
	if errA != nil && errB != nil {
		LogLevel("ERROR", "Error retrieving config file list!")
		return
	}
	files := append(filesA, filesB...)
	if len(files) == 0 {
		fmt.Println(yellow("No profiles found. Please create a profile first."))
		return
	}

	fmt.Println(green("Running all profiles once..."))

	for _, file := range files {
		config, err := loadProfile(file)
		if err != nil {
			LogLevel("ERROR", fmt.Sprintf("Failed to load profile %s: %v", file, err))
			continue
		}

		LogLevel("INFO", fmt.Sprintf("Processing profile: %s", config.Name))
		applyToOpenIPOs(config)
	}

	fmt.Println(green("All profiles processed."))
}

func monitorIPOs(profileFile string, stopChan <-chan bool) {
	config, err := loadProfile(profileFile)
	if err != nil {
		LogLevel("ERROR", fmt.Sprintf("Failed to load profile %s: %v", profileFile, err))
		return
	}

	LogLevel("INFO", fmt.Sprintf("Started monitoring for profile: %s", config.Name))

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-stopChan:
			LogLevel("INFO", fmt.Sprintf("Stopping monitoring for profile: %s", config.Name))
			return
		case <-ticker.C:
			LogLevel("INFO", fmt.Sprintf("Checking IPOs for profile: %s", config.Name))
			applyToOpenIPOs(config)
		}
	}
}

func applyToOpenIPOs(config Config) {
	// Get open IPOs
	ipoList, err := getOpenIPOs()
	if err != nil {
		LogLevel("ERROR", fmt.Sprintf("Failed to get IPO list: %v", err))
		return
	}

	if len(ipoList) == 0 {
		LogLevel("INFO", "No open IPOs found.")
		return
	}

	for _, ipo := range ipoList {
		if ipo.AppliedKittas > 0 {
			LogLevel("INFO", fmt.Sprintf("Already applied to %s (%d kittas)", ipo.CompanyName, ipo.AppliedKittas))
			continue
		}

		kittas := config.DefaultKittas
		if config.AskForKittas {
			fmt.Printf("Enter kittas for %s (default %d): ", ipo.CompanyName, kittas)
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			if input != "" {
				if customKittas, err := strconv.Atoi(input); err == nil {
					kittas = customKittas
				}
			}
		}

		err := applyToIPO(config, ipo, kittas)
		if err != nil {
			LogLevel("ERROR", fmt.Sprintf("Failed to apply to %s: %v", ipo.CompanyName, err))
		} else {
			LogLevel("SUCCESS", fmt.Sprintf("Successfully applied to %s (%d kittas)", ipo.CompanyName, kittas))
		}

		// Small delay between applications
		time.Sleep(2 * time.Second)
	}
}

func getOpenIPOs() ([]IPOData, error) {
	url := "https://iporesult.cdscnp.com.np/result/openIpo"

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result struct {
		Data []IPOData `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

func applyToIPO(config Config, ipo IPOData, kittas int) error {
	url := "https://webbackend.cdscnp.com.np/api/ipo/open/apply"

	payload := map[string]interface{}{
		"accountType":    "Demat",
		"kitta":          kittas,
		"bankId":         config.DefaultBankID,
		"demat":          config.DPID,
		"boid":           config.BOID,
		"clientName":     config.Name,
		"appliedKitta":   kittas,
		"crnNumber":      config.CRN,
		"transactionPIN": config.TransactionPIN,
		"companyShareId": ipo.StockSymbol,
		"password":       config.Password,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var response ApplicationResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("application failed with status %d: %s", resp.StatusCode, response.Object.Message)
	}

	return nil
}

func showLicenseInfo() {
	licenseFile := "license.dat"

	if _, err := os.Stat(licenseFile); os.IsNotExist(err) {
		fmt.Println(red("No license found."))
		return
	}

	licenseData, err := os.ReadFile(licenseFile)
	if err != nil {
		fmt.Println(red("Failed to read license file."))
		return
	}

	var license License
	err = json.Unmarshal(licenseData, &license)
	if err != nil {
		fmt.Println(red("Invalid license format."))
		return
	}

	fmt.Println(cyan("\n=== License Information ==="))
	fmt.Printf("License Key: %s\n", license.Key)
	fmt.Printf("Activated On: %s\n", license.ActivatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Expires On: %s\n", license.ExpiryDate.Format("2006-01-02 15:04:05"))

	daysLeft := int(time.Until(license.ExpiryDate).Hours() / 24)
	if daysLeft > 0 {
		fmt.Printf("Days Remaining: %s\n", green(fmt.Sprintf("%d days", daysLeft)))
	} else {
		fmt.Printf("Status: %s\n", red("EXPIRED"))
		fmt.Println(yellow("Please contact the developer to renew your license."))
		fmt.Println(yellow("Email: planearn01@gmail.com"))
	}
}

func renewLicense() {
	fmt.Println(cyan("\n=== License Renewal ==="))
	fmt.Println("Please contact the developer to obtain a renewal license key.")
	fmt.Println("Email: planearn01@gmail.com")
	fmt.Println("Renewal extends your license for another 5 months.")
	fmt.Print("Enter renewal license key: ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	newKey := strings.TrimSpace(input)

	if newKey == "" {
		fmt.Println(red("License key cannot be empty."))
		return
	}

	// Validate the new license key
	if !validateLicenseKey(newKey) {
		fmt.Println(red("Invalid license key format."))
		return
	}

	// Check if license file exists
	licenseFile := "license.dat"
	if _, err := os.Stat(licenseFile); os.IsNotExist(err) {
		fmt.Println(red("No existing license found. Please activate a license first."))
		return
	}

	// Read current license
	licenseData, err := os.ReadFile(licenseFile)
	if err != nil {
		fmt.Println(red("Failed to read current license."))
		return
	}

	var currentLicense License
	err = json.Unmarshal(licenseData, &currentLicense)
	if err != nil {
		fmt.Println(red("Invalid current license format."))
		return
	}

	// Get current machine ID
	currentMachineID, err := machineid.ID()
	if err != nil {
		fmt.Println(red("Failed to get machine ID."))
		return
	}

	// Verify machine ID matches
	if currentLicense.MachineID != currentMachineID {
		fmt.Println(red("License renewal failed: Machine ID mismatch."))
		return
	}

	// Create renewed license
	renewedLicense := License{
		Key:         newKey,
		MachineID:   currentMachineID,
		ActivatedAt: time.Now(),
		ExpiryDate:  time.Now().AddDate(0, 5, 0), // 5 months from now
	}

	// Save renewed license
	licenseData, err = json.MarshalIndent(renewedLicense, "", "  ")
	if err != nil {
		fmt.Println(red("Failed to create renewed license."))
		return
	}

	err = os.WriteFile(licenseFile, licenseData, 0644)
	if err != nil {
		fmt.Println(red("Failed to save renewed license."))
		return
	}

	fmt.Println(green("License renewed successfully!"))
	fmt.Printf("New expiry date: %s\n", renewedLicense.ExpiryDate.Format("2006-01-02"))
}

func manageProfiles() {
	filesA, errA := filepath.Glob("./profile_*.json")
	filesB, errB := filepath.Glob("./config*.json")
	if errA != nil && errB != nil {
		LogLevel("ERROR", "Error retrieving config file list!")
		return
	}
	files := append(filesA, filesB...)
	if len(files) == 0 {
		fmt.Println(yellow("No profiles found. Please create a profile first."))
		return
	}

	fmt.Println(cyan("\n=== Profile Management ==="))
	for i, file := range files {
		fmt.Printf("%d. %s\n", i+1, getProfileName(file))
	}
	fmt.Println("e. Edit a profile")
	fmt.Println("d. Delete a profile")
	fmt.Println("r. Renew license")
	fmt.Println("b. Back to main menu")
	fmt.Print("Choose an option: ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	choice := strings.TrimSpace(input)

	if choice == "b" {
		return
	}

	if choice == "e" {
		editProfile(files)
	} else if choice == "d" {
		deleteProfile(files)
	} else if choice == "r" {
		renewLicense()
	} else {
		fmt.Println(red("Invalid choice."))
	}
}

func editProfile(files []string) {
	fmt.Print("Enter the number of the profile to edit: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	numStr := strings.TrimSpace(input)
	num, err := strconv.Atoi(numStr)
	if err != nil || num < 1 || num > len(files) {
		fmt.Println(red("Invalid profile number."))
		return
	}

	file := files[num-1]
	fmt.Printf("Editing profile: %s\n", getProfileName(file))

	// Load current config
	configContents, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(red("Error reading profile."))
		return
	}
	var config Config
	err = json.Unmarshal(configContents, &config)
	if err != nil {
		fmt.Println(red("Error parsing profile."))
		return
	}

	// Edit name
	fmt.Printf("Current name: %s\n", config.Name)
	fmt.Print("New name (leave empty to keep current): ")
	input, _ = reader.ReadString('\n')
	newName := strings.TrimSpace(input)
	if newName != "" {
		config.Name = newName
	}

	// Edit other fields
	fmt.Printf("Current DP ID: %s\n", config.DPID)
	fmt.Print("New DP ID (leave empty to keep current): ")
	input, _ = reader.ReadString('\n')
	if newDPID := strings.TrimSpace(input); newDPID != "" {
		config.DPID = newDPID
	}

	fmt.Printf("Current BO ID: %s\n", config.BOID)
	fmt.Print("New BO ID (leave empty to keep current): ")
	input, _ = reader.ReadString('\n')
	if newBOID := strings.TrimSpace(input); newBOID != "" {
		config.BOID = newBOID
	}

	fmt.Print("Update password? (y/n): ")
	input, _ = reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(input)) == "y" {
		fmt.Print("New Password: ")
		input, _ = reader.ReadString('\n')
		config.Password = strings.TrimSpace(input)
	}

	fmt.Printf("Current CRN: %s\n", config.CRN)
	fmt.Print("New CRN (leave empty to keep current): ")
	input, _ = reader.ReadString('\n')
	if newCRN := strings.TrimSpace(input); newCRN != "" {
		config.CRN = newCRN
	}

	fmt.Print("Update Transaction PIN? (y/n): ")
	input, _ = reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(input)) == "y" {
		fmt.Print("New Transaction PIN: ")
		input, _ = reader.ReadString('\n')
		config.TransactionPIN = strings.TrimSpace(input)
	}

	fmt.Printf("Current Default Bank ID: %d\n", config.DefaultBankID)
	fmt.Print("New Default Bank ID (leave empty to keep current): ")
	input, _ = reader.ReadString('\n')
	if newBankIDStr := strings.TrimSpace(input); newBankIDStr != "" {
		if newBankID, err := strconv.Atoi(newBankIDStr); err == nil {
			config.DefaultBankID = newBankID
		}
	}

	fmt.Printf("Current Default Kittas: %d\n", config.DefaultKittas)
	fmt.Print("New Default Kittas (leave empty to keep current): ")
	input, _ = reader.ReadString('\n')
	if newKittasStr := strings.TrimSpace(input); newKittasStr != "" {
		if newKittas, err := strconv.Atoi(newKittasStr); err == nil {
			config.DefaultKittas = newKittas
		}
	}

	fmt.Printf("Current Ask for Kittas: %t\n", config.AskForKittas)
	fmt.Print("Ask for kittas each time? (y/n, leave empty to keep current): ")
	input, _ = reader.ReadString('\n')
	if askStr := strings.TrimSpace(strings.ToLower(input)); askStr != "" {
		config.AskForKittas = askStr == "y" || askStr == "yes"
	}

	// Re-encrypt sensitive data if changed
	sanitizedName := strings.ReplaceAll(config.Name, " ", "_")
	sanitizedName = strings.ReplaceAll(sanitizedName, "/", "_")
	sanitizedName = strings.ReplaceAll(sanitizedName, "\\", "_")
	key := generateKey(sanitizedName)

	encryptedPassword, err := encrypt(config.Password, key)
	if err != nil {
		fmt.Println(red("Error encrypting password."))
		return
	}
	config.Password = encryptedPassword

	encryptedPIN, err := encrypt(config.TransactionPIN, key)
	if err != nil {
		fmt.Println(red("Error encrypting PIN."))
		return
	}
	config.TransactionPIN = encryptedPIN

	// Save updated config
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println(red("Error saving updated profile."))
		return
	}

	err = os.WriteFile(file, configData, 0644)
	if err != nil {
		fmt.Println(red("Error saving updated profile."))
		return
	}

	// Update encryption key file
	keyFilename := "key_" + sanitizedName + ".dat"
	err = os.WriteFile(keyFilename, []byte(key), 0644)
	if err != nil {
		fmt.Println(red("Error saving encryption key."))
		return
	}

	fmt.Println(green("Profile updated successfully!"))
}

func deleteProfile(files []string) {
	fmt.Print("Enter the number of the profile to delete: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	numStr := strings.TrimSpace(input)
	num, err := strconv.Atoi(numStr)
	if err != nil || num < 1 || num > len(files) {
		fmt.Println(red("Invalid profile number."))
		return
	}

	file := files[num-1]
	profileName := getProfileName(file)

	fmt.Printf("Are you sure you want to delete profile '%s'? (y/n): ", profileName)
	input, _ = reader.ReadString('\n')
	confirm := strings.TrimSpace(strings.ToLower(input))

	if confirm != "y" && confirm != "yes" {
		fmt.Println("Deletion cancelled.")
		return
	}

	// Delete profile file
	err = os.Remove(file)
	if err != nil {
		fmt.Println(red("Error deleting profile file."))
		return
	}

	// Delete encryption key file
	sanitizedName := strings.TrimPrefix(strings.TrimSuffix(file, ".json"), "profile_")
	keyFile := "key_" + sanitizedName + ".dat"
	os.Remove(keyFile) // Ignore error if key file doesn't exist

	fmt.Println(green("Profile deleted successfully!"))
}

func GetClientIds() map[string]string {
	clientIdDict := map[string]string{
		"13200": "128",
		"12300": "129",
		"17200": "130",
		"11900": "131",
		"15600": "132",
		"17500": "201",
		"14700": "133",
		"11100": "134",
		"15000": "135",
		"16000": "136",
		"11700": "137",
		"10100": "138",
		"13300": "139",
		"13400": "140",
		"12000": "141",
		"14500": "142",
		"11300": "143",
		"14900": "144",
		"10800": "145",
		"17600": "153",
		"12200": "151",
		"11200": "146",
		"16200": "147",
		"18000": "681",
		"17700": "148",
		"17400": "149",
		"13100": "150",
		"17900": "402",
		"18200": "1182",
		"14300": "154",
		"15200": "156",
		"10700": "157",
		"13800": "158",
		"16100": "159",
		"14100": "155",
		"16700": "160",
		"13600": "161",
		"17300": "162",
		"12500": "199",
		"15900": "163",
		"16800": "198",
		"15100": "166",
		"10400": "164",
		"16400": "165",
		"15700": "167",
		"16300": "168",
		"15500": "169",
		"15300": "170",
		"11500": "171",
		"10200": "172",
		"10600": "173",
		"13700": "174",
		"11000": "175",
		"11800": "176",
		"17000": "177",
		"13900": "178",
		"12600": "179",
		"14800": "180",
		"16900": "181",
		"15400": "152",
		"12800": "182",
		"18600": "1270",
		"16600": "183",
		"16500": "184",
		"18100": "1080",
		"14400": "185",
		"15800": "186",
		"11600": "187",
		"12700": "188",
		"18400": "1189",
		"18500": "1196",
		"12900": "189",
		"10900": "190",
		"14600": "191",
		"13000": "192",
		"14000": "193",
		"14200": "194",
		"17800": "370",
		"12400": "195",
		"18300": "1186",
		"11400": "196",
		"17100": "197",
		"13500": "200",
	}
	return clientIdDict
}

func retryHTTPRequest(method, url, contentType string, body io.Reader, maxRetries int) (*http.Response, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	for i := 0; i < maxRetries; i++ {
		var resp *http.Response
		var err error

		if method == "POST" {
			resp, err = client.Post(url, contentType, body)
		} else {
			req, _ := http.NewRequest(method, url, body)
			if contentType != "" {
				req.Header.Set("Content-Type", contentType)
			}
			resp, err = client.Do(req)
		}

		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		if i < maxRetries-1 {
			Log("Request failed, retrying...")
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}
	return nil, fmt.Errorf("request failed after %d retries", maxRetries)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	showIntroMsg()

	// Check license on startup
	if !checkLicense() {
		fmt.Println(yellow("No valid license found."))
		for !activateLicense() {
			fmt.Println(yellow("License activation required to continue."))
		}
	}

mainMenu:
	for {
		fmt.Println(cyan("\n=== IPO~Master BY dallefx Menu ==="))
		fmt.Println("1. Create a new profile")
		fmt.Println("2. Run selected profiles continuously")
		fmt.Println("3. Run all profiles once")
		fmt.Println("4. Manage profiles (edit/delete)")
		fmt.Println("5. License information")
		fmt.Println("6. Exit")
		fmt.Print("Choose an option (1-6): ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			createNewProfile()
		case "2":
			runSelectedProfilesContinuously()
		case "3":
			runAllProfilesOnce()
		case "4":
			manageProfiles()
		case "5":
			showLicenseInfo()
		case "6":
			fmt.Println("Goodbye!")
			break mainMenu
		default:
			fmt.Println(red("Invalid choice."))
		}
	}
}
