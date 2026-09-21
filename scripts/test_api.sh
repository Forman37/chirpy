#!/usr/bin/env bash

BASE_URL="http://localhost:8080"
USER_ID="550e8400-e29b-41d4-a716-446655440000"

echo "========================================"
echo "RESET DATABASE"
echo "========================================"

curl -i -X POST "$BASE_URL/admin/reset"

echo
echo
echo "========================================"
echo "CREATE USER 1"
echo "========================================"

curl -i -X POST "$BASE_URL/api/users" \
    -H "Content-Type: application/json" \
    -d '{"email":"test@example.com"}'

echo
echo
echo "========================================"
echo "CREATE USER 2"
echo "========================================"

curl -i -X POST "$BASE_URL/api/users" \
    -H "Content-Type: application/json" \
    -d '{"email":"bob@example.com"}'

echo
echo
echo "========================================"
echo "INVALID JSON EOF"
echo "========================================"

curl -i -X POST "$BASE_URL/api/users" \
    -H "Content-Type: application/json" \
    -d '{"email":'

echo
echo
echo "========================================"
echo "CREATE CHIRP 1"
echo "========================================"

curl -i -X POST "$BASE_URL/api/chirps" \
    -H "Content-Type: application/json" \
    -d "{\"body\": \"Hello, world!\", \"user_id\": \"$USER_ID\"}"

echo
echo
echo "========================================"
echo "CREATE CHIRP 2 - ERROR : TOO LONG"
echo "========================================"

curl -i -X POST "$BASE_URL/api/chirps" \
    -H "Content-Type: application/json" \
    -d "{\"body\": \"Hello, world! This post will very likely be over 140 characters, but I really have something I want to say and to get off my chest. Since the beginning of time, I think rocks really have been very against us, and heres why...\", \"user_id\": \"$USER_ID\"}"

echo
echo
echo "========================================"
echo "CREATE CHIRP 3 - ERROR : MISSING BODY"
echo "========================================"

curl -i -X POST "$BASE_URL/api/chirps" \
    -H "Content-Type: application/json" \
    -d "{\"user_id\": \"$USER_ID\"}"

echo
echo
echo "========================================"
echo "CREATE CHIRP 4 - ERROR : MISSING USERID"
echo "========================================"

curl -i -X POST "$BASE_URL/api/chirps" \
    -H "Content-Type: application/json" \
    -d "{\"body\": \"Hello, world!\"}"

echo
echo
echo "========================================"
echo "DONE"
echo "========================================"
