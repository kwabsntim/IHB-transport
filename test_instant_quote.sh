#!/bin/bash

# Test instant quote endpoint
echo "Testing Instant Quote Endpoint..."
echo ""

curl -X POST http://localhost:8080/api/public/instant-quote \
  -H "Content-Type: application/json" \
  -d '{
    "pickup_point": "Copenhagen Central Station",
    "delivery_address": "Aarhus City Hall",
    "weight": "25kg",
    "client_email": "test@example.com"
  }' | jq .

echo ""
echo "Check your logs for email sending status!"
