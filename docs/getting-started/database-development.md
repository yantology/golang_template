# Database Development Guide

Panduan praktis menggunakan database PostgreSQL untuk development lokal. Cocok untuk pemula yang baru kenal database.

## 🎯 Apa itu Database Migration?

**Migration** adalah cara untuk mengubah struktur database secara bertahap dan terkontrol.

**Analogi sederhana:**
- Database = Lemari
- Table = Laci dalam lemari  
- Migration = Instruksi untuk tambah/ubah/hapus laci
- Migration file = Catatan instruksi yang bisa diulang

**Kenapa perlu migration?**
- Tim bisa sync struktur database yang sama
- Bisa rollback kalau ada masalah
- Track perubahan database seperti git untuk code

## 🐘 Setup Database Lokal

### 1. Start Database
```bash
make db-up
```

**Apa yang terjadi?**
- Docker download PostgreSQL image (kalau belum ada)
- Jalankan PostgreSQL di port 5432
- Jalankan Adminer (database UI) di port 8081
- Buat database kosong bernama `golang_template_dev`

**Cek berhasil atau tidak:**
```bash
# Cek container jalan
docker-compose ps

# Should show:
# postgres container running on port 5432
# adminer container running on port 8081
```

### 2. Akses Database UI
Buka browser ke http://localhost:8081

**Login credentials:**
- **System:** PostgreSQL
- **Server:** postgres
- **Username:** postgres
- **Password:** dev_password  
- **Database:** golang_template_dev

**💡 Tips:** Bookmark URL ini, akan sering dipakai untuk lihat data.

### 3. Test Connection
```bash
# Connect via terminal
make db-connect

# Kalau berhasil, akan masuk ke psql prompt:
golang_template_dev=#
```

Ketik `\\q` untuk keluar dari psql.

## 📝 Membuat Migration Pertama

### 1. Buat Migration File
```bash
make migrate-create NAME=create_users_table
```

**Apa yang terjadi?**
- Dibuat 2 file di `internal/data/migrations/`:
  - `000001_create_users_table.up.sql` (untuk maju)
  - `000001_create_users_table.down.sql` (untuk mundur)

### 2. Edit Migration File

**File: `000001_create_users_table.up.sql`**
```sql
-- Buat tabel users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Index untuk pencarian cepat berdasarkan email
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Index untuk filter user aktif
CREATE INDEX IF NOT EXISTS idx_users_active ON users(is_active);
```

**File: `000001_create_users_table.down.sql`**
```sql
-- Hapus index dulu
DROP INDEX IF EXISTS idx_users_active;
DROP INDEX IF EXISTS idx_users_email;

-- Baru hapus tabel
DROP TABLE IF EXISTS users;
```

**💡 Penjelasan:**
- **UP file:** Instruksi untuk maju (bikin tabel baru)
- **DOWN file:** Instruksi untuk mundur (hapus tabel)
- **UUID:** ID unik universal, lebih aman dari auto-increment
- **INDEX:** Biar query email dan status lebih cepat

### 3. Jalankan Migration
```bash
make migrate-up
```

**Apa yang terjadi?**
- Database baca file `.up.sql`
- Jalankan perintah SQL di dalamnya
- Catat di tabel `schema_migrations` bahwa migration sudah jalan

