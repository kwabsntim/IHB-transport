#!/bin/bash

# SMTP Email Testing Script
# Make sure to update REAL_EMAIL with your actual email address

REAL_EMAIL="your-email@gmail.com"  # ⚠️ UPDATE THIS WITH YOUR REAL EMAIL!

if [ "$REAL_EMAIL" = "your-email@gmail.com" ]; then
    echo "❌ ERROR: Please update REAL_EMAIL in this script with your actual email address!"
    exit 1
fi

echo "========================================="
echo "🧪 SMTP EMAIL TESTING SCRIPT"
echo "========================================="
echo ""
echo "This script will test all email notifications"
echo "Emails will be sent to: $REAL_EMAIL"
echo ""
read -p "Press Enter to continue..."
echo ""

# Login to get admin token
echo "1️⃣  Logging in as admin..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"isaabaidoo@yahoo.com","password":"Ihbtransport@424"}')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "❌ Login failed! Check admin credentials."
    exit 1
fi

echo "✅ Login successful"
echo ""

# Test 1: Create Delivery Request (sends REQUEST_RECEIVED email)
echo "========================================="
echo "📧 TEST 1: REQUEST RECEIVED EMAIL"
echo "========================================="
echo "Creating delivery request..."
CREATE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/public/deliveries \
  -H "Content-Type: application/json" \
  -d "{
    \"client_name\": \"Email Test User\",
    \"client_email\": \"$REAL_EMAIL\",
    \"pickup_address\": \"123 Test Street, Copenhagen\",
    \"dropoff_address\": \"456 Delivery Avenue, Aarhus\",
    \"item_description\": \"SMTP Test Package\",
    \"weight\": 5.0
  }")

DELIVERY_ID=$(echo $CREATE_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -z "$DELIVERY_ID" ]; then
    echo "❌ Failed to create delivery"
    echo "Response: $CREATE_RESPONSE"
    exit 1
fi

echo "✅ Delivery created: $DELIVERY_ID"
echo "📬 Check your inbox ($REAL_EMAIL) for 'Delivery Request Received' email"
echo ""
read -p "Did you receive the email? (Press Enter to continue)"
echo ""

sleep 2

# Test 2: Set Price (sends PRICE_SENT email)
echo "========================================="
echo "📧 TEST 2: PRICE QUOTE EMAIL"
echo "========================================="
echo "Setting delivery price..."
curl -s -X POST http://localhost:8080/api/admin/deliveries/$DELIVERY_ID/price \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"price": 599.99}' > /dev/null

echo "✅ Price set: 599.99 DKK"
echo "📬 Check your inbox for 'Price Quote' email"
echo ""
read -p "Did you receive the email? (Press Enter to continue)"
echo ""

sleep 2

# Test 3: Accept Price (sends ACCEPTED email)
echo "========================================="
echo "📧 TEST 3: DELIVERY ACCEPTED EMAIL"
echo "========================================="
echo "Accepting delivery price..."
curl -s -X POST http://localhost:8080/api/public/deliveries/$DELIVERY_ID/accept \
  -H "Content-Type: application/json" > /dev/null

echo "✅ Price accepted"
echo "📬 Check your inbox for 'Delivery Confirmed' email"
echo ""
read -p "Did you receive the email? (Press Enter to continue)"
echo ""

sleep 2

# Test 4: Mark as Picked Up (sends DRIVER_ON_WAY email)
echo "========================================="
echo "📧 TEST 4: DRIVER EN ROUTE EMAIL"
echo "========================================="
echo "Marking package as picked up..."
curl -s -X POST http://localhost:8080/api/driver/deliveries/$DELIVERY_ID/pickup \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" > /dev/null

echo "✅ Package picked up by driver"
echo "📬 Check your inbox for 'Driver En Route' email"
echo ""
read -p "Did you receive the email? (Press Enter to continue)"
echo ""

sleep 2

# Test 5: Mark as Delivered (sends DELIVERED email)
echo "========================================="
echo "📧 TEST 5: DELIVERY COMPLETED EMAIL"
echo "========================================="
echo "Marking delivery as complete..."
curl -s -X POST http://localhost:8080/api/driver/deliveries/$DELIVERY_ID/complete \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" > /dev/null

echo "✅ Delivery completed"
echo "📬 Check your inbox for 'Delivery Completed' email"
echo ""
read -p "Did you receive the email? (Press Enter to continue)"
echo ""

# Test 6: Test Decline Flow
echo "========================================="
echo "📧 TEST 6: PRICE DECLINED EMAIL"
echo "========================================="
echo "Creating another delivery for decline test..."
DECLINE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/public/deliveries \
  -H "Content-Type: application/json" \
  -d "{
    \"client_name\": \"Decline Test User\",
    \"client_email\": \"$REAL_EMAIL\",
    \"pickup_address\": \"789 Decline St\",
    \"dropoff_address\": \"321 Reject Ave\",
    \"item_description\": \"Test Decline Package\",
    \"weight\": 3.0
  }")

DECLINE_ID=$(echo $DECLINE_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -z "$DECLINE_ID" ]; then
    echo "❌ Failed to create delivery for decline test"
    exit 1
fi

echo "✅ Second delivery created: $DECLINE_ID"
sleep 1

echo "Setting price..."
curl -s -X POST http://localhost:8080/api/admin/deliveries/$DECLINE_ID/price \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"price": 999.99}' > /dev/null

sleep 1

echo "Declining the price..."
curl -s -X POST http://localhost:8080/api/public/deliveries/$DECLINE_ID/decline \
  -H "Content-Type: application/json" \
  -d '{"reason": "Price is too high for my budget"}' > /dev/null

echo "✅ Price declined"
echo "📬 Check your inbox for 'Request Declined' email"
echo "📬 Admin should also receive notification at: isaabaidoo@yahoo.com"
echo ""

# Summary
echo ""
echo "========================================="
echo "✅ ALL EMAIL TESTS COMPLETED!"
echo "========================================="
echo ""
echo "You should have received 6 emails total:"
echo "  1. ✉️  Delivery Request Received"
echo "  2. ✉️  Price Quote (599.99 DKK)"
echo "  3. ✉️  Delivery Confirmed"
echo "  4. ✉️  Driver En Route"
echo "  5. ✉️  Delivery Completed"
echo "  6. ✉️  Request Declined (for second delivery)"
echo ""
echo "Admin ($ADMIN_EMAIL) should receive:"
echo "  7. ✉️  Price Declined Notification"
echo ""
echo "========================================="
echo ""
echo "If you didn't receive emails, check:"
echo "  - SMTP configuration in .env"
echo "  - Spam/junk folder"
echo "  - Server logs for error messages"
echo "  - Provider dashboard (SendGrid, Mailgun, etc.)"
echo ""
echo "Test delivery IDs for cleanup:"
echo "  - $DELIVERY_ID (completed)"
echo "  - $DECLINE_ID (declined)"
echo ""
