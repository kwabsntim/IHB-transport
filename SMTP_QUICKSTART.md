# SMTP Quick Start

## 🚀 Fastest Setup (Gmail - 5 minutes)

### 1. Get Gmail App Password
```
1. Visit: https://myaccount.google.com/apppasswords
2. Select: Mail → Other (IHB Transport)
3. Copy the 16-character password
```

### 2. Add to .env
```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=youremail@gmail.com
SMTP_PASSWORD=your-16-char-app-password
SMTP_FROM_EMAIL=youremail@gmail.com
SMTP_FROM_NAME=IHB Transport
```

### 3. Restart & Test
```bash
# Restart server
go run cmd/server/main.go

# In another terminal, test emails
./test_smtp.sh
```

---

## 📧 Email Types

Your app sends 6 types of emails:

1. **REQUEST_RECEIVED** - When delivery is created
2. **PRICE_SENT** - When admin sets price
3. **ACCEPTED** - When client accepts price
4. **PRICE_DECLINED** - When client declines (+ admin notification)
5. **DRIVER_ON_WAY** - When driver picks up package
6. **DELIVERED** - When delivery is complete

All emails are:
- ✅ HTML formatted
- ✅ Mobile responsive
- ✅ Professional design
- ✅ Logged to database

---

## 🎯 Testing Without SMTP

If you don't configure SMTP, emails will log to console:

```
⚠️  SMTP not configured - emails will be logged to console only
📧 EMAIL (Console): Delivery Request Received #abc123
   To: client@example.com
   (SMTP not configured - email not sent)
```

This is perfect for development!

---

## 🔧 Troubleshooting

**"Failed to send email"**
- Check username and password in .env
- For Gmail, use App Password (not regular password)
- Restart server after changing .env

**Emails go to spam**
- Use custom domain instead of gmail.com
- Set up SPF/DKIM records
- See SMTP_SETUP.md for details

**"Authentication failed"**
- Verify credentials are correct
- No spaces in .env values
- For Gmail, enable 2FA first

---

## 📊 Free Tier Limits

| Provider | Free Emails |
|----------|-------------|
| Gmail | 500/day |
| SMTP2GO | 1,000/month |
| SendGrid | 100/day |
| Mailgun | 5,000/month (3 months) |

---

## 🚀 Production Deploy

Add these environment variables to your hosting platform:

**Render:**
```
Dashboard → Environment → Add Environment Variable
```

**Railway:**
```
Dashboard → Variables → New Variable
```

**Fly.io:**
```bash
fly secrets set SMTP_HOST=smtp.gmail.com
fly secrets set SMTP_PORT=587
fly secrets set SMTP_USER=youremail@gmail.com
fly secrets set SMTP_PASSWORD=your-password
```

---

## 📚 Full Documentation

- **SMTP_SETUP.md** - Complete setup guide for all providers
- **test_smtp.sh** - Test all 6 email types
- **.env.example** - Configuration template

Need help? See SMTP_SETUP.md for detailed instructions.
