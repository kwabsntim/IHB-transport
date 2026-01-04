#!/bin/bash

TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImlzYWFiYWlkb29AeWFob28uY29tIiwiZXhwIjoxNzY3NjI0MTA5LCJyb2xlIjoiYWRtaW4ifQ.VmHwyCnOPAjeH6bSF7VEwISA0RH851PFnn-qX5l0tlY"

echo "========================================="
echo "TESTING DECLINE FLOW"
echo "========================================="
echo ""

echo "1️⃣  Creating New Delivery Request..."
RESPONSE=$(curl -s -X POST http://localhost:8080/api/public/deliveries \
  -H "Content-Type: application/json" \
  -d '{"client_name":"Decline Tester","client_email":"decline@example.com","pickup_address":"Start","dropoff_address":"End","item_description":"Test","weight":2.0}')
DECLINE_ID=$(echo $RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "✅ Created delivery: $DECLINE_ID"
echo ""

sleep 1

echo "2️⃣  Admin Sets Price (999.99 DKK - Too expensive!)..."
curl -s -X POST http://localhost:8080/api/admin/deliveries/$DECLINE_ID/price \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"price": 999.99}'
echo ""
echo ""

sleep 1

echo "3️⃣  Client DECLINES the Price..."
curl -s -X POST http://localhost:8080/api/public/deliveries/$DECLINE_ID/decline \
  -H "Content-Type: application/json" \
  -d '{"reason": "Too expensive, cannot afford it"}'
echo ""
echo ""

sleep 1

echo "4️⃣  Check Final Status..."
curl -s http://localhost:8080/api/public/deliveries/$DECLINE_ID | grep -o '"status":"[^"]*' | head -1
echo ""

echo "5️⃣  Check Decline Reason..."
curl -s http://localhost:8080/api/public/deliveries/$DECLINE_ID | grep -o '"decline_reason":"[^"]*'
echo ""
echo ""

echo "========================================="
echo ""
echo "TESTING OTHER ENDPOINTS"
echo "========================================="
echo ""

echo "📋 Get All Deliveries (Admin)..."
curl -s http://localhost:8080/api/admin/deliveries \
  -H "Authorization: Bearer $TOKEN" | grep -o '"count":[0-9]*'
echo ""
echo ""

echo "📊 Get Deliveries by Status (DECLINED)..."
curl -s "http://localhost:8080/api/admin/deliveries/status?status=DECLINED" \
  -H "Authorization: Bearer $TOKEN" | grep -o '"count":[0-9]*'
echo ""
echo ""

echo "📊 Get Deliveries by Status (DELIVERED)..."
curl -s "http://localhost:8080/api/admin/deliveries/status?status=DELIVERED" \
  -H "Authorization: Bearer $TOKEN" | grep -o '"count":[0-9]*'
echo ""
echo ""

echo "========================================="
echo "✅ ALL TESTS COMPLETED"
echo "========================================="
