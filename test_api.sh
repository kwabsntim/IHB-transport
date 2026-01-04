#!/bin/bash

TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImlzYWFiYWlkb29AeWFob28uY29tIiwiZXhwIjoxNzY3NjI0MTA5LCJyb2xlIjoiYWRtaW4ifQ.VmHwyCnOPAjeH6bSF7VEwISA0RH851PFnn-qX5l0tlY"

echo "========================================="
echo "COMPLETE DELIVERY FLOW TEST"
echo "========================================="
echo ""

echo "1️⃣  Creating New Delivery Request..."
RESPONSE=$(curl -s -X POST http://localhost:8080/api/public/deliveries \
  -H "Content-Type: application/json" \
  -d '{"client_name":"Test User","client_email":"test@example.com","pickup_address":"Start Location","dropoff_address":"End Location","item_description":"Test Package","weight":5.0}')
NEW_ID=$(echo $RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "✅ Created delivery: $NEW_ID"
echo ""

sleep 1

echo "2️⃣  Admin Sets Price (199.99 DKK)..."
curl -s -X POST http://localhost:8080/api/admin/deliveries/$NEW_ID/price \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"price": 199.99}'
echo ""
echo ""

sleep 1

echo "3️⃣  Client Accepts Price..."
curl -s -X POST http://localhost:8080/api/public/deliveries/$NEW_ID/accept
echo ""
echo ""

sleep 1

echo "4️⃣  Driver Marks as Picked Up..."
curl -s -X POST http://localhost:8080/api/driver/deliveries/$NEW_ID/pickup \
  -H "Authorization: Bearer $TOKEN"
echo ""
echo ""

sleep 1

echo "5️⃣  Driver Marks as Delivered..."
curl -s -X POST http://localhost:8080/api/driver/deliveries/$NEW_ID/complete \
  -H "Authorization: Bearer $TOKEN"
echo ""
echo ""

echo "6️⃣  Final Delivery Status..."
curl -s http://localhost:8080/api/public/deliveries/$NEW_ID | grep -o '"status":"[^"]*'
echo ""
echo ""

echo "========================================="
echo "✅ COMPLETE FLOW TEST FINISHED"
echo "========================================="
