# Mosaic App Makefile
# Provides convenient shortcuts for common tasks

.PHONY: help build clean test dev deploy docker production

# Default target
help:
	@echo "🎨 Mosaic App - Available Commands"
	@echo "=================================="
	@echo ""
	@echo "Development:"
	@echo "  make dev-setup    - Setup development environment"
	@echo "  make dev-start    - Start development servers with hot reloading"
	@echo "  make dev-test     - Run all tests"
	@echo "  make dev-lint     - Run linting checks"
	@echo "  make dev-clean    - Clean development artifacts"
	@echo ""
	@echo "Build & Deploy:"
	@echo "  make build        - Build both frontend and backend"
	@echo "  make build-backend - Build only backend"
	@echo "  make build-frontend - Build only frontend"
	@echo "  make clean        - Clean all build artifacts"
	@echo ""
	@echo "Deployment:"
	@echo "  make deploy-local - Deploy and run locally"
	@echo "  make deploy-docker - Deploy using Docker"
	@echo "  make deploy-compose - Deploy using Docker Compose"
	@echo "  make deploy-prod  - Deploy to production (systemd)"
	@echo "  make deploy-cloud - Create cloud deployment package"
	@echo ""
	@echo "Utilities:"
	@echo "  make status       - Show deployment status"
	@echo "  make logs         - Show application logs"
	@echo "  make stop         - Stop all running services"

# Development commands
dev-setup:
	@echo "🔧 Setting up development environment..."
	@chmod +x scripts/dev.sh
	@./scripts/dev.sh setup

dev-start:
	@echo "🚀 Starting development environment..."
	@chmod +x scripts/dev.sh
	@./scripts/dev.sh start

dev-test:
	@echo "🧪 Running tests..."
	@chmod +x scripts/dev.sh
	@./scripts/dev.sh test

dev-lint:
	@echo "🔍 Running linting..."
	@chmod +x scripts/dev.sh
	@./scripts/dev.sh lint

dev-clean:
	@echo "🧹 Cleaning development artifacts..."
	@chmod +x scripts/dev.sh
	@./scripts/dev.sh clean

# Build commands
build:
	@echo "🏗️ Building application..."
	@chmod +x scripts/build.sh
	@./scripts/build.sh build

build-backend:
	@echo "🏗️ Building backend..."
	@chmod +x scripts/build.sh
	@./scripts/build.sh backend

build-frontend:
	@echo "🏗️ Building frontend..."
	@chmod +x scripts/build.sh
	@./scripts/build.sh frontend

clean:
	@echo "🧹 Cleaning build artifacts..."
	@chmod +x scripts/build.sh
	@./scripts/build.sh clean

# Deployment commands
deploy-local:
	@echo "🚀 Deploying locally..."
	@chmod +x scripts/deploy.sh
	@./scripts/deploy.sh local

deploy-docker:
	@echo "🐳 Deploying with Docker..."
	@chmod +x scripts/deploy.sh
	@./scripts/deploy.sh docker

deploy-compose:
	@echo "🐳 Deploying with Docker Compose..."
	@chmod +x scripts/deploy.sh
	@./scripts/deploy.sh compose

deploy-prod:
	@echo "🏭 Deploying to production..."
	@chmod +x scripts/deploy.sh
	@./scripts/deploy.sh production

deploy-cloud:
	@echo "☁️ Creating cloud deployment package..."
	@chmod +x scripts/deploy.sh
	@./scripts/deploy.sh cloud

# Utility commands
status:
	@echo "📊 Checking deployment status..."
	@chmod +x scripts/deploy.sh
	@./scripts/deploy.sh status

logs:
	@echo "📋 Application logs:"
	@if pgrep -f "mosaic" >/dev/null 2>&1; then \
		echo "Local process logs:"; \
		ps aux | grep mosaic | grep -v grep; \
	fi
	@if command -v docker >/dev/null 2>&1 && docker ps -q -f name=mosaic-app >/dev/null 2>&1; then \
		echo "Docker container logs:"; \
		docker logs mosaic-app --tail 20; \
	fi
	@if command -v systemctl >/dev/null 2>&1 && systemctl is-active --quiet mosaic-app 2>/dev/null; then \
		echo "Systemd service logs:"; \
		sudo journalctl -u mosaic-app --no-pager -n 20; \
	fi

