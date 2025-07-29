# Development Guide

Selamat datang di Go Backend Template! Panduan ini akan membantu kamu memulai development dengan cepat dan efisien.

## 🎯 Apa yang Mau Kamu Lakukan Hari Ini?

**Pilih aktivitas yang paling sesuai dengan kebutuhan kamu:**

### 🌅 **Mulai Development Harian**
- ✅ Start database dan aplikasi
- ✅ Workflow edit-test-commit  
- ✅ Commands yang sering dipakai
- ✅ Quick troubleshooting

**👉 Baca:** [Daily Workflow Guide](daily-workflow.md)

### 🗄️ **Kerja dengan Database Local**
- ✅ Setup PostgreSQL dengan Docker
- ✅ Buat dan jalankan migration
- ✅ Kelola data development
- ✅ Database management tools

**👉 Baca:** [Database Development Guide](database-development.md)

### 🏭 **Deploy Database ke Production**
- ✅ Production migration yang aman
- ✅ Environment configuration
- ✅ Safety checks dan backup
- ✅ Emergency procedures

**👉 Baca:** [Database Production Guide](database-production.md)

### 🧪 **Testing & Quality Assurance**
- ✅ Unit, integration, E2E testing
- ✅ Test best practices
- ✅ Coverage dan automation
- ✅ Mock dan fixtures

**👉 Baca:** [Testing Guide](testing-guide.md)

### ✨ **Code Quality & Standards**
- ✅ Linting, formatting, vetting
- ✅ Pre-commit hooks
- ✅ Quality metrics
- ✅ Editor integration

**👉 Baca:** [Code Quality Guide](code-quality.md)

### 🚨 **Troubleshooting Masalah**
- ✅ Aplikasi tidak start
- ✅ Database connection issues
- ✅ Port conflicts
- ✅ Test failures

**👉 Baca:** [Troubleshooting Guide](troubleshooting.md)

## 🚀 Quick Start (Paling Sering Dibutuhkan)

### Hari Pertama Setup
```bash
# 1. Clone project
git clone <your-repo>
cd golang_template

# 2. Install dependencies
go mod tidy

# 3. Start database
make db-up

# 4. Run migrations
make migrate-up

# 5. Start application
make dev

# 6. Test API
curl http://localhost:8080/health
```

### Workflow Harian
```bash
# Start development session
make db-up && make dev

# Run tests setelah coding
make test

# Check code quality
make verify

# End session (opsional)
make db-down
```

## 📚 Learning Path untuk Pemula

### Level 1: Basics (Hari 1-3)
1. **Setup environment** → [Daily Workflow Guide](daily-workflow.md)
2. **Understand database** → [Database Development Guide](database-development.md)
3. **Basic troubleshooting** → [Troubleshooting Guide](troubleshooting.md)

### Level 2: Development (Minggu 1-2)  
1. **Write tests** → [Testing Guide](testing-guide.md)
2. **Code quality** → [Code Quality Guide](code-quality.md)
3. **Advanced database** → [Database Development Guide](database-development.md)

### Level 3: Production (Minggu 3+)
1. **Production deployment** → [Database Production Guide](database-production.md)
2. **Performance tuning** → [Advanced guides](../architecture/)
3. **Monitoring & ops** → [Operations guides](../deployment/)

