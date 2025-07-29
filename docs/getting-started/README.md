# Getting Started with Go Backend Template

Selamat datang! Panduan ini akan membantu kamu memulai development dengan Go Backend Template. Semua dokumentasi sudah dipecah berdasarkan kebutuhan specific supaya mudah dicari dan dipahami.

## 🎯 Mau Ngapain Hari Ini?

**Langsung pilih sesuai kebutuhan kamu:**

### 🏠 **Baru Pertama Kali Setup?**
**👉 Start here:** [Development Guide](development.md) - Overview lengkap dengan navigation

### 🌅 **Mulai Development Harian**
**👉 Baca:** [Daily Workflow Guide](daily-workflow.md)
- Workflow start-to-finish setiap hari
- Commands yang sering dipakai
- Quick start untuk pemula

### 🗄️ **Kerja dengan Database**

#### Local Development
**👉 Baca:** [Database Development Guide](database-development.md)  
- Setup PostgreSQL lokal dengan Docker
- Migration untuk pemula (dengan analogi mudah)
- Database management tools

#### Production Deployment  
**👉 Baca:** [Database Production Guide](database-production.md)
- Production migration yang aman
- Environment configuration
- Safety checklist & backup procedures

### 🧪 **Testing & Quality**
**👉 Baca:** [Testing Guide](testing-guide.md)
- Unit, integration, E2E testing untuk pemula
- Test best practices dengan contoh code
- Coverage dan automation

**👉 Baca:** [Code Quality Guide](code-quality.md)
- Linting, formatting, vetting tools
- Pre-commit hooks setup
- Quality metrics & CI integration

### 🚨 **Ada Masalah?**
**👉 Baca:** [Troubleshooting Guide](troubleshooting.md)
- Quick problem finder dengan solusi
- Common issues & fixes
- Emergency commands

## 🚀 Super Quick Start

**Kalau terburu-buru, ini minimum commands:**

```bash
# Setup pertama kali
go mod tidy && make db-up && make migrate-up && make dev

# Workflow harian
make db-up && make dev

# Test setelah coding
make test && make verify
```

## 📚 Learning Path

### 👶 **Pemula (Hari 1-3)**
1. [Development Guide](development.md) - Overview
2. [Daily Workflow Guide](daily-workflow.md) - Basic workflow
3. [Database Development Guide](database-development.md) - Database basics
4. [Troubleshooting Guide](troubleshooting.md) - When stuck

### 🧑‍💻 **Intermediate (Minggu 1-2)**
1. [Testing Guide](testing-guide.md) - Write good tests
2. [Code Quality Guide](code-quality.md) - Quality tools
3. [Database Production Guide](database-production.md) - Production ready

### 🚀 **Advanced (Minggu 3+)**
- [Architecture guides](../architecture/) - Deep dive architecture
- [API Development](../api-development/) - Build APIs
- [Configuration](../configuration/) - Advanced config
- [Deployment](../deployment/) - Production deployment

## 🎯 Quick Navigation

**Pilih berdasarkan situasi kamu:**

| Situasi | Dokumen |
|---------|---------|
| 🆕 **First time here** | [Development Guide](development.md) |
| 🌅 **Daily development** | [Daily Workflow](daily-workflow.md) |
| 🗄️ **Database local** | [Database Development](database-development.md) |
| 🏭 **Database production** | [Database Production](database-production.md) |
| 🧪 **Testing** | [Testing Guide](testing-guide.md) |
| ✨ **Code quality** | [Code Quality Guide](code-quality.md) |
| 🚨 **Something broken** | [Troubleshooting Guide](troubleshooting.md) |

## 🎨 Visual Workflow

```
┌───────────────────────────────────────────────────────────────┐
│                     Development Journey                       │
├───────────────────────────────────────────────────────────────┤
│                                                               │
│  🆕 New Developer                                             │
│  │                                                            │
│  ├─ 🏠 Read: Development Guide (Overview)                     │
│  │                                                            │
│  ├─ 🌅 Read: Daily Workflow (Basic commands)                 │
│  │                                                            │
│  ├─ 🗄️ Read: Database Development (Local setup)              │
│  │                                                            │
│  ├─ 🧪 Read: Testing Guide (Quality practices)               │
│  │                                                            │
│  ├─ ✨ Read: Code Quality (Professional standards)           │
│  │                                                            │
│  └─ 🏭 Read: Database Production (Deploy safely)             │
│                                                               │
│  🎓 Experienced Developer                                     │
│  └─ 📚 Advanced Topics (Architecture, Performance, etc.)     │
│                                                               │
└───────────────────────────────────────────────────────────────┘
```

## 🛠️ Essential Commands Cheat Sheet

### Daily Development
```bash
make db-up          # Start database
make dev            # Start application  
make test           # Run tests
make verify         # Code quality checks
```

### Database
```bash
make migrate-up     # Run migrations
make migrate-status # Check migration status
make db-connect     # Connect to database
make db-logs        # View database logs
```

### Troubleshooting
```bash
lsof -i :8080              # Check port 8080
docker-compose ps          # Check containers
make db-down && make db-up # Restart database
```

## 🆘 Quick Help

### Common Problems & Quick Fixes

| Problem | Quick Solution | Full Guide |
|---------|----------------|------------|
| **App won't start** | `lsof -i :8080` → `kill -9 <PID>` | [Troubleshooting](troubleshooting.md) |
| **Database error** | `make db-down && make db-up` | [Database Dev](database-development.md) |
| **Tests failing** | `make db-clean && make db-up` | [Testing Guide](testing-guide.md) |
| **First time setup** | `go mod tidy && make db-up` | [Development Guide](development.md) |

### Emergency Reset (Development Only)
```bash
# Nuclear option - resets everything
make db-clean && docker system prune -f && go clean -cache
go mod tidy && make db-up && make migrate-up && make dev
```

## 📋 Documentation Standards

Semua dokumentasi di folder ini mengikuti prinsip:

- ✅ **Bahasa sederhana** untuk orang awam
- ✅ **Step-by-step** dengan penjelasan
- ✅ **Real examples** dengan expected output  
- ✅ **Common mistakes** dan cara mengatasinya
- ✅ **Cross-references** antar dokumen
- ✅ **Quick decision trees** untuk situasi specific

## 🔗 Other Documentation

Setelah comfortable dengan getting-started guides, explore:

- **[Architecture](../architecture/)** - How the system is structured
- **[API Development](../api-development/)** - Building REST APIs
- **[Configuration](../configuration/)** - Environment & config management
- **[Database](../database/)** - Advanced database topics
- **[Deployment](../deployment/)** - Production deployment
- **[Examples](../examples/)** - Practical code examples

## 💡 Pro Tips

1. **Bookmark [Development Guide](development.md)** - Starting point untuk semua
2. **Keep [Troubleshooting Guide](troubleshooting.md) handy** - Untuk emergency
3. **Follow the learning path** - Jangan skip basics
4. **Practice with real examples** - Theory + practice = mastery
5. **Ask questions** - Better ask than stuck

## 📞 Need More Help?

1. **Start with relevant guide** dari list di atas
2. **Check troubleshooting** kalau ada masalah
3. **Search existing issues** di repository
4. **Ask with specific error messages** dan steps to reproduce
5. **Contribute back** kalau kamu solve something

---

**🎯 Remember:** Dokumentasi ini dibuat untuk **memudahkan hidup developer**. Kalau ada yang kurang jelas atau missing, let us know supaya bisa diperbaiki!

**🚀 Ready to start?** Pilih guide yang sesuai dan happy coding! 🎉