# IHB Transport - Deployment Guide

## 🚀 Recommended: Supabase + Render/Railway

This guide covers deploying with Supabase as the database (free 500MB) and Render or Railway for hosting (both have generous free tiers).

### Prerequisites
1. GitHub account with this repository pushed
2. Supabase account (https://supabase.com) - **Recommended**
3. Hosting account: Render (https://render.com) or Railway (https://railway.app)

## Database Setup

### Option A: Supabase (Recommended - Free 500MB + Auto Backups)

1. Go to [Supabase](https://supabase.com) and create a new project
2. Choose a **Database Password** (save this!)
3. Select region closest to your users (e.g., `eu-central-1`)
4. Wait 2-3 minutes for project creation
5. Get connection string:
   - Go to Settings → Database
   - Find "Connection string" section
   - Copy **URI** format: `postgresql://postgres.xxxxx:[PASSWORD]@db.xxxxx.supabase.co:5432/postgres`
   - Replace `[PASSWORD]` with your database password

**See `SUPABASE_SETUP.md` for detailed Supabase configuration and features.**

### Option B: Render PostgreSQL (Alternative)

1. Go to Render Dashboard
2. Click "New +" → "PostgreSQL"
3. Fill in:
   - Name: `ihb-transport-db`
   - Database: `logisticsDB`
   - Region: Choose closest to your users
   - Plan: Free tier
4. Click "Create Database"
5. Copy the **Internal Database URL**

## Hosting Setup

### Deploy to Render with Supabase

### Step 1: Create Web Service on Render

1. Go to [Render Dashboard](https://dashboard.render.com)
2. Click "New +" → "Web Service"
3. Connect your GitHub repository
4. Fill in:
   - Name: `ihb-transport-api`
   - Region: Choose closest (same as Supabase if possible)
   - Branch: `main`
   - Root Directory: (leave empty)
   - Runtime: `Docker`
   - Instance Type: Free tier

### Step 2: Set Environment Variables

In the Environment section, add:

```
DATABASE_URL=postgresql://postgres.xxxxx:yourpassword@db.xxxxx.supabase.co:5432/postgres
ADMIN_EMAIL=isaabaidoo@yahoo.com
ADMIN_PASSWORD=your-secure-password
JWT_SECRET=your-super-secret-jwt-key-here-change-this
PORT=8080
GIN_MODE=release
```

**Important**: Use your actual Supabase connection string for `DATABASE_URL`

### Step 3: Deploy

1. Click "Create Web Service"
2. Render will automatically:
   - Build your Docker image
   - Connect to Supabase
   - Start the server

### Step 4: Run Migrations (First Time)

After first deployment, go to Shell tab and run:
```bash
./migrate
./seed
```

Or trigger manually via API after deployment.
```bash
go run cmd/migrate/main.go && go run cmd/seed/main.go
```

### Step 6: Seed Admin User

Run in Shell:
```bash
go run cmd/seed/main.go
```

---

## 🔗 Your API URLs

Once deployed, your API will be available at:
```
https://ihb-transport-api.onrender.com
```

Test it:
```bash
curl https://ihb-transport-api.onrender.com/ping
```

---

## 📝 Important Notes

1. **Free Tier Limitations:**
   - Service spins down after 15 minutes of inactivity
   - First request after spin-down takes 30-60 seconds
   - Database has 90-day expiration (backup your data!)

2. **Production Checklist:**
   - ✅ Change JWT_SECRET to a strong random string
   - ✅ Use strong admin password
   - ✅ Enable CORS for your frontend domain
   - ✅ Set up proper error logging
   - ✅ Regular database backups

3. **Health Check:**
   - Render automatically pings `/ping` endpoint
   - Make sure it returns 200 OK

---

## 🐳 Local Docker Testing

Build and test locally before deploying:

```bash
# Build image
docker build -t ihb-transport .

# Run with environment variables
docker run -p 8080:8080 \
  -e DATABASE_URL="postgresql://user:pass@host:5432/db" \
  -e ADMIN_EMAIL="admin@example.com" \
  -e ADMIN_PASSWORD="password" \
  -e JWT_SECRET="secret" \
  ihb-transport
```

Test:
```bash
curl http://localhost:8080/ping
```

---

## 🔄 Update Deployment

Push to GitHub main branch:
```bash
git add .
git commit -m "Update deployment"
git push origin main
```

Render will automatically rebuild and deploy! 🚀
