# Quick Start with Supabase

Get your IHB Transport API running in 10 minutes!

## 🎯 Quick Setup (Fastest Path)

### 1. Create Supabase Project (3 minutes)

```bash
# Go to: https://supabase.com
# Click: New Project
# Set: Name = "ihb-transport"
# Set: Database Password = (generate strong password)
# Set: Region = (closest to you)
# Wait: ~2 minutes for creation
```

### 2. Get Connection String (1 minute)

```bash
# In Supabase Dashboard:
# Go to: Settings → Database
# Find: "Connection string" section
# Copy: URI format
# It looks like: postgresql://postgres.xxxxx:PASSWORD@db.xxxxx.supabase.co:5432/postgres
```

### 3. Create .env File (1 minute)

```bash
# In your project root, create .env:
cat > .env << 'EOF'
DATABASE_URL=postgresql://postgres.xxxxx:yourpassword@db.xxxxx.supabase.co:5432/postgres
ADMIN_EMAIL=isaabaidoo@yahoo.com
ADMIN_PASSWORD=YourSecurePassword123!
JWT_SECRET=your-super-secret-jwt-key-minimum-32-chars-long
PORT=8080
GIN_MODE=debug
EOF
```

Replace `DATABASE_URL` with your actual Supabase connection string!

### 4. Run Locally (2 minutes)

```bash
# Install dependencies
go mod download

# Run migrations (create tables)
go run cmd/migrate/main.go

# Seed admin user
go run cmd/seed/main.go

# Start server
go run cmd/server/main.go
```

You should see:
```
📡 Using DATABASE_URL connection string
✅ Database connected successfully
✅ Migrations completed successfully
✅ Server running on port 8080
```

### 5. Test API (1 minute)

```bash
# Test ping
curl http://localhost:8080/ping

# Login as admin
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"isaabaidoo@yahoo.com","password":"YourSecurePassword123!"}'

# Create a delivery request
curl -X POST http://localhost:8080/api/public/deliveries \
  -H "Content-Type: application/json" \
  -d '{
    "client_name": "Test User",
    "client_email": "test@example.com",
    "pickup_address": "123 Start St",
    "dropoff_address": "456 End Ave",
    "item_description": "Test Package",
    "weight": 5.0
  }'
```

## 🚀 Deploy to Production (5 minutes)

### Option 1: Railway (Easiest - Auto-detects Dockerfile)

```bash
# 1. Push to GitHub
git add .
git commit -m "Add Supabase integration"
git push origin main

# 2. Go to: https://railway.app
# 3. Click: "New Project" → "Deploy from GitHub repo"
# 4. Select: Your repository
# 5. Add environment variables:
#    DATABASE_URL = <your-supabase-connection-string>
#    ADMIN_EMAIL = isaabaidoo@yahoo.com
#    ADMIN_PASSWORD = <your-password>
#    JWT_SECRET = <your-jwt-secret>
#    PORT = 8080
#    GIN_MODE = release

# 6. Railway auto-deploys! ✅
# 7. After deployment, open Shell tab and run:
./migrate
./seed
```

### Option 2: Render (Free Tier Available)

```bash
# 1. Push to GitHub (same as above)

# 2. Go to: https://dashboard.render.com
# 3. Click: "New +" → "Web Service"
# 4. Connect: Your GitHub repo
# 5. Set:
#    - Name: ihb-transport-api
#    - Runtime: Docker
#    - Instance Type: Free
# 6. Add same environment variables as Railway
# 7. Click: "Create Web Service"
# 8. After deployment, open Shell and run:
./migrate
./seed
```

## 📊 Verify Deployment

```bash
# Replace YOUR_DEPLOYED_URL with your actual URL
export API_URL="https://ihb-transport-api.up.railway.app"

# Test ping
curl $API_URL/ping

# Login
curl -X POST $API_URL/login \
  -H "Content-Type: application/json" \
  -d '{"email":"isaabaidoo@yahoo.com","password":"YourPassword"}'
```

## 🎨 Next Steps

### 1. Check Supabase Dashboard
- Go to Table Editor
- See your tables: `admins`, `delivery_requests`, `status_logs`, `email_logs`
- View data in real-time

### 2. Test All Endpoints
Use the test scripts:
```bash
# Update TOKEN in test files with your login token
./test_api.sh
./test_decline.sh
```

### 3. Build Frontend
You have 13 API endpoints ready:
- Public: Create delivery, track by ID/email, accept/decline
- Admin: Login, view all, set price
- Driver: Mark as picked up, mark as delivered

See `DEPLOYMENT.md` for complete endpoint documentation.

### 4. Configure Real Emails (Optional)
Currently emails log to console. To send real emails:
- Choose provider: SendGrid, Mailgun, AWS SES
- Update `internal/services/email_service.go`
- Add SMTP credentials to .env

## 🔧 Troubleshooting

### "Missing required environment variable"
- Check `.env` file exists
- Verify `DATABASE_URL` is set
- Make sure no extra spaces in .env

### "Failed to connect to database"
- Verify Supabase connection string is correct
- Check password doesn't have special characters that need escaping
- Test connection in Supabase dashboard

### "Port already in use"
- Stop other processes: `lsof -ti:8080 | xargs kill -9`
- Or change PORT in .env to 8081

### Migrations not running
- Check database connection first
- Run manually: `go run cmd/migrate/main.go`
- Check Supabase logs in dashboard

## 📚 Documentation

- `SUPABASE_SETUP.md` - Detailed Supabase guide with features
- `DEPLOYMENT.md` - Complete deployment guide
- `README.md` - Project overview and API endpoints

## 💡 Tips

**Local Development:**
- Use `GIN_MODE=debug` for detailed logs
- Use Supabase connection string even locally
- Check Supabase dashboard to see database changes in real-time

**Production:**
- Use `GIN_MODE=release` for better performance
- Generate strong JWT_SECRET: `openssl rand -base64 32`
- Enable Supabase Row Level Security (RLS) for extra security

**Database:**
- Supabase free tier: 500MB (handles 100,000+ deliveries)
- Auto-backups daily
- No need to manage PostgreSQL yourself

---

**You're all set! 🎉**

Your mini-trans style delivery platform is ready to go. Start building your frontend using the API endpoints!
