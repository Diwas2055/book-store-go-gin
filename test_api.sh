#!/bin/bash
set -e

echo "Starting API tests..."
echo "---------------------"

# Base URL
BASE_URL=http://localhost:5000

# Test variables
BOOK_ID=""
USER_ID=""

# Create a user
echo "Testing CREATE user endpoint..."
CREATE_USER_RESPONSE=$(curl -s -X POST ${BASE_URL}/users \
  -H "Content-Type: application/json" \
  -d '{
        "name":"Alice", "email":"alice@example.com"
  }')
echo "Response: ${CREATE_USER_RESPONSE}"
# Extract book ID (this will depend on your API response format)
USER_ID=$(echo "$CREATE_USER_RESPONSE" | jq -r '.data.id')
if [ -z "$USER_ID" ]; then
  echo "Failed to extract user ID. Exiting."
  exit 1
fi

echo "Extracted User ID: $USER_ID"
echo ""

# Create a book
echo "Testing CREATE book endpoint..."
CREATE_RESPONSE=$(curl -s -X POST ${BASE_URL}/books \
  -H "Content-Type: application/json" \
  -H "X-User-ID: $USER_ID" \
  -d '{
    "title": "Test Book",
    "author": "Test Author"
  }')
echo "Response: ${CREATE_RESPONSE}"
# Extract book ID (this will depend on your API response format)
BOOK_ID=$(echo "$CREATE_RESPONSE" | jq -r '.data.id')
echo "Created book with ID: ${BOOK_ID}"
echo ""

# List all books
echo "Testing GET all books endpoint..."
curl -s -X GET ${BASE_URL}/books -H "X-User-ID: $USER_ID"
echo -e "\n"

# Get specific book
echo "Testing GET specific book endpoint..."
curl -s -X GET ${BASE_URL}/books/${BOOK_ID} -H "X-User-ID: $USER_ID"
echo -e "\n"

# Update book
echo "Testing UPDATE book endpoint..."
curl -s -X PUT ${BASE_URL}/books/${BOOK_ID} \
  -H "Content-Type: application/json" \
  -H "X-User-ID: $USER_ID" \
  -d '{
    "title": "Updated Test Book",
    "author": "Test Author Updated"
  }'
echo -e "\n"

# Get updated book
echo "Getting updated book..."
curl -s -X GET ${BASE_URL}/books/${BOOK_ID} -H "X-User-ID: $USER_ID"
echo -e "\n"

# Get specific user
echo "Testing GET specific user endpoint..."
curl -s -X GET ${BASE_URL}/users/${USER_ID}
echo -e "\n"

# Delete book
echo "Testing DELETE book endpoint..."
curl -s -X DELETE ${BASE_URL}/books/${BOOK_ID} -H "X-User-ID: $USER_ID"
echo -e "\n"

# Verify deletion
echo "Verifying book was deleted..."
curl -s -X GET ${BASE_URL}/books/${BOOK_ID} -H "X-User-ID: $USER_ID"
echo -e "\n"
# Delete users
echo "Testing Delete user endpoint"
curl -s -X DELETE ${BASE_URL}/users/${USER_ID}
echo -e "\n"

echo "API tests completed."
