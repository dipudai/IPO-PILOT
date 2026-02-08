#!/bin/bash

echo "=== IPO Pilot Login Test ==="
echo ""

# Test 1: Register user
echo "1️⃣ Testing Registration..."
REG_RESPONSE=$(curl -s -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Debug User","email":"debug@test.org","password":"Debug@12345"}')
echo "Response: $REG_RESPONSE"
echo ""

# Test 2: Login
echo "2️⃣ Testing Login..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"debug@test.org","password":"Debug@12345"}')
echo "Response: $LOGIN_RESPONSE"
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | head -1 | cut -d'"' -f4)
echo "Extracted Token: ${TOKEN:0:50}..."
echo ""

# Test 3: Access dashboard with auth header
echo "3️⃣ Testing Dashboard Access (with Authorization header)..."
DASH_HEADER=$(curl -s -i -H "Authorization: Bearer $TOKEN" http://localhost:8080/dashboard | head -20)
echo "$DASH_HEADER"
echo ""

# Test 4: Access dashboard with cookie
echo "4️⃣ Testing Dashboard Access (with Cookie)..."
DASH_COOKIE=$(curl -s -i -b "auth_token=$TOKEN" http://localhost:8080/dashboard | head -20)
echo "$DASH_COOKIE"
