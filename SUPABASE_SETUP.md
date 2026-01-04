# Supabase Setup Guide

This guide will help you set up your IHB Transport backend with Supabase.

## Why Supabase?

- ✅ Free PostgreSQL database (500MB included)
- ✅ Built-in authentication (optional - we use our custom JWT)
- ✅ Real-time subscriptions (future feature)
- ✅ Automatic backups
- ✅ Easy connection string management
- ✅ Dashboard for database management

## Step 1: Create Supabase Project

1. Go to [https://supabase.com](https://supabase.com)
2. Sign up or log in
3. Click "New Project"
4. Fill in:
   - **Name**: `ihb-transport` (or your preferred name)
   - **Database Password**: Generate a strong password (save this!)
   - **Region**: Choose closest to your users (e.g., `eu-central-1` for Europe)
   - **Pricing Plan**: Free tier is perfect to start

5. Wait 2-3 minutes for project to be created

## Step 2: Get Database Connection String

1. In your Supabase project dashboard, click **Settings** (gear icon)
2. Click **Database** in the left sidebar
3. Scroll to **Connection string** section
4. Select **URI** tab
5. Copy the connection string (looks like this):
   ```
   postgresql://postgres.xxxxx:[YOUR-PASSWORD]@db.xxxxx.supabase.co:5432/postgres
   ```
6. Replace `[YOUR-PASSWORD]` with the database password you created in Step 1

## Step 3: Update Your .env File

Create or update `.env` file in your project root:

```bash
# Supabase Database Connection
DATABASE_URL=postgresql://postgres.xxxxx:yourpassword@db.xxxxx.supabase.co:5432/postgres

# Admin Configuration
ADMIN_EMAIL=isaabaidoo@yahoo.com
ADMIN_PASSWORD=your_secure_admin_password

# JWT Secret (generate a random string)
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Server Configuration
PORT=8080
GIN_MODE=release
```

**Important**: The `DATABASE_URL` will automatically enable SSL. Supabase requires SSL connections.

## Step 4: Run Migrations

With Supabase connected, run migrations to create your tables:

```bash
# Create tables
go run cmd/migrate/main.go

# Seed admin user
go run cmd/seed/main.go
```

You should see:
```
📡 Using Supabase connection string
✅ Database connected successfully
✅ Migrations completed successfully
```

## Step 5: Verify in Supabase Dashboard

1. Go to your Supabase project
2. Click **Table Editor** in the left sidebar
3. You should see your tables:
   - `admins`
   - `delivery_requests`
   - `status_logs`
   - `email_logs`

## Step 6: Run Your Server

```bash
go run cmd/server/main.go
```

Test with:
```bash
curl http://localhost:8080/ping
```

## Deployment Options

### Option A: Deploy to Render with Supabase

1. Create `render.yaml` (auto-deployment config):

```yaml
services:
  - type: web
    name: ihb-transport-api
    env: docker
    plan: free
    envVars:
      - key: DATABASE_URL
        sync: false  # You'll set this manually
      - key: ADMIN_EMAIL
        value: isaabaidoo@yahoo.com
      - key: ADMIN_PASSWORD
        sync: false
      - key: JWT_SECRET
        generateValue: true
      - key: PORT
        value: 8080
      - key: GIN_MODE
        value: release
```

2. Push to GitHub
3. Connect to Render
4. Set `DATABASE_URL` environment variable with your Supabase connection string
5. Deploy!

### Option B: Deploy to Railway

Railway has native Supabase integration:

1. Push to GitHub
2. Go to [railway.app](https://railway.app)
3. Click "New Project" → "Deploy from GitHub repo"
4. Select your repository
5. Add environment variables:
   - `DATABASE_URL`: Your Supabase connection string
   - `ADMIN_EMAIL`, `ADMIN_PASSWORD`, `JWT_SECRET`, `PORT`, `GIN_MODE`
6. Railway will auto-detect Dockerfile and deploy

### Option C: Deploy to Fly.io

1. Install flyctl: `curl -L https://fly.io/install.sh | sh`
2. Login: `fly auth login`
3. Create app: `fly launch --no-deploy`
4. Set secrets:
   ```bash
   fly secrets set DATABASE_URL="your-supabase-connection-string"
   fly secrets set ADMIN_EMAIL="isaabaidoo@yahoo.com"
   fly secrets set ADMIN_PASSWORD="your-password"
   fly secrets set JWT_SECRET="your-jwt-secret"
   ```
5. Deploy: `fly deploy`

## Local Development vs Production

### Local Development
Use individual environment variables in `.env`:
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=localpassword
DB_NAME=ihb_transport
DB_SSLMODE=disable
```

The code automatically falls back to these if `DATABASE_URL` is not set.

### Production
Always use `DATABASE_URL` with Supabase connection string.

## Supabase Features You Can Add Later

### 1. Row Level Security (RLS)
Add policies to restrict data access at the database level:
- Clients can only view their own deliveries
- Drivers can only update assigned deliveries

### 2. Real-time Subscriptions
Add WebSocket support for live delivery tracking:
```javascript
// Frontend example
const subscription = supabase
  .from('delivery_requests')
  .on('UPDATE', payload => {
    console.log('Delivery updated!', payload)
  })
  .subscribe()
```

### 3. Storage for Photos
Use Supabase Storage for delivery photos:
- Photo of item before pickup
- Proof of delivery photo

### 4. Edge Functions (Optional)
Deploy serverless functions for background tasks:
- Send scheduled reminder emails
- Generate delivery reports
- Calculate pricing based on distance

## Monitoring and Maintenance

### View Logs in Supabase
1. Go to **Logs** in left sidebar
2. Select **Postgres Logs** to see database queries
3. Select **API Logs** if using Supabase Auth/Storage

### Backup Database
Supabase automatically backs up your database daily. To manually backup:
1. Go to **Database** → **Backups**
2. Click "Start a backup"

### Check Database Size
1. Go to **Settings** → **Database**
2. View usage under **Disk** section
3. Free tier includes 500MB (plenty for thousands of deliveries)

## Troubleshooting

### Connection Timeout
- Check if `DATABASE_URL` is correct
- Ensure password doesn't contain special characters that need URL encoding
- Verify your IP isn't blocked (Supabase allows all IPs by default)

### SSL Required Error
- Supabase requires SSL. Your connection string should work automatically.
- If issues persist, add `?sslmode=require` to connection string

### Migration Errors
- Check Supabase dashboard **Logs** section
- Ensure you're using latest GORM version
- Verify table names don't conflict with Supabase internal tables

## Cost Estimation

**Free Tier Limits:**
- Database: 500MB
- Bandwidth: 2GB
- API Requests: 50,000/month
- Storage: 1GB (if using Supabase Storage)

**For mini-trans style single driver:**
- ~1000 deliveries = ~5MB database space
- You can handle **100,000+ deliveries** on free tier!

## Next Steps

1. ✅ Connect to Supabase
2. ✅ Run migrations
3. ✅ Test locally
4. 🚀 Deploy to hosting platform
5. 🎨 Build frontend
6. 📧 Configure SMTP for real emails
7. 📱 Add mobile app (optional)

## Support

- **Supabase Docs**: https://supabase.com/docs
- **Supabase Discord**: https://discord.supabase.com
- **Your API Docs**: See `DEPLOYMENT.md` for endpoint documentation

---

**Happy building! 🚀**
