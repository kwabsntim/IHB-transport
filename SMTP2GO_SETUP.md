# SMTP2GO Complete Setup Guide for IHB Transport

## 📊 Email Budget Planning

### Monthly Allocation (1,000 emails/month)
- **Daily Average**: 1,000 ÷ 30 days = ~33 emails/day
- **Safe Daily Limit**: **30 emails/day** (buffer for spikes)
- **Per-Client Limit**: **3 emails/day** (prevents abuse)

### Why These Limits?
1. **Startup Phase**: Low initial traffic
2. **Buffer for Spikes**: Weekend/holiday catch-up
3. **Anti-Abuse**: Prevents spam from single user
4. **Cost Control**: Stays within free tier

### Email Distribution Example:
```
10 deliveries/day × 6 emails each = 60 emails needed
But with rate limiting:
- Each client max 3 emails/day
- System max 30 emails/day
- Deliveries process over 2-3 days
- Perfect for startup pace! ✅
```

---

## 🚀 SMTP2GO Setup (Step-by-Step)

### Step 1: Create SMTP2GO Account

1. Go to [https://www.smtp2go.com](https://www.smtp2go.com)
2. Click **Sign Up Free**
3. Fill in your details:
   - Email: `isaabaidoo@yahoo.com` (or your preferred email)
   - Company: `IHB Transport`
   - Country: Denmark
4. Verify your email address (check inbox)

### Step 2: Create SMTP User

1. After login, go to **Settings** → **Users**
2. Click **Add SMTP User**
3. Enter username: `ihb-transport-api`
4. Click **Generate Password** (copy this immediately!)
5. Save the credentials:
   ```
   Username: ihb-transport-api
   Password: [generated password - save this!]
   ```

### Step 3: Configure Your Application

Update your `.env` file:

```bash
# SMTP2GO Configuration
SMTP_HOST=mail.smtp2go.com
SMTP_PORT=2525
SMTP_USER=ihb-transport-api
SMTP_PASSWORD=your-generated-password-here
SMTP_FROM_EMAIL=noreply@ihbtransport.com
SMTP_FROM_NAME=IHB Transport
```

**Port Options:**
- `2525` - Recommended (works everywhere)
- `587` - Alternative (standard SMTP)
- `8025` - Alternative (firewall-friendly)
- `80` - Last resort (if others blocked)

### Step 4: Test Connection

```bash
# Restart your server
go run cmd/server/main.go

# You should see:
✅ SMTP enabled - using mail.smtp2go.com:2525
📊 Rate limiter initialized: 30 emails/day, 3 per client
```

### Step 5: Send Test Email

```bash
# Create a delivery request with YOUR real email
curl -X POST http://localhost:8080/api/public/deliveries \
  -H "Content-Type: application/json" \
  -d '{
    "client_name": "Test User",
    "client_email": "YOUR_EMAIL@gmail.com",
    "pickup_address": "Copenhagen Central Station",
    "dropoff_address": "Aarhus Station",
    "item_description": "Test Package",
    "weight": 5.0
  }'

# Check your email inbox! 📧
```

---

## 📊 Rate Limiting System

### How It Works

**Daily Limit: 30 emails**
- Resets at midnight (UTC)
- Prevents exceeding free tier
- Buffer for unexpected spikes

**Per-Client Limit: 3 emails/day**
- Prevents spam/abuse
- One client can't use all quota
- Fair distribution

### Email Flow with Limits

**Scenario 1: Normal Delivery**
```
Client creates delivery    → Email 1 ✅ (1/3 for client, 1/30 total)
Admin sets price          → Email 2 ✅ (2/3 for client, 2/30 total)
Client accepts            → Email 3 ✅ (3/3 for client, 3/30 total)
Driver picks up           → ❌ BLOCKED (client limit reached)
Delivery complete         → ❌ BLOCKED (client limit reached)

Next day: Client can receive more emails ✅
```

**Scenario 2: Protection**
```
Malicious user tries to spam:
- Create 10 deliveries
- Only first 3 trigger emails
- Remaining 7 blocked
- System protected! ✅
```

### Monitoring Rate Limits

Your server logs show:
```
📧 Email sent (5/30 today) to client@example.com (2/3 for this client)
⚠️  Rate limit exceeded: email limit reached for client@example.com (3/3 emails today)
📊 Rate limiter reset for new day: 2026-01-05
```

---

## 🔧 SMTP2GO Dashboard Features

### 1. Sending Statistics
- Go to **Reports** → **Email Activity**
- See: Sent, Delivered, Bounced, Opened
- Monitor deliverability rate (aim for >95%)

### 2. Activity Log
- View every email sent
- Check delivery status
- See bounce reasons
- Debug issues

### 3. API Usage
- Track against 1,000/month limit
- Current usage displayed on dashboard
- Alerts when approaching limit

### 4. Sender Domains (Optional - Recommended for Production)
- Go to **Settings** → **Sender Domains**
- Click **Add Domain**
- Enter your domain (e.g., `ihbtransport.dk`)
- Add DNS records (SPF, DKIM, DMARC)
- Wait for verification (~1 hour)
- Now use `noreply@ihbtransport.dk` as FROM email

---

## 📈 Scaling Your Email Limits

### When to Increase Limits

**Signs you need more:**
- Rate limit warnings every day
- Legitimate emails being blocked
- Customer complaints about missing emails
- Growing business (10+ deliveries/day)

### How to Scale

**Option 1: Increase Daily Limit**
```go
// In utils/rate_limiter.go
dailyLimit: 50, // Instead of 30
```

**Option 2: Upgrade SMTP2GO Plan**
- 10,000 emails/month: $10/month
- 50,000 emails/month: $50/month
- 100,000 emails/month: $80/month

**Option 3: Smart Email Batching**
- Combine notifications (fewer emails)
- Summary emails (daily digest)
- SMS for critical updates

---

## 🎯 Best Practices for SMTP2GO

### 1. Warm Up Your Account
```
Week 1: Send 10-20 emails/day
Week 2: Send 30-50 emails/day
Week 3+: Full capacity
```

### 2. Monitor Bounce Rate
- Keep below 5%
- Remove invalid emails
- Verify email format before sending

### 3. Handle Errors Gracefully
```go
// Already implemented! ✅
- Failed emails logged to database
- Error messages captured
- Retry logic possible
```

### 4. Use Webhooks (Optional)
- Get real-time delivery notifications
- Track opens and clicks
- Update delivery status automatically

---

## 🚨 Troubleshooting SMTP2GO

### "Authentication Failed"
```bash
# Check credentials
echo $SMTP_USER      # Should be: ihb-transport-api
echo $SMTP_PASSWORD  # Should be your generated password

# No spaces in .env file
# Username is case-sensitive
```

### "Connection Timeout"
```bash
# Try different port
SMTP_PORT=587   # Instead of 2525
SMTP_PORT=8025  # Or this
```

### "Daily Limit Reached"
```bash
# Check current usage
curl http://localhost:8080/api/admin/email-stats \
  -H "Authorization: Bearer $TOKEN"

# Temporarily increase limit (in rate_limiter.go)
dailyLimit: 50  # Emergency increase
```

### "Emails Going to Spam"
1. Verify sender domain (see Dashboard → Sender Domains)
2. Add SPF record: `v=spf1 include:smtp2go.com ~all`
3. Test with mail-tester.com
4. Check content (avoid spam words)

---

## 💰 Cost Breakdown

### Free Tier (What you get)
- ✅ 1,000 emails/month
- ✅ SMTP + API access
- ✅ Email tracking
- ✅ Activity logs
- ✅ 99.9% uptime SLA

### When you grow (Paid plans)
| Plan | Emails | Cost | Cost per 1,000 |
|------|--------|------|----------------|
| Free | 1,000 | $0 | $0 |
| Basic | 10,000 | $10/mo | $1.00 |
| Plus | 50,000 | $50/mo | $1.00 |
| Pro | 100,000 | $80/mo | $0.80 |

**Your Startup Budget:**
- Month 1-6: Free tier ($0)
- Month 7+: Basic plan ($10/mo) if growing
- Very affordable! 🎉

---

## 📧 Email Templates Included

All 6 emails are professionally designed:

1. **Request Received** - Confirmation + tracking info
2. **Price Quote** - Clear pricing + accept/decline CTA
3. **Delivery Accepted** - Next steps + timeline
4. **Price Declined** - Thank you + admin notification
5. **Driver En Route** - Pickup confirmation + ETA
6. **Delivery Complete** - Success message + feedback request

Each email:
- ✅ Mobile responsive
- ✅ Professional HTML
- ✅ Clear branding
- ✅ Deliverability optimized

---

## 🔐 Security Best Practices

### 1. Protect Your Credentials
```bash
# Never commit .env to git
echo ".env" >> .gitignore

# Use environment variables in production
# Render: Dashboard → Environment Variables
# Railway: Dashboard → Variables
```

### 2. Monitor for Abuse
```bash
# Check logs daily for suspicious activity
grep "Rate limit" server.log
grep "Failed to send" server.log
```

### 3. Validate Email Addresses
```go
// Already implemented! ✅
// utils/input_validation.go validates emails
```

---

## ✅ Setup Checklist

- [ ] Create SMTP2GO account
- [ ] Create SMTP user credentials
- [ ] Add credentials to .env file
- [ ] Restart server
- [ ] Send test email
- [ ] Verify email received
- [ ] Check SMTP2GO dashboard
- [ ] Monitor rate limits
- [ ] (Optional) Verify sender domain
- [ ] Deploy to production
- [ ] Add SMTP env vars to hosting platform

---

## 🎉 You're Ready!

With SMTP2GO configured:
- ✅ Professional email delivery
- ✅ Rate limiting prevents abuse
- ✅ 1,000 free emails/month
- ✅ Excellent deliverability
- ✅ Production-ready setup

Start sending emails and growing your business! 🚀

---

**Need Help?**
- SMTP2GO Support: support@smtp2go.com
- SMTP2GO Docs: https://www.smtp2go.com/docs
- Test mail: mail-tester.com
