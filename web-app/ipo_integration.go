package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Multi-IPO source integration

// Get open IPOs from all configured sources
func getOpenIPOsFromAllSources() ([]IPOData, error) {
	var sources []IPOSource
	db.Where("is_active = ?", true).Order("priority DESC").Find(&sources)

	allIPOs := make([]IPOData, 0)

	for _, source := range sources {
		ipos, err := fetchIPOsFromSource(&source)
		if err != nil {
			fmt.Printf("Error fetching from %s: %v\n", source.Name, err)
			continue
		}
		allIPOs = append(allIPOs, ipos...)
		
		// Update last checked time
		now := time.Now()
		source.LastChecked = &now
		db.Save(&source)
	}

	// Deduplicate IPOs based on company share ID
	return deduplicateIPOs(allIPOs), nil
}

// Fetch IPOs from a specific source
func fetchIPOsFromSource(source *IPOSource) ([]IPOData, error) {
	switch source.Type {
	case "meroshare":
		return fetchFromMeroShare(source)
	case "iporesult":
		return fetchFromIPOResult(source)
	case "cts":
		return fetchFromCTS(source)
	case "custom":
		return fetchFromCustomAPI(source)
	default:
		return nil, fmt.Errorf("unknown source type: %s", source.Type)
	}
}

// Fetch from MeroShare API
func fetchFromMeroShare(source *IPOSource) ([]IPOData, error) {
	url := source.BaseURL + "/api/meroShare/companyShare/applicableIssue/"
	
	reqBody := []byte(`{
		"filterFieldParams": [],
		"page": 1,
		"size": 10,
		"searchRoleViewConstants": "VIEW_APPLICABLE_SHARE",
		"filterDateParams": []
	}`)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if source.APIKey != "" {
		req.Header.Set("Authorization", source.APIKey)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Object []map[string]interface{} `json:"object"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	ipos := make([]IPOData, 0)
	for _, item := range result.Object {
		ipo := IPOData{
			SourceID:       source.ID,
			CompanyName:    getString(item, "companyName"),
			StockSymbol:    getString(item, "scrip"),
			CompanyShareID: fmt.Sprintf("%v", item["companyShareId"]),
			ShareType:      getString(item, "shareTypeName"),
			ShareGroup:     getString(item, "shareGroupName"),
			IssueOpenDate:  getString(item, "issueOpenDate"),
			IssueCloseDate: getString(item, "issueCloseDate"),
			Status:         "open",
			LastUpdated:    time.Now(),
		}
		ipos = append(ipos, ipo)
	}

	return ipos, nil
}

// Fetch from IPO Result API
func fetchFromIPOResult(source *IPOSource) ([]IPOData, error) {
	url := source.BaseURL + "/result/openIpo"

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
		Data []struct {
			CompanyName     string `json:"companyName"`
			StockSymbol     string `json:"stockSymbol"`
			SectorName      string `json:"sectorName"`
			StockPrice      string `json:"stockPrice"`
			MinUnits        string `json:"minUnits"`
			MaxUnits        string `json:"maxUnits"`
			TotalUnits      string `json:"totalUnits"`
			IssueOpenDate   string `json:"issueOpenDate"`
			IssueCloseDate  string `json:"issueCloseDate"`
			ShareType       string `json:"shareType"`
			ShareGroup      string `json:"shareGroup"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	ipos := make([]IPOData, 0)
	for _, item := range result.Data {
		ipo := IPOData{
			SourceID:       source.ID,
			CompanyName:    item.CompanyName,
			StockSymbol:    item.StockSymbol,
			CompanyShareID: item.StockSymbol,
			SectorName:     item.SectorName,
			StockPrice:     item.StockPrice,
			MinUnits:       item.MinUnits,
			MaxUnits:       item.MaxUnits,
			TotalUnits:     item.TotalUnits,
			IssueOpenDate:  item.IssueOpenDate,
			IssueCloseDate: item.IssueCloseDate,
			ShareType:      item.ShareType,
			ShareGroup:     item.ShareGroup,
			Status:         "open",
			LastUpdated:    time.Now(),
		}
		ipos = append(ipos, ipo)
	}

	return ipos, nil
}

// Fetch from CTS (Capital Market) API
func fetchFromCTS(source *IPOSource) ([]IPOData, error) {
	// Placeholder for CTS integration
	// TODO: Implement CTS API integration
	return []IPOData{}, nil
}

// Fetch from custom API
func fetchFromCustomAPI(source *IPOSource) ([]IPOData, error) {
	// Generic custom API integration
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", source.BaseURL, nil)
	if err != nil {
		return nil, err
	}

	if source.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+source.APIKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ipos []IPOData
	if err := json.NewDecoder(resp.Body).Decode(&ipos); err != nil {
		return nil, err
	}

	for i := range ipos {
		ipos[i].SourceID = source.ID
		ipos[i].LastUpdated = time.Now()
	}

	return ipos, nil
}

// Deduplicate IPOs based on company share ID
func deduplicateIPOs(ipos []IPOData) []IPOData {
	seen := make(map[string]bool)
	result := make([]IPOData, 0)

	for _, ipo := range ipos {
		key := ipo.CompanyShareID
		if !seen[key] {
			seen[key] = true
			result = append(result, ipo)
		}
	}

	return result
}

// Get live IPOs handler
func getLiveIPOsHandler(c *gin.Context) {
	ipos, err := getOpenIPOsFromAllSources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch IPOs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(ipos),
		"ipos":  ipos,
	})
}

