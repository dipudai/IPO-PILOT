package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Nepal Payment Gateways Integration
// Supports: eSewa, Khalti, ConnectIPS

// eSewa Integration
type EsewaPaymentRequest struct {
	Amount      string `json:"amount"`
	FailureURL  string `json:"failure_url"`
	ProductCode string `json:"product_code"`
	RefundURL   string `json:"refund_url"`
	ServiceCode string `json:"service_code"`
	SuccessURL  string `json:"success_url"`
	TaxAmount   string `json:"tax_amount"`
	TotalAmount string `json:"total_amount"`
	TransactionUUID string `json:"transaction_uuid"`
}

// Khalti Integration
type KhaltiPaymentRequest struct {
	PublicKey    string `json:"public_key"`
	Amount       int64  `json:"amount"` // In paisa (amount in NPR * 100)
	ProductName  string `json:"product_name"`
	ProductID    string `json:"product_id"`
	Returner     string `json:"returner"`
	Website      string `json:"website"`
	MerchantName string `json:"merchant_name"`
}

// Nepal Payment Gateway Configuration
type NepalPaymentConfig struct {
	EsewaEnabled      bool
	EsewaServiceCode  string
	EsewaProductCode  string
	EsewaSecret       string
	
	KhaltiEnabled     bool
	KhaltiPublicKey   string
	KhaltiSecretKey   string
	
	ConnectIPSEnabled bool
	ConnectIPSURL     string
	ConnectIPSMerchantCode string
}

// Initialize Nepal payment config from environment
func getNepalPaymentConfig() NepalPaymentConfig {
	return NepalPaymentConfig{
		EsewaEnabled:     true,
		EsewaServiceCode: "EPAYTEST",
		EsewaProductCode: "IPO-PILOT",
		EsewaSecret:      "8gBm/:&EnhH.1/q",  // Get from eSewa
		
		KhaltiEnabled:    true,
		KhaltiPublicKey:  "test_public_key_dc74e0fd57cb46cd93832722edca97c3",
		KhaltiSecretKey:  "test_secret_key_dc74e0fd57cb46cd93832722edca97c3",
		
		ConnectIPSEnabled: false,
		ConnectIPSURL:    "https://connectips.com/api/",
		ConnectIPSMerchantCode: "YOUR_MERCHANT_CODE",
	}
}

// eSewa Initiate Payment
func initiateEsewaPayment(c *gin.Context, subscriptionID uint, amount string, userEmail string) (string, error) {
	config := getNepalPaymentConfig()
	
	transactionUUID := fmt.Sprintf("IPO-%d-%d", subscriptionID, time.Now().Unix())
	
	esewaReq := EsewaPaymentRequest{
		Amount:          amount,
		FailureURL:      "http://localhost:8080/payment/esewa/failure",
		ProductCode:     config.EsewaProductCode,
		RefundURL:       "http://localhost:8080/payment/esewa/refund",
		ServiceCode:     config.EsewaServiceCode,
		SuccessURL:      "http://localhost:8080/payment/esewa/success",
		TaxAmount:       "0",
		TotalAmount:     amount,
		TransactionUUID: transactionUUID,
	}
	
	// Create eSewa request URL
	params := url.Values{}
	params.Add("amount", esewaReq.Amount)
	params.Add("failure_url", esewaReq.FailureURL)
	params.Add("product_code", esewaReq.ProductCode)
	params.Add("refund_url", esewaReq.RefundURL)
	params.Add("service_code", esewaReq.ServiceCode)
	params.Add("success_url", esewaReq.SuccessURL)
	params.Add("tax_amount", esewaReq.TaxAmount)
	params.Add("total_amount", esewaReq.TotalAmount)
	params.Add("transaction_uuid", esewaReq.TransactionUUID)
	
	esewaURL := "https://rc-epay.esewa.com.np/api/epay/initiate/?" + params.Encode()
	
	return esewaURL, nil
}

