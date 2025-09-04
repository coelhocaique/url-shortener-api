#!/bin/bash

echo "Testing URL Shortener API Error Handling..."
echo "=========================================="

BASE_URL="http://localhost:8080"

# Test 1: Invalid URL format (should return 400)
echo "1. Testing invalid URL format..."
RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" -X POST $BASE_URL/urls \
  -H "Content-Type: application/json" \
  -d '{"url": "caique is bad", "expiration_ms": 3600000}')

HTTP_STATUS=$(echo $RESPONSE | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
BODY=$(echo $RESPONSE | sed -e 's/HTTPSTATUS\:.*//g')

echo "Status: $HTTP_STATUS"
echo "Response: $BODY"
echo "Expected: 400"
echo ""

# Test 2: Invalid alias characters (should return 400)
echo "2. Testing invalid alias characters..."
RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" -X POST $BASE_URL/urls \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.example.com", "alias": "invalid@alias", "expiration_ms": 3600000}')

HTTP_STATUS=$(echo $RESPONSE | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
BODY=$(echo $RESPONSE | sed -e 's/HTTPSTATUS\:.*//g')

echo "Status: $HTTP_STATUS"
echo "Response: $BODY"
echo "Expected: 400"
echo ""

# Test 3: Valid request (should return 201)
echo "3. Testing valid request..."
RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" -X POST $BASE_URL/urls \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.example.com", "alias": "test-alias", "expiration_ms": 3600000}')

HTTP_STATUS=$(echo $RESPONSE | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
BODY=$(echo $RESPONSE | sed -e 's/HTTPSTATUS\:.*//g')

echo "Status: $HTTP_STATUS"
echo "Response: $BODY"
echo "Expected: 201"
echo ""

# Test 4: Duplicate alias (should return 409)
echo "4. Testing duplicate alias..."
RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" -X POST $BASE_URL/urls \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.example2.com", "alias": "test-alias", "expiration_ms": 3600000}')

HTTP_STATUS=$(echo $RESPONSE | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
BODY=$(echo $RESPONSE | sed -e 's/HTTPSTATUS\:.*//g')

echo "Status: $HTTP_STATUS"
echo "Response: $BODY"
echo "Expected: 409"
echo ""

echo "Error handling tests completed!"
