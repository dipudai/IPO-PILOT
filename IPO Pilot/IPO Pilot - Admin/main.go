package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

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

type License struct {
	Key         string
	MachineID   string
	ExpiryDate  time.Time
	ActivatedAt time.Time
}

type AvailableIssueObject struct {
	CompanyShareId int
	SubGroup       string
	Scrip          string
	CompanyName    string
	ShareTypeName  string
	ShareGroupName string
	StatusName     string
	IssueOpenDate  string
	IssueCloseDate string
}

type BankBrief struct {
	Code string
	Id   int
	Name string
}

type ScripToApply struct {
	CompanyShareId string
	KittasToApply  string
	BankIdToApply  int
	CompanyName    string
}

type BankDetail struct {
	AccountBranchId int
	AccountNumber   string
	AccountTypeId   int
	AccountTypeName string
	BranchName      string
	Id              int
}

type OwnDetail struct {
	Address                string
	Boid                   string
	ClientCode             string
	Contact                string
	CreatedApproveDate     string
	CreatedApproveDateStr  string
	CustomerTypeCode       string
	Demat                  string
	DematExpiryDate        string
	Email                  string
	ExpiredDate            string
	ExpiredDateStr         string
	Gender                 string
	Id                     int
	ImagePath              string
	MeroShareEmail         string
	Name                   string
	PanNumber              string
	PasswordChangeDate     string
	PasswordChangedDateStr string
	PasswordExpiryDate     string
	PasswordExpiryDateStr  string
	ProfileName            string
	RenderDashboard        bool
	RenewedDate            string
	RenewedDateStr         string
	Username               string
}

type ApplyScripPayloadJSON struct {
	AccountBranchId int    `json:"accountBranchId"`
	AccountNumber   string `json:"accountNumber"`
	AppliedKitta    string `json:"appliedKitta"`
	AccountTypeId   int    `json:"accountTypeId"`
	BankId          int    `json:"bankId"`
	Boid            string `json:"boid"`
	CompanyShareId  string `json:"companyShareId"`
	CrnNumber       string `json:"crnNumber"`
	CustomerId      int    `json:"customerId"`
	Demat           string `json:"demat"`
	TransactionPIN  string `json:"transactionPIN"`
}

func checkLicense() bool {
	licenseFile := "license.dat"

	// Check if activation file exists (one-time activation)
	if _, err := os.Stat(licenseFile); os.IsNotExist(err) {
		return false // Not activated yet
	}

	// Read activation status
	activationData, err := os.ReadFile(licenseFile)
	if err != nil {
		LogLevel("ERROR", "Failed to read activation file")
		return false
	}

	var license License
	err = json.Unmarshal(activationData, &license)
	if err != nil {
		LogLevel("ERROR", "Invalid activation data")
		return false
	}

	// Check machine ID
	machineID, err := machineid.ID()
	if err != nil {
		LogLevel("ERROR", "Failed to get machine ID")
		return false
	}

	if license.MachineID != machineID {
		LogLevel("ERROR", "License is not valid for this machine")
		return false
	}

	LogLevel("SUCCESS", fmt.Sprintf("App activated on %s", license.ActivatedAt.Format("2006-01-02 15:04")))
	return true
}

func activateLicense() bool {
	fmt.Println(cyan("\n=== App Activation ==="))
	fmt.Print("Enter activation password: ")

	reader := bufio.NewReader(os.Stdin)
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	// Check password
	correctPassword := "Earn&Grow0309"
	if password != correctPassword {
		fmt.Println(red("Invalid password"))
		return false
	}

	// Generate activation data
	machineID, err := machineid.ID()
	if err != nil {
		LogLevel("ERROR", "Failed to get machine ID")
		return false
	}

	// Create activation record
	license := License{
		Key:         "activated",
		MachineID:   machineID,
		ActivatedAt: time.Now(),
	}

	// Save activation status
	activationData, _ := json.Marshal(license)
	err = os.WriteFile("license.dat", activationData, 0644)
	if err != nil {
		LogLevel("ERROR", "Failed to save activation")
		return false
	}

	fmt.Println(green("App activated successfully!"))
	fmt.Println(cyan("No license renewal needed. Enjoy unlimited access!"))
	fmt.Println(cyan("For password changes, contact developer: dallefx"))
	fmt.Println(cyan("Email: sahmahesh2077@gmail.com"))
	return true
}