## 🎨 Development Workflow Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    Daily Development Flow                    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  🌅 Start Day          🔄 Development Loop    🌙 End Day   │
│  ├─ make db-up         ├─ Edit code           ├─ Commit    │
│  ├─ make dev           ├─ Restart app         ├─ Push      │  
│  └─ Check health       ├─ Test manually       └─ Clean up  │
│                        ├─ make test                        │
│                        └─ Fix issues                       │
└─────────────────────────────────────────────────────────────┘
```

## 🛠️ Essential Commands Reference

### Database
| Command | Purpose | When to Use |
|---------|---------|-------------|
| `make db-up` | Start database | Awal development |
| `make db-down` | Stop database | Akhir hari (opsional) |
| `make migrate-up` | Run migrations | After pull, new migrations |
| `make migrate-status` | Check migration status | Debug migration issues |

### Application  
| Command | Purpose | When to Use |
|---------|---------|-------------|
| `make dev` | Start application | Development |
| `make build` | Build binary | Before deployment |
| `make test` | Run tests | After coding |
| `make verify` | Code quality checks | Before commit |

### Troubleshooting
| Command | Purpose | When to Use |
|---------|---------|-------------|
| `lsof -i :8080` | Check port usage | App won't start |
| `make db-logs` | Database logs | Database issues |
| `docker-compose ps` | Container status | Services not running |

## 📖 Documentation Structure

```
docs/getting-started/
├── 🏠 development.md          # This overview (START HERE)
├── 🌅 daily-workflow.md       # Daily development workflow  
├── 🗄️ database-development.md # Local database & migrations
├── 🏭 database-production.md  # Production deployment
├── 🧪 testing-guide.md        # Testing practices
├── ✨ code-quality.md         # Code quality tools
└── 🚨 troubleshooting.md      # Problem solving
```

## 🎯 Quick Decision Tree

**Stuck? Follow this decision tree:**

```
❓ What's your situation?

├─ 🆕 First time setup
│  └─ 👉 Read: Daily Workflow Guide
│
├─ 🐛 Something is broken  
│  └─ 👉 Read: Troubleshooting Guide
│
├─ 🗄️ Need to work with database
│  ├─ Local development → Database Development Guide
│  └─ Production deployment → Database Production Guide
│
├─ 🧪 Want to write/run tests
│  └─ 👉 Read: Testing Guide
│
├─ ✨ Code quality issues
│  └─ 👉 Read: Code Quality Guide
│
└─ 📚 Want to learn more
   └─ 👉 Check: Architecture guides in ../architecture/
```

## 🆘 Need Help?

### Common Scenarios & Solutions

| Problem | Quick Fix | Detailed Guide |
|---------|-----------|----------------|
| **"App won't start"** | `lsof -i :8080` then `kill -9 <PID>` | [Troubleshooting](troubleshooting.md#aplikasi-tidak-bisa-start) |
| **"Database error"** | `make db-down && make db-up` | [Troubleshooting](troubleshooting.md#database-issues) |
| **"Tests failing"** | `make db-clean && make db-up` | [Testing Guide](testing-guide.md) |
| **"Migration stuck"** | `make migrate-status` | [Database Development](database-development.md) |
| **"Code quality fails"** | `make verify` | [Code Quality Guide](code-quality.md) |

### Emergency Reset (Development Only)
```bash
# Nuclear option - resets everything
make db-clean
docker system prune -f
go clean -cache -modcache
go mod tidy
make db-up
make migrate-up
make dev
```

## 🔗 Advanced Topics

Setelah nyaman dengan basics, explore topics advanced:

- **Architecture**: [Architecture Overview](../architecture/overview.md)
- **API Development**: [Creating APIs](../api-development/creating-apis.md)  
- **Configuration**: [Configuration Guide](../configuration/overview.md)
- **Deployment**: [Deployment Guide](../deployment/)
- **Performance**: [Performance Guide](../optimization/)

## 📞 Getting Support

1. **Check documentation** relevant di atas
2. **Try troubleshooting steps** 
3. **Search issues** di repository
4. **Ask team members** dengan specific error messages
5. **Create issue** dengan complete information

---

**💡 Pro Tip:** Bookmark halaman ini sebagai starting point. Setiap kali butuh bantuan development, mulai dari sini dan follow links ke dokumentasi yang specific.

**🎯 Remember:** Dokumentasi ini dibuat untuk memudahkan hidup developer. Kalau masih bingung atau ada yang kurang jelas, let us know!