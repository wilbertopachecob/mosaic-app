#!/bin/bash

# Mosaic App Build Script
# This script builds both the frontend and backend components

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    # Check Go
    if ! command_exists go; then
        print_error "Go is not installed. Please install Go 1.22+"
        exit 1
    fi
    
    # Check Node.js
    if ! command_exists node; then
        print_error "Node.js is not installed. Please install Node.js 18+"
        exit 1
    fi
    
    # Check npm
    if ! command_exists npm; then
        print_error "npm is not installed. Please install npm"
        exit 1
    fi
    
    print_success "All prerequisites are satisfied"
}

# Function to build backend
build_backend() {
    print_status "Building backend..."
    
    # Clean previous builds
    if [ -f "mosaic" ]; then
        rm mosaic
        print_status "Removed previous backend binary"
    fi
    
    # Update Go modules
    print_status "Updating Go modules..."
    go mod tidy
    
    # Run tests
    print_status "Running backend tests..."
    go test ./...
    
    # Build binary
    print_status "Compiling backend binary..."
    go build -ldflags="-s -w" -o mosaic .
    
    if [ -f "mosaic" ]; then
        print_success "Backend built successfully: mosaic"
        # Show binary info
        ls -lh mosaic
    else
        print_error "Backend build failed"
        exit 1
    fi
}

# Function to build frontend
build_frontend() {
    print_status "Building frontend..."
    
    cd frontend
    
    # Check if node_modules exists
    if [ ! -d "node_modules" ]; then
        print_status "Installing frontend dependencies..."
        npm install
    fi
    
    # Run tests
    print_status "Running frontend tests..."
    npm test -- --watchAll=false --passWithNoTests
    
    # Type check
    print_status "Running TypeScript type check..."
    npm run type-check
    
    # Build for production
    print_status "Building frontend for production..."
    npm run build
    
    # Check if build was successful
    if [ -d "build" ]; then
        print_success "Frontend built successfully"
        # Show build info
        du -sh build/
    else
        print_error "Frontend build failed"
        exit 1
    fi
    
    cd ..
}

# Function to create deployment package
create_deployment_package() {
    print_status "Creating deployment package..."
    
    # Create dist directory
    mkdir -p dist
    
    # Copy backend binary
    cp mosaic dist/
    
    # Copy frontend build
    cp -r frontend/build dist/
    
    # Copy configuration files
    cp env.example dist/
    cp README.md dist/
    
    # Create tiles directory if it doesn't exist
    mkdir -p dist/tiles
    
    # Create startup script
    cat > dist/start.sh << 'EOF'
#!/bin/bash
# Mosaic App Startup Script

# Load environment variables
if [ -f ".env" ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Set default port if not set
export SERVER_PORT=${SERVER_PORT:-8080}

echo "Starting Mosaic App on port $SERVER_PORT..."
./mosaic
EOF
    
    chmod +x dist/start.sh
    
    # Create Dockerfile
    cat > dist/Dockerfile << 'EOF'
FROM alpine:latest

# Install dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create app directory
WORKDIR /app

# Copy binary and assets
COPY mosaic .
COPY build/ ./build/
COPY start.sh .

# Make binary executable
RUN chmod +x mosaic start.sh

# Create tiles directory
RUN mkdir -p tiles

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

# Start the application
CMD ["./start.sh"]
EOF
    
    # Create docker-compose.yml
    cat > dist/docker-compose.yml << 'EOF'
version: '3.8'

services:
  mosaic-app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - SERVER_PORT=8080
      - MAX_FILE_SIZE=10485760
      - TILES_DIR=tiles
      - LOG_LEVEL=info
    volumes:
      - ./tiles:/app/tiles
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/api/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
EOF
    
    print_success "Deployment package created in dist/ directory"
    
    # Show package contents
    print_status "Deployment package contents:"
    ls -la dist/
}

# Function to clean build artifacts
clean_build() {
    print_status "Cleaning build artifacts..."
    
    # Remove backend binary
    if [ -f "mosaic" ]; then
        rm mosaic
    fi
    
    # Remove frontend build
    if [ -d "frontend/build" ]; then
        rm -rf frontend/build
    fi
    
    # Remove dist directory
    if [ -d "dist" ]; then
        rm -rf dist
    fi
    
    print_success "Build artifacts cleaned"
}

# Main build process
main() {
    echo "🎨 Mosaic App Build Script"
    echo "=========================="
    
    # Parse command line arguments
    case "${1:-build}" in
        "build")
            check_prerequisites
            build_backend
            build_frontend
            create_deployment_package
            print_success "Build completed successfully!"
            ;;
        "backend")
            check_prerequisites
            build_backend
            print_success "Backend build completed!"
            ;;
        "frontend")
            check_prerequisites
            build_frontend
            print_success "Frontend build completed!"
            ;;
        "clean")
            clean_build
            ;;
        "test")
            check_prerequisites
            print_status "Running all tests..."
            go test ./...
            cd frontend && npm test -- --watchAll=false --passWithNoTests && cd ..
            print_success "All tests passed!"
            ;;
        *)
            echo "Usage: $0 [build|backend|frontend|clean|test]"
            echo ""
            echo "Commands:"
            echo "  build    - Build both frontend and backend (default)"
            echo "  backend  - Build only the backend"
            echo "  frontend - Build only the frontend"
            echo "  clean    - Clean all build artifacts"
            echo "  test     - Run all tests"
            exit 1
            ;;
    esac
}

# Run main function
main "$@" 