func validateLicenseKey(key string) bool {
	// Deprecated - password-based activation now in use
	return true
}

func generateLicenseKey() string {
	// Deprecated - password-based activation now in use
	return "password-based"
}

func showIntroMsg() {
	Log("----------------------------------------------------------------------------")
	Log("IPO Pilot - Automated IPO Application System")
	Log("Professional IPO monitoring and automatic application for MeroShare")
	Log("----------------------------------------------------------------------------")
}

func createNewProfile() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a name for this profile: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "" {
		fmt.Println(red("Profile name cannot be empty."))
		return
	}
	if len(name) > 50 {
		fmt.Println(red("Profile name too long (max 50 characters)."))
		return
	}

	// create a filename from profile name
	san := sanitizeFilename(name)
	baseName := fmt.Sprintf("profile_%s.json", san)
	// avoid overwriting existing file
	if _, err := os.Stat(baseName); err == nil {
		baseName = fmt.Sprintf("profile_%s_%s.json", san, GetTimestamp())
	}
	configFileName := baseName
	LogLevel("INFO", fmt.Sprintf("Creating new profile '%s' in %s", name, configFileName))

	// Create a basic config with name
	var config Config
	config.Name = name

	// Save initial config with name
	serializedData, _ := json.MarshalIndent(config, "", " ")
	_ = os.WriteFile(configFileName, serializedData, 0666)

	// Now run DoWork to fill in the rest
	DoWork(configFileName)
}

func getProfileName(file string) string {
	configContents, err := os.ReadFile(file)
	var config Config
	if err == nil {
		json.Unmarshal(configContents, &config)
	}
	if config.Name != "" {
		return config.Name
	}
	base := filepath.Base(file)
	// Fallback to BOID or DPID if name not set
	if config.BOID != "" {
		return fmt.Sprintf("%s (BOID: %s)", base, config.BOID)
	}
	if config.DPID != "" {
		return fmt.Sprintf("%s (DPID: %s)", base, config.DPID)
	}
	return base
}

// sanitizeFilename makes a safe filename from a profile name
func sanitizeFilename(name string) string {
	name = strings.TrimSpace(name)
	var b strings.Builder
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
			b.WriteRune(r)
		} else if unicode.IsSpace(r) || r == '+' || r == '.' {
			b.WriteRune('_')
		}
		// else drop other characters
	}
	s := b.String()
	if s == "" {
		s = GetTimestamp()
	}
	return s
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
			LogLevel("WARN", fmt.Sprintf("Request failed, retrying in %d seconds... (%d/%d)", i+1, i+1, maxRetries))
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}
	return nil, fmt.Errorf("request failed after %d retries", maxRetries)
}

func runAllProfilesOnce() {
	filesA, errA := filepath.Glob("./profile_*.json")
	filesB, errB := filepath.Glob("./config*.json")
	if errA != nil && errB != nil {
		Panic("Error retrieving config file list!")
	}
	files := append(filesA, filesB...)
	if len(files) == 0 {
		fmt.Println("No profiles found. Please create a profile first.")
		return
	}
	for _, file := range files {
		Log(fmt.Sprintf("Working on profile: %s", getProfileName(file)))
		DoWork(file)
	}
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
		fmt.Println(yellow("Please renew your license to continue using the software."))
	}
}

