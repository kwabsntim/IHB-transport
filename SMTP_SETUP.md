# SMTP Email Setup Guide

This guide will help you configure real email sending for IHB Transport using SMTP.

## 🎯 Quick Overview

The application supports multiple SMTP providers. Choose the one that works best for you:

| Provider | Free Tier | Recommended For | Setup Difficulty |
|----------|-----------|-----------------|------------------|
| **Gmail** | 500/day | Testing & Personal | Easy |
| **SMTP2GO** | 1,000/month | Small Business | Easy |
| **SendGrid** | 100/day | Production | Medium |
| **Mailgun** | 5,000/month (3 months) | Production | Medium |
| **AWS SES** | 62,000/month (free with EC2) | Enterprise | Hard |

## Option 1: Gmail (Easiest for Testing)

### Step 1: Enable 2-Factor Authentication
1. Go to [Google Account Security](https://myaccount.google.com/security)
2. Enable **2-Step Verification** if not already enabled

### Step 2: Create App Password
1. Go to [App Passwords](https://myaccount.google.com/apppasswords)
2. Select App: **Mail**
3. Select Device: **Other (Custom name)** → Enter "IHB Transport"
4. Click **Generate**
5. Copy the 16-character password (spaces will be removed automatically)

### Step 3: Update .env
```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=youremail@gmail.com
SMTP_PASSWORD=your-16-char-app-password
SMTP_FROM_EMAIL=youremail@gmail.com
SMTP_FROM_NAME=IHB Transport
```

### Limitations:
- ⚠️ 500 emails per day limit
- ⚠️ May show "sent via gmail.com" in recipient inbox
- ✅ Perfect for testing and development

---

## Option 2: SMTP2GO (Best Free Tier)

### Step 1: Create Account
1. Go to [SMTP2GO](https://www.smtp2go.com)
2. Sign up for free account
3. Verify your email address

### Step 2: Get SMTP Credentials
1. Go to **Settings** → **Users**
2. Click **Add SMTP User**
3. Enter username (e.g., `ihb-transport`)
4. Copy the generated password
5. Save the credentials

### Step 3: Update .env
```bash
SMTP_HOST=mail.smtp2go.com
SMTP_PORT=2525
SMTP_USER=ihb-transport
SMTP_PASSWORD=your-smtp2go-password
SMTP_FROM_EMAIL=noreply@ihbtransport.com
SMTP_FROM_NAME=IHB Transport
```

### Benefits:
- ✅ 1,000 emails per month free
- ✅ Custom "from" address
- ✅ Detailed analytics dashboard
- ✅ Better deliverability than Gmail

---

## Option 3: SendGrid (Production Ready)

### Step 1: Create Account
1. Go to [SendGrid](https://sendgrid.com)
2. Sign up for free account
3. Complete sender verification

### Step 2: Create API Key
1. Go to **Settings** → **API Keys**
2. Click **Create API Key**
3. Name: `IHB Transport`
4. Permissions: **Full Access** or **Restricted (Mail Send only)**
5. Copy the API key (starts with `SG.`)

### Step 3: Update .env
```bash
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASSWORD=SG.your-sendgrid-api-key-here
SMTP_FROM_EMAIL=noreply@ihbtransport.com
SMTP_FROM_NAME=IHB Transport
```

### Benefits:
- ✅ 100 emails per day free forever
- ✅ Professional email infrastructure
- ✅ Advanced analytics
- ✅ Better deliverability

---

## Option 4: Mailgun (Enterprise Grade)

### Step 1: Create Account
1. Go to [Mailgun](https://www.mailgun.com)
2. Sign up for free account
3. Verify your phone number

### Step 2: Get SMTP Credentials
1. Go to **Sending** → **Domain Settings**
2. Click on your sandbox domain (or add custom domain)
3. Find **SMTP Credentials** section
4. Copy the username and password

### Step 3: Update .env
```bash
SMTP_HOST=smtp.mailgun.org
SMTP_PORT=587
SMTP_USER=postmaster@sandboxXXXX.mailgun.org
SMTP_PASSWORD=your-mailgun-password
SMTP_FROM_EMAIL=noreply@ihbtransport.com
SMTP_FROM_NAME=IHB Transport
```

### Benefits:
- ✅ 5,000 emails/month for 3 months (then $35/month or $0.80/1000 emails)
- ✅ Excellent deliverability
- ✅ Webhook support
- ✅ Email validation API

---

## Testing Your SMTP Setup

### 1. Restart Your Server
```bash
# Stop the server if running
# Then start again to load new env vars
go run cmd/server/main.go
```

You should see:
```
✅ SMTP enabled - using smtp.gmail.com:587
```

Or if not configured:
```
⚠️  SMTP not configured - emails will be logged to console only
```

### 2. Test with Delivery Request
```bash
# Create a test delivery request
curl -X POST http://localhost:8080/api/public/deliveries \
  -H "Content-Type: application/json" \
  -d '{
    "client_name": "Test User",
    "client_email": "YOUR_REAL_EMAIL@gmail.com",
    "pickup_address": "123 Start St",
    "dropoff_address": "456 End Ave",
    "item_description": "Test Package",
    "weight": 5.0
  }'
```

Check your email inbox for "Delivery Request Received" email!

### 3. Test Price Email
```bash
# Login to get token
TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"isaabaidoo@yahoo.com","password":"YOUR_ADMIN_PASSWORD"}' \
  | grep -o '"token":"[^"]*' | cut -d'"' -f4)

# Get the delivery ID from previous test
DELIVERY_ID="paste-id-here"

# Set price
curl -X POST http://localhost:8080/api/admin/deliveries/$DELIVERY_ID/price \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"price": 299.99}'
```

Check your email for "Price Quote" email!

---

## Email Templates Included

Your application now sends beautiful HTML emails for:

1. **Request Received** - Confirmation when delivery is created
2. **Price Quote** - Professional price quote with clear call-to-action
3. **Delivery Accepted** - Confirmation and next steps
4. **Delivery Declined** - Thank you message + admin notification
5. **Driver En Route** - Package picked up notification
6. **Delivery Complete** - Success confirmation

---

## Troubleshooting

### "Failed to connect to SMTP server"

**Gmail:**
- ✅ Check if 2FA is enabled
- ✅ Use App Password, not regular password
- ✅ Remove spaces from app password

**Other providers:**
- ✅ Verify host and port are correct
- ✅ Check firewall isn't blocking port 587
- ✅ Try port 2525 if 587 doesn't work

### "Authentication failed"

- ✅ Double-check username and password
- ✅ For SendGrid, username must be exactly `apikey`
- ✅ Ensure no extra spaces in .env file

### "Recipient rejected"

**Gmail:**
- ✅ You can only send to verified addresses initially
- ✅ Build sending reputation gradually

**Mailgun Sandbox:**
- ✅ Add recipient email to "Authorized Recipients" list
- ✅ Or verify a custom domain

### Emails going to spam

- ✅ Add SPF record to your domain
- ✅ Add DKIM signature (provider handles this)
- ✅ Use a verified custom domain
- ✅ Avoid spammy words in subject/body

---

## Production Recommendations

### 1. Use Custom Domain
Instead of `noreply@ihbtransport.com`, use a real domain you own:
```bash
SMTP_FROM_EMAIL=noreply@yourdomain.com
```

### 2. Verify Domain with Provider
Most providers offer domain verification:
- SendGrid: Settings → Sender Authentication
- Mailgun: Sending → Domains → Add Domain
- SMTP2GO: Settings → Sending Domains

### 3. Set Up SPF and DKIM
Add these DNS records (provider will give you values):
```
TXT record: v=spf1 include:sendgrid.net ~all
TXT record: DKIM signature (provided by email service)
```

### 4. Monitor Deliverability
- Check provider dashboard for bounce rates
- Keep bounce rate < 5%
- Remove invalid email addresses

### 5. Rate Limiting
Add rate limiting to prevent abuse:
```go
// Limit: 10 delivery requests per hour per email
// Implement in handlers/handlers.go
```

---

## Environment Variables Summary

Add these to your `.env` file:

```bash
# Required for SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# Optional (will use defaults if not set)
SMTP_FROM_EMAIL=noreply@ihbtransport.com  # default
SMTP_FROM_NAME=IHB Transport              # default
```

For deployment platforms (Render, Railway, etc.), add these as environment variables in the dashboard.

---

## Cost Comparison (After Free Tier)

| Provider | Cost |
|----------|------|
| Gmail | Free (500/day limit) |
| SMTP2GO | $10/month (10,000 emails) |
| SendGrid | $19.95/month (50,000 emails) |
| Mailgun | $35/month (50,000 emails) |
| AWS SES | $0.10 per 1,000 emails |

---

## Next Steps

1. ✅ Choose your SMTP provider
2. ✅ Set up account and get credentials
3. ✅ Update `.env` file with SMTP settings
4. ✅ Restart server
5. ✅ Test with real delivery request
6. ✅ Deploy to production with SMTP configured

**Need help?** Check provider documentation or contact their support.

---

**Happy emailing! 📧**
