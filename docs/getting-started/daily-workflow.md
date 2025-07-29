# Daily Development Workflow

Panduan sederhana untuk workflow harian development. Ikuti langkah-langkah ini setiap kali mulai coding.

## 🌅 Memulai Hari Development

### 1. Nyalakan Database
Sebelum coding, database harus jalan dulu:

```bash
make db-up
```

**Apa yang terjadi?**
- Docker akan menjalankan PostgreSQL di port 5432
- Adminer (database UI) akan jalan di port 8081
- Data database tersimpan supaya tidak hilang saat restart

**Cek apakah berhasil:**
- Buka http://localhost:8081 (Adminer - untuk lihat database)
- Kalau muncul login page, berarti database sudah jalan ✅

### 2. Jalankan Aplikasi
Sekarang jalankan aplikasi Go:

```bash
make dev
```

**Apa yang terjadi?**
- Aplikasi Go akan compile dan jalan di port 8080
- Akan konek ke database PostgreSQL
- Siap menerima HTTP request

**Cek apakah berhasil:**
- Buka http://localhost:8080/health
- Kalau muncul `{"status":"healthy"}`, berarti aplikasi jalan ✅

### 3. Tes API Dasar
Pastikan API bisa diakses:

```bash
# Test via browser atau curl
curl http://localhost:8080/api/v1/public/ping
```

**Expected response:**
```json
{
  "success": true,
  "message": "pong",
  "data": null
}
```

## 🔄 Loop Development Harian

### Langkah 1: Edit Code
- Buat perubahan di file Go
- Edit handler, service, model, atau migration
- Simpan file

### Langkah 2: Restart Aplikasi
Karena tidak pakai hot reload, restart manual:

```bash
# Stop aplikasi (Ctrl+C di terminal yang jalan make dev)
# Lalu jalankan lagi
make dev
```

**💡 Tips:** Biarkan terminal database tetap jalan, hanya restart aplikasi saja.

### Langkah 3: Test Perubahan
Test perubahan yang dibuat:

```bash
# Test manual via browser/curl
curl http://localhost:8080/your-new-endpoint

# Atau jalankan automated tests
make test
```

### Langkah 4: Commit (Opsional)
Kalau perubahan sudah oke:

```bash
git add .
git commit -m "feat: add new feature"
```

## 🧪 Testing Cepat

### Run Tests
```bash
# Test semua
make test

# Test dengan coverage
make test-coverage
```

**Kapan harus test?**
- Setelah bikin fitur baru
- Sebelum commit code
- Kalau ada bug aneh

### Check Code Quality
```bash
# Format code
make fmt

# Check linting
make lint

# Run semua quality checks
make verify
```

## 🗄️ Database Tasks

### Lihat Database
```bash
# Buka Adminer di browser
open http://localhost:8081

# Atau connect via terminal
make db-connect
```

**Login Adminer:**
- Server: `postgres`
- Username: `postgres` 
- Password: `dev_password`
- Database: `golang_template_dev`

### Migration Sederhana
```bash
# Buat migration baru
make migrate-create NAME=add_something

# Jalankan migration
make migrate-up

# Cek status migration
make migrate-status
```

## 🌙 Mengakhiri Development

### 1. Stop Aplikasi
```bash
# Di terminal yang jalan make dev, tekan:
Ctrl+C
```

### 2. Stop Database (Opsional)
```bash
# Kalau mau matikan database juga
make db-down
```

**💡 Tips:** Database bisa tetap jalan, tidak masalah. Data akan tetap tersimpan.

## ⚡ Quick Commands Reference

| Command | Fungsi | Kapan Dipakai |
|---------|--------|---------------|
| `make db-up` | Nyalakan database | Awal development |
| `make dev` | Jalankan aplikasi | Setiap restart app |
| `make test` | Run tests | Setelah coding |
| `make fmt` | Format code | Sebelum commit |
| `make db-down` | Matikan database | Akhir hari (opsional) |

## 🚨 Quick Troubleshooting

### Aplikasi tidak mau start
```bash
# Cek apakah port 8080 sudah dipakai
lsof -i :8080

# Kalau ada, kill processnya
kill -9 <PID>
```

### Database connection error
```bash
# Cek apakah database jalan
make db-logs

# Restart database
make db-down && make db-up
```

### Tests gagal
```bash
# Lihat error message detail
make test

# Kalau database issue, reset migration
make migrate-down && make migrate-up
```

## 📚 Next Steps

- **Perlu setup database?** → [Database Development Guide](database-development.md)
- **Mau deploy ke production?** → [Database Production Guide](database-production.md)
- **Mau belajar testing?** → [Testing Guide](testing-guide.md)
- **Ada masalah?** → [Troubleshooting Guide](troubleshooting.md)

---

**💡 Pro Tips:**
- Simpan bookmark untuk http://localhost:8080/health dan http://localhost:8081
- Buat alias di shell: `alias devstart="make db-up && make dev"`
- Gunakan multiple terminal: satu untuk database, satu untuk aplikasi