func renewLicense() {
	fmt.Println(cyan("\n=== About App Activation ==="))
	fmt.Println(green("✓ App activated (one-time)"))
	fmt.Println(yellow("No renewal needed!"))
	fmt.Println(cyan("For password or account changes, contact: dallefx"))
	fmt.Println(cyan("Email: sahmahesh2077@gmail.com"))
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
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name != "" {
		config.Name = name
	}

	// Save updated config
	serializedData, _ := json.MarshalIndent(config, "", " ")
	err = os.WriteFile(file, serializedData, 0666)
	if err != nil {
		fmt.Println(red("Error saving profile."))
	} else {
		fmt.Println(green("Profile updated successfully."))
	}
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

	fmt.Printf("Are you sure you want to delete profile '%s'? (y/N): ", profileName)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm == "y" || confirm == "yes" {
		// Delete config file
		err := os.Remove(file)
		if err != nil {
			fmt.Println(red("Error deleting profile file."))
			return
		}

		// Delete key file
		keyFile := "key_" + file + ".dat"
		os.Remove(keyFile) // Ignore error if key file doesn't exist

		fmt.Println(green("Profile deleted successfully."))
	} else {
		fmt.Println("Deletion cancelled.")
	}
}

func runSelectedProfilesContinuously() {
	filesA, errA := filepath.Glob("./profile_*.json")
	filesB, errB := filepath.Glob("./config*.json")
	if errA != nil && errB != nil {
		Panic("Error retrieving config file list!")
	}
	files := append(filesA, filesB...)
	if len(files) == 0 {
		fmt.Println("No profiles found. Please create a profile first.")
		return
	}

	fmt.Println("Available profiles:")
	for i, file := range files {
		fmt.Printf("%d. %s\n", i+1, getProfileName(file))
	}
	fmt.Print("Enter the numbers of profiles to run continuously (comma-separated, e.g. 1,2,3): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	selections := strings.TrimSpace(input)

	selectedFiles := []string{}
	for _, sel := range strings.Split(selections, ",") {
		index := strings.TrimSpace(sel)
		if num, err := strconv.Atoi(index); err == nil && num >= 1 && num <= len(files) {
			selectedFiles = append(selectedFiles, files[num-1])
		}
	}

	if len(selectedFiles) == 0 {
		fmt.Println("No valid profiles selected.")
		return
	}

	fmt.Println("Starting continuous monitoring for selected profiles. Press Ctrl+C to stop.")
	var wg sync.WaitGroup
	for _, file := range selectedFiles {
		wg.Add(1)
		go func(configFile string) {
			defer wg.Done()
			runProfileContinuously(configFile)
		}(file)
	}
	wg.Wait()
}

func setupProfile(configFileName string) (*Config, string, OwnDetail, BankDetail) {
	clientIdDict := GetClientIds() //load dpid:clientid dictionary

	key, keyerror := GetKey(configFileName)
	if keyerror != nil {
		_ = os.Remove(configFileName) //since keyerror occurred, we don't have the key for decrypting config, so delete the old config

		makeKeyErr := MakeKey(configFileName)
		if makeKeyErr != nil {
			Panic("Could not create key file!") //we won't proceed without creating a key file
		}
		var getKeyErr error
		key, getKeyErr = GetKey(configFileName)
		if getKeyErr != nil {
			Panic("Could not read key file!") //we won't proceed without a key
		}
	}

	//at this point, we will have either the old key or a newly created one
	var config Config

	configContents, err := os.ReadFile(configFileName)
	if err == nil {
		err = json.Unmarshal(configContents, &config)                             //read the config file into config variable
		config.Password, _ = DecryptAES([]byte(key), config.Password)             //decrypt the saved password
		config.TransactionPIN, _ = DecryptAES([]byte(key), config.TransactionPIN) //decrypt the saved transaction PIN
		config.CRN, _ = DecryptAES([]byte(key), config.CRN)                       //decrypt the saved CRN
	}

	if config.BOID == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter your BOID (last eight digits): ")
		boid, _ := reader.ReadString('\n')
		config.BOID = strings.TrimSpace(boid)
	}

	if config.Password == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter your password: ")
		pwd, _ := reader.ReadString('\n')
		config.Password = strings.TrimSpace(pwd)
	}

	if config.TransactionPIN == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter your transaction PIN: ")
		tpin, _ := reader.ReadString('\n')
		config.TransactionPIN = strings.TrimSpace(tpin)
	}

	if config.DPID == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter your Depository Participant ID: ")
		dpid, _ := reader.ReadString('\n')
		config.DPID = strings.TrimSpace(dpid)
	}

	//ask for default number of kittas to apply
	if config.DefaultKittas == 0 {
		fmt.Println("Please enter the default no. of kittas to apply: ")
		reader := bufio.NewReader(os.Stdin)
		kittasStr, _ := reader.ReadString('\n')
		kittas := 0
		fmt.Sscan(strings.TrimSpace(kittasStr), &kittas)
		config.DefaultKittas = kittas
	}

	//login to get the auth token(JWT)
	clientIdStr := clientIdDict[config.DPID]
	clientId := 0
	fmt.Sscan(clientIdStr, &clientId)
	authRequestBody := map[string]interface{}{"clientId": clientId, "username": config.BOID, "password": config.Password}
	req, _ := json.Marshal(authRequestBody)
	resp, err := retryHTTPRequest("POST", "https://webbackend.cdsc.com.np/api/meroShare/auth/", "application/json", bytes.NewBuffer(req), 3)
	if err != nil {
		LogLevel("ERROR", "Authentication failed after retries")
		return nil, "", OwnDetail{}, BankDetail{}
	}
	authToken := resp.Header.Get("Authorization")

	//retrieve demat
	request, _ := http.NewRequest("GET", "https://webbackend.cdsc.com.np/api/meroShare/ownDetail/", nil)
	request.Header.Add("Authorization", authToken)
	request.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil || response.StatusCode != 200 {
		LogLevel("ERROR", "Failed to get own details")
		return nil, "", OwnDetail{}, BankDetail{}
	}
	defer response.Body.Close()
	var ownDetail OwnDetail
	err = json.NewDecoder(response.Body).Decode(&ownDetail)
	if err != nil {
		Panic("Error parsing own details JSON!")
	}
	Log(fmt.Sprint("Obtained auth token for ", ownDetail.Name))

	//ask for default bank and its corresponding CRN to apply from
	if config.DefaultBankID == 0 {
		//load bank brief
		request, _ := http.NewRequest("GET", "https://webbackend.cdsc.com.np/api/meroShare/bank/", bytes.NewBufferString(""))
		request.Header.Add("Authorization", authToken)
		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			Panic("Error getting banks!")
		}
		defer response.Body.Close()

		var bankBriefs []BankBrief
		err = json.NewDecoder(response.Body).Decode(&bankBriefs)
		if err != nil {
			Panic("Error parsing bank briefs JSON!")
		}
		fmt.Println("Please select the bank you'd like to use:")
		for i, bankBrief := range bankBriefs {
			fmt.Printf("%d. %s (ID: %d)\n", i+1, bankBrief.Name, bankBrief.Id)
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter the number: ")
		bankSelectionStr, _ := reader.ReadString('\n')
		bankSelection := 0
		fmt.Sscan(strings.TrimSpace(bankSelectionStr), &bankSelection)
		if bankSelection < 1 || bankSelection > len(bankBriefs) {
			Panic("Invalid bank selection!")
		}
		selectedBank := bankBriefs[bankSelection-1]
		config.DefaultBankID = selectedBank.Id
		fmt.Printf("Selected: %s\n", selectedBank.Name)

		//ask for the bank's CRN
		fmt.Println("Please enter the CRN of this bank: ")
		bankCRN, _ := reader.ReadString('\n')
		config.CRN = strings.TrimSpace(bankCRN)
	}

	//serialize the populated config variable to config file after encrypting them
	encryptedConfig := config //make a copy of config
	encryptedConfig.Password, _ = EncryptAES([]byte(key), config.Password)
	encryptedConfig.TransactionPIN, _ = EncryptAES([]byte(key), config.TransactionPIN)
	encryptedConfig.CRN, _ = EncryptAES([]byte(key), config.CRN)
	serializedData, _ := json.MarshalIndent(encryptedConfig, "", " ")
	_ = os.WriteFile(configFileName, serializedData, 0666)

	//retrieve accountNumber and customerId fields (required to apply) from an API call to get the default bank's details
	request, _ = http.NewRequest("GET", "https://webbackend.cdsc.com.np/api/meroShare/bank/"+strconv.Itoa(config.DefaultBankID), nil)
	request.Header.Add("Authorization", authToken)
	request.Header.Add("Content-Type", "application/json")
	client = &http.Client{}
	response, err = client.Do(request)
	if err != nil || response.StatusCode != 200 {
		Panic("Error getting bank details!")
	}
	defer response.Body.Close()

	var bankDetails []BankDetail
	err = json.NewDecoder(response.Body).Decode(&bankDetails)
	if err != nil {
		Panic("Error parsing bank details JSON!")
	}

	//take the first element of the JSON array, which is the required BankDetail JSON object
	bankDetail := bankDetails[0]

	return &config, authToken, ownDetail, bankDetail
}

