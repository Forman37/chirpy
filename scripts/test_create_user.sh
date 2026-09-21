#!/usr/bin/env bash

curl -X POST http://localhost:8080/api/users \
    -H "Content-Type: application/json" \
    -d '{"email": "test@example.com"}'
