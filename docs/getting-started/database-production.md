# Database Production Migration Guide

Panduan aman untuk menjalankan database migration di server production. **Sangat penting dibaca** sebelum deploy ke production!

## ⚠️ PERINGATAN PENTING

Migration di production bisa **merusak data** kalau salah. Selalu:
- ✅ **Test di staging dulu**
- ✅ **Backup database sebelum migration**
- ✅ **Jalankan saat traffic rendah**
- ✅ **Siapkan rollback plan**

## 🏭 Development vs Production

| Aspek | Development | Production |
|-------|-------------|------------|
| **Database** | Local Docker | Server Database |
| **Confirmation** | Tidak ada | **2x confirmation** |
| **SSL** | Disabled | **Required** |
| **Password** | Simple | **Strong password** |
| **Backup** | Tidak perlu | **Wajib backup** |
| **Downtime** | Tidak masalah | **Harus minimal** |
| **Rollback** | Gampang reset | **Harus hati-hati** |

## 🔐 Setup Environment Production

### 1. Environment Variables

**Set environment variables ini sebelum migration:**

```bash
# Production environment flag
export APP_SERVER_ENV=production

# Database connection
export APP_DATABASE_HOST=your-prod-db.example.com
export APP_DATABASE_PORT=5432
export APP_DATABASE_USER=app_user
export APP_DATABASE_PASSWORD=your-very-secure-password
export APP_DATABASE_NAME=your_app_production
export APP_DATABASE_SSLMODE=require

# Migration path (opsional, default sudah benar)
export MIGRATION_PATH=./internal/data/migrations
```

### 2. Verifikasi Connection

Test connection dulu sebelum migration:

```bash
# Test connection ke production database
psql \"postgres://$APP_DATABASE_USER:$APP_DATABASE_PASSWORD@$APP_DATABASE_HOST:$APP_DATABASE_PORT/$APP_DATABASE_NAME?sslmode=require\"
```

**Kalau berhasil connect:**
```
your_app_production=>
```

Ketik `\\q` untuk keluar.

**Kalau gagal:**
- Cek credentials
- Cek firewall database server
- Cek SSL certificate

## 📋 Pre-Migration Checklist

### ✅ Before Migration (Wajib!)

- [ ] **Backup database production**
- [ ] **Test migration di staging environment**
- [ ] **Set semua environment variables**
- [ ] **Verify aplikasi bisa handle schema changes**
- [ ] **Schedule maintenance window (kalau perlu)**
- [ ] **Notify tim tentang maintenance**
- [ ] **Prepare rollback plan**
- [ ] **Test database connection**

### 💾 Backup Database

```bash
# Backup production database
pg_dump \"postgres://$APP_DATABASE_USER:$APP_DATABASE_PASSWORD@$APP_DATABASE_HOST:$APP_DATABASE_PORT/$APP_DATABASE_NAME?sslmode=require\" > backup_$(date +%Y%m%d_%H%M%S).sql

# Compress backup (opsional)
gzip backup_*.sql
```

**Simpan backup di tempat aman!** (cloud storage, server terpisah)

## 🚀 Menjalankan Production Migration

### 1. Check Current Status

Selalu cek status migration sebelum mulai:

```bash
# Using production environment
make migrate-prod ARGS=\"status\"

# Or using script directly
./scripts/migrate.sh status
```

Output example:
```
Migration status for environment: production
Database: postgres://app_user:***@prod-db.example.com:5432/app_prod
Version    Path
---------- ------------------------------------------
5          internal/data/migrations/000005_add_products_table.up.sql
```

### 2. Run Migration (Dengan Safety Checks)

```bash
# Method 1: Using make (recommended)
make migrate-prod ARGS=\"up\"

# Method 2: Using script directly
APP_SERVER_ENV=production ./scripts/migrate.sh up
```

**Apa yang terjadi:**

1. **Environment Check**: Script detect production environment
2. **First Confirmation**: 
   ```
   [WARNING] You are about to run migration on PRODUCTION environment!
   Database: postgres://app_user:***@prod-db.example.com:5432/app_prod
   Are you sure? (yes/no):
   ```
   Ketik `yes` dan Enter.

3. **Migration Execution**: Script jalankan migration
4. **Status Update**: Confirmation migration berhasil

### 3. Verify Migration Success

```bash
# Check migration status
make migrate-prod ARGS=\"status\"

# Check application health
curl https://your-production-app.com/health

# Check database structure via psql
psql \"postgres://$APP_DATABASE_USER:$APP_DATABASE_PASSWORD@$APP_DATABASE_HOST:$APP_DATABASE_PORT/$APP_DATABASE_NAME?sslmode=require\" -c \"\\dt\"
```

## 🎯 Step-by-Step Migration

Kalau ada banyak migration dan mau jalankan satu per satu:

### 1. Run One Migration
```bash
# Run exactly 1 migration
make migrate ARGS=\"up 1\"

# Check what happened
make migrate-prod ARGS=\"status\"

# Test application
curl https://your-app.com/health
```

### 2. Repeat for Each Migration
```bash
# Run next migration
make migrate ARGS=\"up 1\"

# Always check after each step
make migrate-prod ARGS=\"status\"
```

**Kenapa step-by-step?**
- Bisa stop kalau ada masalah
- Easier troubleshooting
- Smaller risk per step