func runProfileContinuously(configFile string) {
	config, authToken, ownDetail, bankDetail := setupProfile(configFile)
	if config == nil {
		return
	}
	Log(fmt.Sprintf("Starting continuous monitoring for profile '%s'", config.Name))
	for {
		applyIPOs(config, authToken, ownDetail, bankDetail)
		Log("Sleeping for 5 minutes before next check...")
		time.Sleep(5 * time.Minute)
	}
}

func main() {

	showIntroMsg()

	// Check license first
	if !checkLicense() {
		fmt.Println(yellow("No valid license found. Please activate your license to continue."))
		for {
			fmt.Println("\n1. Activate License")
			fmt.Println("2. Generate Sample License Key (for testing)")
			fmt.Println("3. Exit")
			fmt.Print("Choose an option: ")

			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			choice := strings.TrimSpace(input)

			switch choice {
			case "1":
				if activateLicense() {
					goto mainMenu
				}
			case "2":
				key := generateLicenseKey()
				fmt.Printf("Sample License Key: %s\n", green(key))
				fmt.Println(yellow("Note: This is for testing only. Use a proper license key for production."))
			case "3":
				fmt.Println("Exiting...")
				return
			default:
				fmt.Println(red("Invalid choice."))
			}
		}
	}

mainMenu:
	for {
		fmt.Println(cyan("\n=== IPO Pilot - Main Menu ==="))
		fmt.Println("1. Create a new profile")
		fmt.Println("2. Run selected profiles continuously")
		fmt.Println("3. Run all profiles once")
		fmt.Println("4. Manage profiles (edit/delete)")
		fmt.Println("5. Activation information")
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
			fmt.Println(green("Exiting..."))
			return
		default:
			fmt.Println(red("Invalid choice. Please try again."))
		}
	}
}