stop:
	@echo "🛑 Stopping all services..."
	@if pgrep -f "mosaic" >/dev/null 2>&1; then \
		echo "Stopping local process..."; \
		pkill -f mosaic; \
	fi
	@if command -v docker >/dev/null 2>&1 && docker ps -q -f name=mosaic-app >/dev/null 2>&1; then \
		echo "Stopping Docker container..."; \
		docker stop mosaic-app; \
	fi
	@if command -v docker-compose >/dev/null 2>&1; then \
		echo "Stopping Docker Compose services..."; \
		cd dist && docker-compose down 2>/dev/null || true; \
	fi
	@if command -v systemctl >/dev/null 2>&1 && systemctl is-active --quiet mosaic-app 2>/dev/null; then \
		echo "Stopping systemd service..."; \
		sudo systemctl stop mosaic-app; \
	fi
	@echo "✅ All services stopped"

# Quick development workflow
dev: dev-setup dev-start

# Quick build and deploy
deploy: build deploy-local

# Quick Docker deployment
docker: build deploy-docker

# Quick production deployment
production: build deploy-prod

# Install dependencies
install:
	@echo "📦 Installing dependencies..."
	@go mod tidy
	@cd frontend && npm install && cd ..

# Run tests
test:
	@echo "🧪 Running all tests..."
	@go test ./...
	@cd frontend && npm test -- --watchAll=false --passWithNoTests && cd ..

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -cover ./...
	@cd frontend && npm run test:coverage && cd ..

# Format code
format:
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@cd frontend && npm run lint:fix && cd ..

# Check code quality
check: format test lint

# Lint code
lint:
	@echo "🔍 Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi
	@cd frontend && npm run lint && cd ..

# Type check
type-check:
	@echo "🔍 Type checking..."
	@cd frontend && npm run type-check && cd ..

# Security audit
security:
	@echo "🔒 Security audit..."
	@go list -json -deps . | nancy sleuth
	@cd frontend && npm audit && cd ..

# Generate documentation
docs:
	@echo "📚 Generating documentation..."
	@if command -v godoc >/dev/null 2>&1; then \
		godoc -http=:6060 & \
		echo "Go documentation available at http://localhost:6060"; \
	else \
		echo "godoc not installed. Install with: go install golang.org/x/tools/cmd/godoc@latest"; \
	fi

# Performance benchmark
benchmark:
	@echo "⚡ Running benchmarks..."
	@go test -bench=. ./...

# Create release
release:
	@echo "🏷️ Creating release..."
	@if [ -z "$(VERSION)" ]; then \
		echo "Please specify VERSION=1.0.0"; \
		exit 1; \
	fi
	@git tag -a v$(VERSION) -m "Release v$(VERSION)"
	@git push origin v$(VERSION)
	@make build
	@mkdir -p releases
	@tar -czf releases/mosaic-app-v$(VERSION).tar.gz dist/
	@echo "Release v$(VERSION) created: releases/mosaic-app-v$(VERSION).tar.gz"

# Show system info
info:
	@echo "ℹ️ System Information:"
	@echo "Go version: $(shell go version)"
	@echo "Node version: $(shell node --version)"
	@echo "npm version: $(shell npm --version)"
	@echo "Docker version: $(shell docker --version 2>/dev/null || echo 'Not installed')"
	@echo "OS: $(shell uname -s) $(shell uname -m)"
	@echo "Architecture: $(shell go env GOOS)/$(shell go env GOARCH)"

# Show help for specific command
help-%:
	@echo "Help for command: $*"
	@case "$*" in \
		dev-setup) echo "Sets up the development environment with all dependencies and configuration files." ;; \
		dev-start) echo "Starts both frontend and backend development servers with hot reloading." ;; \
		build) echo "Builds both frontend and backend for production deployment." ;; \
		deploy-local) echo "Deploys and runs the application locally." ;; \
		deploy-docker) echo "Deploys the application using Docker." ;; \
		test) echo "Runs all tests for both frontend and backend." ;; \
		*) echo "No help available for command: $*" ;; \
	esac 