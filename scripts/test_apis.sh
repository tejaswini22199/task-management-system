#!/bin/bash

# Base URL
AUTH_URL="http://localhost:8000"
TASKS_URL="http://localhost:8001"

# Register two users
echo "📝 Registering users..."
USER1=$(curl -s -X POST "$AUTH_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d '{
          "username": "user1",
          "email": "user1@example.com",
          "password": "password123"
        }')

USER2=$(curl -s -X POST "$AUTH_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d '{
          "username": "user2",
          "email": "user2@example.com",
          "password": "password123"
        }')

echo "✅ Users registered!"

# Extract user IDs (if needed)
USER1_ID=$(echo $USER1 | jq -r '.user_id')
USER2_ID=$(echo $USER2 | jq -r '.user_id')

# Login users & get tokens
echo "🔑 Logging in users..."
TOKEN1=$(curl -s -X POST "$AUTH_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
          "email": "user1@example.com",
          "password": "password123"
        }' | jq -r '.token')

TOKEN2=$(curl -s -X POST "$AUTH_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
          "email": "user2@example.com",
          "password": "password123"
        }' | jq -r '.token')

echo "✅ Users logged in! Tokens acquired."

# Create tasks
echo "📝 Creating tasks..."
TASK1=$(curl -s -X POST "$TASKS_URL/tasks/" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN1" \
    -d '{
          "title": "Task 1 for User 1",
          "description": "Description of Task 1",
          "status": "pending",
          "user_ids": [1]
        }')

TASK2=$(curl -s -X POST "$TASKS_URL/tasks/" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN2" \
    -d '{
          "title": "Task 2 for User 2",
          "description": "Description of Task 2",
          "status": "pending",
          "user_ids": [2]
        }')

TASK3=$(curl -s -X POST "$TASKS_URL/tasks/" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN1" \
    -d '{
          "title": "Task 3 shared by User 1 and User 2",
          "description": "Shared task",
          "status": "pending",
          "user_ids": [1, 2]
        }')

echo "✅ Tasks created!"

# Fetch all tasks
echo "📋 Fetching all tasks..."
ALL_TASKS=$(curl -s -X GET "$TASKS_URL/tasks/" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN1")

echo "📌 Tasks List:"
echo "$ALL_TASKS" | jq

echo "🎉 API testing completed!"