func DoWork(configFileName string) {
	config, authToken, ownDetail, bankDetail := setupProfile(configFileName)
	if config == nil {
		return
	}
	applyIPOs(config, authToken, ownDetail, bankDetail)
}

func applyIPOs(config *Config, authToken string, ownDetail OwnDetail, bankDetail BankDetail) {
	//retrieve available issues and create a slice of scrips to apply to
	reqBodyForAvailableIssues := []byte(`{
		"filterFieldParams": [
		  {
			"key": "companyIssue.companyISIN.script",
			"alias": "Scrip"
		  },
		  {
			"key": "companyIssue.companyISIN.company.name",
			"alias": "Company Name"
		  },
		  {
			"key": "companyIssue.assignedToClient.name",
			"value": "",
			"alias": "Issue Manager"
		  }
		],
		"page": 1,
		"size": 10,
		"searchRoleViewConstants": "VIEW_APPLICABLE_SHARE",
		"filterDateParams": [
		  {
			"key": "minIssueOpenDate",
			"condition": "",
			"alias": "",
			"value": ""
		  },
		  {
			"key": "maxIssueCloseDate",
			"condition": "",
			"alias": "",
			"value": ""
		  }
		]
	  }`)
	request, _ := http.NewRequest("POST", "https://webbackend.cdsc.com.np/api/meroShare/companyShare/applicableIssue/", bytes.NewBuffer(reqBodyForAvailableIssues))
	request.Header.Add("Authorization", authToken)
	request.Header.Add("Content-Type", "application/json")
	response, err := retryHTTPRequest("POST", "https://webbackend.cdsc.com.np/api/meroShare/companyShare/applicableIssue/", "application/json", bytes.NewBuffer(reqBodyForAvailableIssues), 3)
	if err != nil {
		LogLevel("ERROR", "Failed to get applicable issues after retries")
		return
	}
	defer response.Body.Close()

	var responseJson map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&responseJson)
	if err != nil {
		Log("Error parsing applicable issues JSON!")
		return
	}

	//damned conversion from []interface{} to []AvailableIssueObject{}
	availableScrips := []AvailableIssueObject{}
	availables := responseJson["object"].([]interface{})
	for _, v := range availables {
		tmp := v.(map[string]interface{})
		if _, isKeyPresent := tmp["action"]; isKeyPresent { //there's an "action" key if this scrip was already applied to
			Log(fmt.Sprint("Skipping recently applied company ", tmp["companyName"].(string)))
			continue
		}
		jsonString, _ := json.Marshal(tmp)
		var availableIssue AvailableIssueObject
		err := json.Unmarshal(jsonString, &availableIssue)
		if err != nil {
			Log("Error parsing available scrips!")
			continue
		}
		availableScrips = append(availableScrips, availableIssue)
	}

	scripsToApply := []ScripToApply{}
	for _, scrip := range availableScrips {
		if !config.AskForKittas {
			if scrip.ShareGroupName == "Ordinary Shares" {
				var scripToApply ScripToApply
				scripToApply.BankIdToApply = config.DefaultBankID
				scripToApply.KittasToApply = strconv.Itoa(config.DefaultKittas)
				scripToApply.CompanyShareId = strconv.Itoa(scrip.CompanyShareId)
				scripToApply.CompanyName = scrip.CompanyName
				scripsToApply = append(scripsToApply, scripToApply)
			}
		} else {
			Log(fmt.Sprint(scrip.CompanyName, "-", scrip.ShareGroupName, "(", scrip.IssueOpenDate, "-", scrip.IssueCloseDate, ")"))
			fmt.Println("Enter the no. of kittas to apply (0 for none): ")
			reader := bufio.NewReader(os.Stdin)
			kittasStr, _ := reader.ReadString('\n')
			kittas := 0
			fmt.Sscan(strings.TrimSpace(kittasStr), &kittas)
			if kittas > 0 {
				var scripToApply ScripToApply
				scripToApply.BankIdToApply = config.DefaultBankID
				scripToApply.KittasToApply = strconv.Itoa(kittas)
				scripToApply.CompanyShareId = strconv.Itoa(scrip.CompanyShareId)
				scripToApply.CompanyName = scrip.CompanyName
				scripsToApply = append(scripsToApply, scripToApply)
			}

		}
	}

	if len(scripsToApply) == 0 {
		Log("No new IPOs open right now.")
		return
	}

	//apply to the companies in scripstoApply slice
	applyReqJson := &ApplyScripPayloadJSON{
		AccountBranchId: bankDetail.AccountBranchId,
		AccountNumber:   bankDetail.AccountNumber,
		AccountTypeId:   bankDetail.AccountTypeId,
		Boid:            ownDetail.Boid,
		CrnNumber:       config.CRN,
		Demat:           ownDetail.Demat,
		CustomerId:      bankDetail.Id,
		TransactionPIN:  config.TransactionPIN,
	}
	for _, scrip := range scripsToApply {
		applyReqJson.AppliedKitta = scrip.KittasToApply
		applyReqJson.BankId = scrip.BankIdToApply
		applyReqJson.CompanyShareId = scrip.CompanyShareId

		reqjson, err := json.Marshal(applyReqJson)
		if err != nil {
			LogLevel("ERROR", "Cannot build apply JSON payload!")
			continue
		}
		response, err := retryHTTPRequest("POST", "https://webbackend.cdsc.com.np/api/meroShare/applicantForm/share/apply", "application/json", bytes.NewBuffer(reqjson), 3)
		if err != nil || response.StatusCode != 201 {
			LogLevel("ERROR", fmt.Sprintf("Failed to apply to scrip - %s", scrip.CompanyName))
		} else {
			LogLevel("SUCCESS", fmt.Sprintf("Applied %s kittas to - %s", applyReqJson.AppliedKitta, scrip.CompanyName))
		}
	}

	//closing tag
	Log("\n")
}