**Cek berhasil:**
- Buka Adminer (http://localhost:8081)
- Pilih database `golang_template_dev`
- Harusnya muncul tabel `users` dan `schema_migrations`

### 4. Cek Status Migration
```bash
make migrate-status
```

Output yang diharapkan:
```
Version    Path
---------- ------------------------------------------
1          internal/data/migrations/000001_create_users_table.up.sql
```

## 🔄 Development Migration Workflow

### Workflow Harian
```bash
# 1. Start database (kalau belum)
make db-up

# 2. Cek migration status
make migrate-status

# 3. Jalankan pending migrations
make migrate-up

# 4. Mulai coding...
```

### Buat Migration Baru
```bash
# Buat migration untuk tabel products
make migrate-create NAME=create_products_table

# Edit file yang dibuat di internal/data/migrations/
# Lalu jalankan
make migrate-up
```

### Rollback Migration
```bash
# Rollback 1 migration terakhir
make migrate-down

# Cek status setelah rollback
make migrate-status
```

**⚠️ Hati-hati:** Rollback akan **hapus data** di tabel yang di-drop!

## 🔍 Melihat dan Mengelola Data

### Via Adminer (GUI)
1. Buka http://localhost:8081
2. Login dengan credentials di atas
3. Pilih tabel untuk lihat data
4. Bisa insert, update, delete data via interface

### Via Terminal (CLI)
```bash
# Connect ke database
make db-connect

# Lihat semua tabel
\\dt

# Lihat struktur tabel users
\\d users

# Query data
SELECT * FROM users;

# Insert data test
INSERT INTO users (email, password_hash, first_name, last_name) 
VALUES ('test@example.com', 'hashed_password', 'John', 'Doe');

# Keluar
\\q
```

### Via Code (Go)
Kalau mau insert data dari aplikasi Go, tinggal jalankan:
```bash
make dev
```

Lalu test endpoint yang buat/baca data users.

## 🛠️ Commands Reference

### Database Management
```bash
make db-up          # Start database
make db-down        # Stop database
make db-logs        # Lihat database logs
make db-connect     # Connect via psql
make db-clean       # Reset database (hapus semua data!)
```

### Migration Commands
```bash
make migrate-create NAME=nama_migration  # Buat migration baru
make migrate-up                          # Jalankan semua pending migrations
make migrate-down                        # Rollback 1 migration
make migrate-status                      # Lihat migration status
make migrate-force VERSION=1             # Force ke versi tertentu (bahaya!)
```

### Advanced Migration (Helper Script)
```bash
./scripts/migrate.sh status              # Status migration
./scripts/migrate.sh up                  # Jalankan migrations
./scripts/migrate.sh up 2                # Jalankan 2 migrations
./scripts/migrate.sh down               # Rollback 1 migration
./scripts/migrate.sh create add_table   # Buat migration
```

## 📚 Migration Best Practices

### ✅ DO (Lakukan)
- **Selalu buat UP dan DOWN migration**
- **Test migration sebelum commit**
- **Gunakan `IF EXISTS` / `IF NOT EXISTS`**
- **Buat index untuk kolom yang sering di-query**
- **Gunakan UUID untuk primary key**
- **Tambahkan timestamps (created_at, updated_at)**

### ❌ DON'T (Jangan)
- **Jangan edit migration yang sudah jalan di production**
- **Jangan lupa buat DOWN migration**
- **Jangan migration yang terlalu besar sekaligus**
- **Jangan hapus migration file yang sudah commit**

### Migration Example Template
```sql
-- UP migration template
CREATE TABLE IF NOT EXISTS table_name (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_table_name_field ON table_name(field);

-- DOWN migration template  
DROP INDEX IF EXISTS idx_table_name_field;
DROP TABLE IF EXISTS table_name;
```

## 🚨 Troubleshooting

### Database tidak mau start
```bash
# Cek apakah port 5432 sudah dipakai
lsof -i :5432

# Kalau ada, stop container lain atau ganti port
make db-down && make db-up
```

### Migration error
```bash
# Lihat error detail
make migrate-status

# Reset migration (development only!)
make db-clean && make db-up && make migrate-up
```

### Connection refused
```bash
# Cek container status
docker-compose ps

# Restart database
make db-down
make db-up
```

### Adminer tidak bisa login
- Pastikan credentials benar: postgres/postgres/dev_password
- Cek container postgres jalan: `docker-compose ps`
- Restart adminer: `make db-down && make db-up`

## 🔗 Next Steps

- **Mau deploy ke production?** → [Database Production Guide](database-production.md)
- **Mau tau testing database?** → [Testing Guide](testing-guide.md)
- **Ada masalah lain?** → [Troubleshooting Guide](troubleshooting.md)

---

**💡 Pro Tips:**
- Backup database sebelum eksperimen: `pg_dump > backup.sql`
- Gunakan Adminer untuk explore data structure
- Buat migration kecil-kecil, jangan langsung besar
- Selalu test rollback migration di development