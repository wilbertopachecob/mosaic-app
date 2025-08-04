# 🚀 Build & Deployment Scripts

This directory contains comprehensive build and deployment scripts for the Mosaic App.

## 📁 Scripts Overview

### 🔧 Development Scripts

#### `dev.sh` - Development Environment
Sets up and manages the development environment with hot reloading.

```bash
# Setup development environment
./scripts/dev.sh setup

# Start development servers
./scripts/dev.sh start

# Run tests
./scripts/dev.sh test

# Run linting
./scripts/dev.sh lint

# Type checking
./scripts/dev.sh type-check

# Clean development artifacts
./scripts/dev.sh clean

# Show development status
./scripts/dev.sh status

# Open development URLs
./scripts/dev.sh open
```

### 🏗️ Build Scripts

#### `build.sh` - Build System
Handles building both frontend and backend components.

```bash
# Build everything (default)
./scripts/build.sh build

# Build only backend
./scripts/build.sh backend

# Build only frontend
./scripts/build.sh frontend

# Clean build artifacts
./scripts/build.sh clean

# Run all tests
./scripts/build.sh test
```

### 🚀 Deployment Scripts

#### `deploy.sh` - Deployment System
Handles deployment to different environments.

```bash
# Deploy locally
./scripts/deploy.sh local

# Deploy with Docker
./scripts/deploy.sh docker

# Deploy with Docker Compose
./scripts/deploy.sh compose

# Deploy to production (systemd)
./scripts/deploy.sh production

# Create cloud deployment package
./scripts/deploy.sh cloud

# Show deployment status
./scripts/deploy.sh status
```

### 🔄 CI/CD Scripts

#### `ci.sh` - Continuous Integration
Designed for CI/CD pipelines and automated testing.

```bash
# Run complete CI pipeline
./scripts/ci.sh full

# Run only tests
./scripts/ci.sh test

# Run build pipeline
./scripts/ci.sh build

# Run security checks
./scripts/ci.sh security

# Cleanup artifacts
./scripts/ci.sh cleanup
```

## 🛠️ Makefile Commands

For convenience, a Makefile provides shortcuts for common tasks:

```bash
# Development
make dev-setup      # Setup development environment
make dev-start      # Start development servers
make dev-test       # Run tests
make dev-lint       # Run linting
make dev-clean      # Clean development artifacts

# Build
make build          # Build everything
make build-backend  # Build only backend
make build-frontend # Build only frontend
make clean          # Clean build artifacts

# Deployment
make deploy-local   # Deploy locally
make deploy-docker  # Deploy with Docker
make deploy-compose # Deploy with Docker Compose
make deploy-prod    # Deploy to production
make deploy-cloud   # Create cloud package

# Utilities
make status         # Show deployment status
make logs           # Show application logs
make stop           # Stop all services
```

## 🎯 Quick Start Workflows

### Development Workflow
```bash
# First time setup
make dev-setup

# Start development
make dev-start

# Run tests
make dev-test

# Clean up
make dev-clean
```

### Production Deployment Workflow
```bash
# Build and deploy locally
make deploy

# Build and deploy with Docker
make docker

# Build and deploy to production
make production
```

### CI/CD Workflow
```bash
# Run complete CI pipeline
./scripts/ci.sh full

# Run only tests
./scripts/ci.sh test

# Run security checks
./scripts/ci.sh security
```

## 🌍 Deployment Options

### Local Deployment
```bash
./scripts/deploy.sh local
```
- Runs the application directly on your machine
- Good for development and testing
- Requires Go and Node.js installed

### Docker Deployment
```bash
./scripts/deploy.sh docker
```
- Containerized deployment
- Isolated environment
- Easy to manage and scale
- Requires Docker installed

### Docker Compose Deployment
```bash
./scripts/deploy.sh compose
```
- Multi-service deployment
- Easy configuration management
- Good for complex setups
- Requires Docker and Docker Compose

### Production Deployment
```bash
./scripts/deploy.sh production
```
- Systemd service installation
- Automatic startup on boot
- Production-ready configuration
- Requires sudo privileges

### Cloud Deployment
```bash
./scripts/deploy.sh cloud
```
- Creates deployment package
- Supports multiple cloud platforms
- Includes platform-specific instructions
- Ready for AWS, GCP, Azure, etc.

## 📋 Prerequisites

### For All Scripts
- Bash shell
- Go 1.22+
- Node.js 18+
- npm

### For Docker Deployment
- Docker
- Docker Compose (for compose deployment)

### For Production Deployment
- sudo privileges
- systemd (Linux)

### For Development
- Air (auto-installed by dev script)
- golangci-lint (optional, for linting)

## 🔧 Configuration

### Environment Variables
Scripts use environment variables for configuration:

```bash
# Server Configuration
SERVER_PORT=8080

# File Upload Settings
MAX_FILE_SIZE=10485760

# Tiles Configuration
TILES_DIR=tiles

# Logging
LOG_LEVEL=info
```

### Configuration Files
- `.env` - Environment variables
- `.air.toml` - Go hot reloading configuration
- `docker-compose.yml` - Docker Compose configuration
- `Dockerfile` - Docker container configuration

## 📊 Monitoring & Logs

### View Application Status
```bash
./scripts/deploy.sh status
make status
```

### View Application Logs
```bash
make logs
```

### Stop All Services
```bash
make stop
```

## 🧪 Testing

### Run All Tests
```bash
./scripts/build.sh test
make test
```

### Run Tests with Coverage
```bash
./scripts/ci.sh test
make test-coverage
```

### Run Integration Tests
```bash
./scripts/ci.sh full
```

## 🔒 Security

### Security Audit
```bash
./scripts/ci.sh security
make security
```

### Dependency Scanning
```bash
# Go dependencies
go list -json -deps . | nancy sleuth

# Frontend dependencies
cd frontend && npm audit
```

## 📈 Performance

### Run Benchmarks
```bash
make benchmark
```

### Bundle Analysis
```bash
cd frontend && npm run build -- --analyze
```

## 🚨 Troubleshooting

### Common Issues

#### Script Permission Denied
```bash
chmod +x scripts/*.sh
```

#### Docker Not Running
```bash
# Start Docker
sudo systemctl start docker  # Linux
open -a Docker              # macOS
```

#### Port Already in Use
```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>
```

#### Build Failures
```bash
# Clean and rebuild
make clean
make build
```

### Getting Help

#### Script Help
```bash
./scripts/build.sh help
./scripts/deploy.sh help
./scripts/dev.sh help
./scripts/ci.sh help
```

#### Makefile Help
```bash
make help
make help-<command>
```

## 📝 Examples

### Complete Development Setup
```bash
# Clone repository
git clone <repository-url>
cd mosaic-app

# Setup development environment
make dev-setup

# Start development
make dev-start

# Open in browser
make open
```

### Production Deployment
```bash
# Build application
make build

# Deploy to production
sudo make deploy-prod

# Check status
make status

# View logs
make logs
```

### Docker Deployment
```bash
# Build and deploy with Docker
make docker

# Check container status
docker ps

# View container logs
docker logs mosaic-app
```

### CI/CD Pipeline
```bash
# Run complete pipeline
./scripts/ci.sh full

# Check artifacts
ls -la artifacts/

# Deploy from artifacts
cd artifacts/deploy
./deploy.sh
```

## 🤝 Contributing

When adding new scripts or modifying existing ones:

1. Follow the existing naming conventions
2. Add proper error handling
3. Include colored output for better UX
4. Add help documentation
5. Update this README
6. Test on multiple platforms

## 📄 License

These scripts are part of the Mosaic App project and follow the same license terms. 