func EncryptAES(key []byte, message string) (encoded string, err error) {
	//Create byte array from the input string
	plainText := []byte(message)

	//Create a new AES cipher using the key
	block, err := aes.NewCipher(key)

	//IF NewCipher failed, exit:
	if err != nil {
		return
	}

	//Make the cipher text a byte array of size BlockSize + the length of the message
	cipherText := make([]byte, aes.BlockSize+len(plainText))

	//iv is the ciphertext up to the blocksize (16)
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(crand.Reader, iv); err != nil {
		return
	}

	//Encrypt the data:
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//Return string encoded in base64
	return base64.RawStdEncoding.EncodeToString(cipherText), err
}

func DecryptAES(key []byte, secure string) (decoded string, err error) {
	//Remove base64 encoding:
	cipherText, err := base64.RawStdEncoding.DecodeString(secure)

	//IF DecodeString failed, exit:
	if err != nil {
		return
	}

	//Create a new AES cipher with the key and encrypted message
	block, err := aes.NewCipher(key)

	//IF NewCipher failed, exit:
	if err != nil {
		return
	}

	//IF the length of the cipherText is less than 16 Bytes:
	if len(cipherText) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short!")
		return
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	//Decrypt the message
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), err
}

// reads the key file corresponding to the given profile number
func GetKey(configFileName string) ([]byte, error) {
	keyPath := "key_" + configFileName + ".dat"
	key, err := os.ReadFile(keyPath)
	return key, err
}

