package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	crand "crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your-secret-key-change-in-production")

// JWT Claims
type Claims struct {
	UserID  uint `json:"user_id"`
	IsAdmin bool `json:"is_admin"`
	jwt.RegisteredClaims
}

// Generate JWT token
func generateJWT(userID uint, isAdmin bool) (string, error) {
	claims := &Claims{
		UserID:  userID,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Validate JWT token
func validateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// Hash password using bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Check password hash
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Generate encryption key from name
func generateEncryptionKey(name string) string {
	hash := md5.Sum([]byte(name + "IPO_Pilot_Secret"))
	return fmt.Sprintf("%x", hash)[:16]
}

// Encrypt using AES
func encryptAES(key, text string) (string, error) {
	keyBytes := []byte(key)
	if len(keyBytes) != 16 && len(keyBytes) != 24 && len(keyBytes) != 32 {
		// Pad or truncate key to 16 bytes
		if len(keyBytes) < 16 {
			keyBytes = append(keyBytes, make([]byte, 16-len(keyBytes))...)
		} else {
			keyBytes = keyBytes[:16]
		}
	}

	block, err := aes.NewCipher(keyBytes)
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

// Decrypt using AES
func decryptAES(key, cryptoText string) (string, error) {
	keyBytes := []byte(key)
	if len(keyBytes) != 16 && len(keyBytes) != 24 && len(keyBytes) != 32 {
		// Pad or truncate key to 16 bytes
		if len(keyBytes) < 16 {
			keyBytes = append(keyBytes, make([]byte, 16-len(keyBytes))...)
		} else {
			keyBytes = keyBytes[:16]
		}
	}

	block, err := aes.NewCipher(keyBytes)
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
