#!/bin/bash

# Test script for logout endpoint
echo "Testing logout endpoint..."

# First, let's login to get a token
echo "Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:3000/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' \
  -c cookies.txt)

echo "Login response: $LOGIN_RESPONSE"

# Now, let's test the logout endpoint
echo "Logging out..."
LOGOUT_RESPONSE=$(curl -s -X POST http://localhost:3000/auth/logout \
  -H "Content-Type: application/json" \
  -b cookies.txt)

echo "Logout response: $LOGOUT_RESPONSE"

# Clean up
rm -f cookies.txt

echo "Test completed."