// creates a random 32-digit key and writes to a key file corresponding to the given profile number
func MakeKey(configFileName string) error {
	keyPath := "key_" + configFileName + ".dat"
	randStr := randString(32)
	error_ := os.WriteFile(keyPath, []byte(randStr), 0644)
	return error_
}

func randString(n int) string {
	const alphanum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@+=[{}];:,./~!@#$%^&*()_"
	symbols := big.NewInt(int64(len(alphanum)))
	states := big.NewInt(0)
	states.Exp(symbols, big.NewInt(int64(n)), nil)
	r, err := crand.Int(crand.Reader, states)
	if err != nil {
		panic(err)
	}
	var bytes = make([]byte, n)
	r2 := big.NewInt(0)
	symbol := big.NewInt(0)
	for i := range bytes {
		r2.DivMod(r, symbols, symbol)
		r, r2 = r2, r
		bytes[i] = alphanum[symbol.Int64()]
	}
	return string(bytes)
}

func GetClientIds() map[string]string {
	//map of "depository participant id" : "clientId"
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

func GetTimestamp() string {
	now := time.Now()
	timestamp := now.Unix()
	return strconv.Itoa(int(timestamp))
}

func Log(msg string) {
	LogLevel("INFO", msg)
}

func LogLevel(level, msg string) {
	var coloredMsg string
	switch level {
	case "ERROR":
		coloredMsg = red("[ERROR] " + msg)
	case "SUCCESS":
		coloredMsg = green("[SUCCESS] " + msg)
	case "WARN":
		coloredMsg = yellow("[WARN] " + msg)
	case "INFO":
		coloredMsg = blue("[INFO] " + msg)
	default:
		coloredMsg = msg
	}

	fmt.Println(coloredMsg)
	file, err := os.OpenFile("IPO_Pilot-Log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(red("Cannot write to log file!"))
	} else {
		t := time.Now()
		timestamp := t.Format("2006-01-02 03:04:05 PM")
		file.WriteString("[" + timestamp + "] [" + level + "] " + msg + "\n")
	}
}

func Panic(msg string) {
	LogLevel("ERROR", msg)
	panic(msg)
}