// Khalti Initiate Payment
func initiateKhaltiPayment(c *gin.Context, subscriptionID uint, amountNPR string) (map[string]interface{}, error) {
	config := getNepalPaymentConfig()
	
	amountInt, _ := strconv.ParseInt(strings.TrimSpace(amountNPR), 10, 64)
	amountInPaisa := amountInt * 100 // Convert to paisa
	
	khaltiReq := KhaltiPaymentRequest{
		PublicKey:   config.KhaltiPublicKey,
		Amount:      amountInPaisa,
		ProductName: "IPO Pilot Subscription",
		ProductID:   fmt.Sprintf("IPO-%d", subscriptionID),
		Returner:    "User",
		Website:     "http://localhost:8080",
		MerchantName: "IPO Pilot Nepal",
	}
	
	data, _ := json.Marshal(khaltiReq)
	
	return map[string]interface{}{
		"publicKey": khaltiReq.PublicKey,
		"amount":    khaltiReq.Amount,
		"productName": khaltiReq.ProductName,
		"productId": khaltiReq.ProductID,
		"payload": string(data),
	}, nil
}

// eSewa Payment Success Handler
func esewaSuccessHandler(c *gin.Context) {
	transactionUUID := c.Query("transaction_uuid")
	status := c.Query("status")
	refID := c.Query("refId")
	
	if status != "COMPLETE" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Payment verification failed",
		})
		return
	}
	
	// Verify payment with eSewa
	config := getNepalPaymentConfig()
	verifyData := fmt.Sprintf("total_amount=%s,transaction_uuid=%s,product_code=%s,%s",
		c.Query("total_amount"), transactionUUID, config.EsewaProductCode, config.EsewaSecret)
	
	hash := md5.Sum([]byte(verifyData))
	hashStr := fmt.Sprintf("%x", hash)
	
	if c.Query("signature") != hashStr {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Payment signature verification failed",
		})
		return
	}
	
	// Update subscription as paid
	log.Printf("eSewa Payment Success - TxnID: %s, RefID: %s\n", transactionUUID, refID)
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payment successful! Your subscription is now active.",
		"refId": refID,
		"redirectURL": "/dashboard",
	})
}

// Khalti Payment Success Handler
func khaltiSuccessHandler(c *gin.Context) {
	tokenID := c.PostForm("token")
	amount := c.PostForm("amount")
	
	if tokenID == "" || amount == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid payment data",
		})
		return
	}
	
	// Verify with Khalti
	config := getNepalPaymentConfig()
	
	req, _ := http.NewRequest("POST", "https://khalti.com/api/v2/payment/verify/", nil)
	req.Header.Add("Authorization", "Key "+config.KhaltiSecretKey)
	
	payload := url.Values{}
	payload.Add("token", tokenID)
	payload.Add("amount", amount)
	
	req.Body = io.NopCloser(strings.NewReader(payload.Encode()))
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Payment verification failed",
		})
		return
	}
	defer resp.Body.Close()
	
	var khaltiResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&khaltiResp)
	
	if khaltiResp["state"].(map[string]interface{})["name"] != "Complete" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Payment not complete",
		})
		return
	}
	
	log.Printf("Khalti Payment Success - Token: %s, Amount: %s\n", tokenID, amount)
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payment successful! Your subscription is now active.",
		"redirectURL": "/dashboard",
	})
}

// ConnectIPS Payment Handler (Future)
func connectIPSHandler(c *gin.Context) {
	// ConnectIPS integration for bank transfers
	c.JSON(http.StatusOK, gin.H{
		"message": "ConnectIPS integration coming soon",
	})
}

// Nepal Payment Handler - Unified entry point
func nepaliPaymentHandler(c *gin.Context) {
	var input struct {
		SubscriptionID uint   `json:"subscription_id"`
		Amount         string `json:"amount"`
		PaymentMethod  string `json:"payment_method"` // esewa, khalti, connectips
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	
	switch input.PaymentMethod {
	case "esewa":
		esewaURL, err := initiateEsewaPayment(c, input.SubscriptionID, input.Amount, "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate eSewa payment"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"url": esewaURL,
			"method": "redirect",
		})
		
	case "khalti":
		khaltiData, err := initiateKhaltiPayment(c, input.SubscriptionID, input.Amount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate Khalti payment"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": khaltiData,
			"method": "widget",
		})
		
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported payment method"})
	}
}