// Get upcoming IPOs handler
func getUpcomingIPOsHandler(c *gin.Context) {
	// Filter for IPOs that haven't opened yet
	allIPOs, _ := getOpenIPOsFromAllSources()
	upcomingIPOs := make([]IPOData, 0)

	now := time.Now()
	for _, ipo := range allIPOs {
		openDate, err := time.Parse("2006-01-02", ipo.IssueOpenDate)
		if err == nil && openDate.After(now) {
			upcomingIPOs = append(upcomingIPOs, ipo)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(upcomingIPOs),
		"ipos":  upcomingIPOs,
	})
}

// Start monitoring handler
func startMonitoringHandler(c *gin.Context) {
	userID := c.GetUint("userID")

	var input struct {
		ProfileID uint `json:"profile_id" binding:"required"`
		Interval  int  `json:"interval"` // in seconds
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if input.Interval == 0 {
		input.Interval = 300 // Default 5 minutes
	}

	session := MonitoringSession{
		UserID:    userID,
		ProfileID: input.ProfileID,
		IsActive:  true,
		StartedAt: time.Now(),
		Interval:  input.Interval,
	}

	if err := db.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start monitoring"})
		return
	}

	// Start background monitoring goroutine
	go monitorIPOsForSession(&session)

	c.JSON(http.StatusOK, gin.H{
		"message": "Monitoring started",
		"session": session,
	})
}

// Stop monitoring handler
func stopMonitoringHandler(c *gin.Context) {
	userID := c.GetUint("userID")

	var input struct {
		SessionID uint `json:"session_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	now := time.Now()
	db.Model(&MonitoringSession{}).
		Where("id = ? AND user_id = ?", input.SessionID, userID).
		Updates(map[string]interface{}{
			"is_active":  false,
			"stopped_at": now,
		})

	c.JSON(http.StatusOK, gin.H{"message": "Monitoring stopped"})
}

// Monitor status handler
func monitorStatusHandler(c *gin.Context) {
	userID := c.GetUint("userID")

	var sessions []MonitoringSession
	db.Preload("Profile").Where("user_id = ? AND is_active = ?", userID, true).Find(&sessions)

	c.JSON(http.StatusOK, gin.H{
		"count":    len(sessions),
		"sessions": sessions,
	})
}

// Background monitoring function
func monitorIPOsForSession(session *MonitoringSession) {
	ticker := time.NewTicker(time.Duration(session.Interval) * time.Second)
	defer ticker.Stop()

	for {
		// Check if session is still active
		var activeSession MonitoringSession
		if err := db.First(&activeSession, session.ID).Error; err != nil || !activeSession.IsActive {
			return
		}

		// Get open IPOs
		ipos, err := getOpenIPOsFromAllSources()
		if err != nil {
			fmt.Printf("Error fetching IPOs for session %d: %v\n", session.ID, err)
			<-ticker.C
			continue
		}

		// Auto-apply to new IPOs
		var profile Profile
		db.First(&profile, session.ProfileID)

		for _, ipo := range ipos {
			// Check if already applied
			var existingApp IPOApplication
			result := db.Where("user_id = ? AND profile_id = ? AND company_share_id = ?",
				session.UserID, session.ProfileID, ipo.CompanyShareID).First(&existingApp)

			if result.Error != nil { // Not applied yet
				// Create application
				app := IPOApplication{
					UserID:         session.UserID,
					ProfileID:      session.ProfileID,
					IPOSourceID:    ipo.SourceID,
					CompanyName:    ipo.CompanyName,
					CompanyShareID: ipo.CompanyShareID,
					KittasApplied:  profile.DefaultKittas,
					BankID:         profile.DefaultBankID,
					Status:         "pending",
					AppliedAt:      time.Now(),
				}
				db.Create(&app)

				// Process application
				go processIPOApplication(&app, &profile)
			}
		}

		<-ticker.C
	}
}

// Helper function to safely get string from interface map
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
		return fmt.Sprintf("%v", val)
	}
	return ""
}