## 🔄 Rollback Production Migration

**⚠️ HATI-HATI: Rollback bisa hapus data!**

### 1. Check What to Rollback
```bash
# Lihat migration history
make migrate-prod ARGS=\"status\"
```

### 2. Rollback One Migration
```bash
# Rollback 1 migration terakhir
make migrate-prod ARGS=\"down\"
```

**Double confirmation akan muncul:**
```
[WARNING] Running migrations DOWN for environment: production
[WARNING] This will rollback 1 migration(s)
Confirm rollback in production (yes/no): yes
```

### 3. Force Migration Version (Emergency)
```bash
# Kalau migration stuck, force ke versi tertentu
make migrate-prod ARGS=\"force 4\"
```

**⚠️ BAHAYA:** Force bisa bikin data inconsistent. Hanya untuk emergency!

## 📊 Production Migration Strategies

### Strategy 1: Blue-Green Deployment
1. Setup database schema di green environment
2. Run migration di green database
3. Switch traffic ke green
4. Keep blue as backup

### Strategy 2: Rolling Migration
1. Run backward-compatible migration
2. Deploy new application version
3. Run cleanup migration (remove old columns)

### Strategy 3: Maintenance Window
1. Schedule downtime (misal: 2 AM)
2. Stop application
3. Run migration
4. Start application
5. Monitor

## 🔍 Monitoring Production Migration

### During Migration
```bash
# Monitor database connections
# (di server database)
SELECT * FROM pg_stat_activity WHERE datname = 'your_app_production';

# Monitor application logs
tail -f /var/log/your-app/app.log

# Monitor system resources
htop
```

### After Migration
- **Response time**: Apakah lebih lambat?
- **Error rate**: Ada error baru?
- **Database size**: Sesuai ekspektasi?
- **Query performance**: Index berfungsi?

## 🚨 Emergency Procedures

### Migration Stuck/Failed
```bash
# 1. Check migration status
make migrate-prod ARGS=\"status\"

# 2. Check database locks
psql \"postgres://...\" -c \"SELECT * FROM pg_locks WHERE granted = false;\"

# 3. Check long-running queries
psql \"postgres://...\" -c \"SELECT pid, query, query_start FROM pg_stat_activity WHERE state = 'active';\"

# 4. Force migration version (last resort)
make migrate-prod ARGS=\"force <last_good_version>\"
```

### Application Error After Migration
```bash
# 1. Rollback migration
make migrate-prod ARGS=\"down\"

# 2. Or restore from backup
psql \"postgres://...\" < backup_20240129_120000.sql

# 3. Restart application
systemctl restart your-app
```

### Database Connection Issues
```bash
# Check connection
psql \"postgres://$APP_DATABASE_USER:$APP_DATABASE_PASSWORD@$APP_DATABASE_HOST:$APP_DATABASE_PORT/$APP_DATABASE_NAME?sslmode=require\"

# Check environment variables
env | grep APP_DATABASE

# Check SSL certificate
openssl s_client -connect $APP_DATABASE_HOST:$APP_DATABASE_PORT
```

## 📋 Post-Migration Checklist

### ✅ Immediate (0-15 minutes)
- [ ] Migration status shows correct version
- [ ] Application health check passes
- [ ] No error spikes in logs
- [ ] Database connections stable

### ✅ Short-term (15-60 minutes)
- [ ] Key user flows working
- [ ] API response times normal
- [ ] No customer complaints
- [ ] Database performance stable

### ✅ Long-term (1-24 hours)
- [ ] All features working correctly
- [ ] Performance metrics normal
- [ ] No data corruption
- [ ] Monitor error rates

## 🎓 Production Migration Best Practices

### Planning
- **Test everything in staging first**
- **Use same database version as production**
- **Test with production-like data volume**
- **Plan rollback strategy**

### Execution
- **Run during low-traffic hours**
- **Monitor application during migration**
- **Keep database backup recent**
- **Have team on standby**

### Communication
- **Notify stakeholders beforehand**
- **Status updates during migration**
- **Post-migration report**

## 📚 Advanced Topics

### Large Table Migrations
Untuk tabel dengan jutaan record:
```sql
-- Use CONCURRENTLY for index creation (tidak lock table)
CREATE INDEX CONCURRENTLY idx_large_table_field ON large_table(field);

-- Use small batches for data updates
UPDATE large_table SET new_field = 'value' WHERE id BETWEEN 1 AND 10000;
-- Repeat for other ranges
```

### Zero-Downtime Migrations
1. **Add column with default value**
2. **Deploy app yang support both old/new schema**
3. **Migrate data in background**
4. **Deploy app yang cuma pakai new schema**
5. **Remove old column**

### Database Sharding Migration
- Plan shard key carefully
- Migrate one shard at a time
- Test cross-shard queries

## 🔗 Next Steps

- **Butuh troubleshooting?** → [Troubleshooting Guide](troubleshooting.md)
- **Mau setup CI/CD migration?** → [Advanced Migration Guide](../database/migrations/README.md)
- **Performance issues?** → [Database Performance Guide](../database/performance.md)

---

**⚠️ Remember:**
- Production migration = high stakes
- Always have backup & rollback plan
- Test everything in staging first
- Monitor closely after migration