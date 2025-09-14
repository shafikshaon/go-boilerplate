#!/bin/bash

BASE_URL="http://localhost:8080"

echo "🧪 Testing API Endpoints"
echo "========================"

# Health check
echo "1. Health Check:"
curl -s "$BASE_URL/health" | jq .

echo -e "\n2. Register User:"
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123"
  }')
echo $REGISTER_RESPONSE | jq .

echo -e "\n3. Login:"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }')
echo $LOGIN_RESPONSE | jq .

# Extract token
TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.token')

echo -e "\n4. Get Current User Profile (Protected):"
if [ "$TOKEN" != "null" ] && [ "$TOKEN" != "" ]; then
    curl -s "$BASE_URL/api/v1/auth/me" \
      -H "Authorization: Bearer $TOKEN" | jq .
fi

echo -e "\n5. Get Users (cached in Redis):"
curl -s "$BASE_URL/api/v1/users" | jq .

echo -e "\n6. Get User by ID (first time - DB, second time - Redis cache):"
curl -s "$BASE_URL/api/v1/users/1" | jq .
echo -e "\n   Getting same user again (should be from cache):"
curl -s "$BASE_URL/api/v1/users/1" | jq .

if [ "$TOKEN" != "null" ] && [ "$TOKEN" != "" ]; then
    echo -e "\n7. Update User (Protected):"
    curl -s -X PUT "$BASE_URL/api/v1/users/1" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "name": "Updated Test User"
      }' | jq .

    echo -e "\n8. Logout (Protected):"
    curl -s -X POST "$BASE_URL/api/v1/auth/logout" \
      -H "Authorization: Bearer $TOKEN" | jq .
fi

echo -e "\n✅ API testing completed!"
echo ""
echo "🔗 Available endpoints:"
echo "  GET  /health                     - Health check"
echo "  POST /api/v1/auth/register      - Register user"
echo "  POST /api/v1/auth/login         - Login user"
echo "  GET  /api/v1/auth/me            - Get current user (Protected)"
echo "  POST /api/v1/auth/logout        - Logout user (Protected)"
echo "  GET  /api/v1/users              - Get all users"
echo "  GET  /api/v1/users/:id          - Get user by ID (cached)"
echo "  PUT  /api/v1/users/:id          - Update user (Protected)"
echo "  DELETE /api/v1/users/:id        - Delete user (Protected)"
