#!/bin/bash

# Mosaic App Development Script
# This script sets up the development environment with hot reloading

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

# Function to install development dependencies
install_dev_deps() {
    print_status "Installing development dependencies..."
    
    # Install Air for Go hot reloading
    if ! command_exists air; then
        print_status "Installing Air for Go hot reloading..."
        go install github.com/cosmtrek/air@latest
    fi
    
    # Install frontend dependencies
    if [ ! -d "frontend/node_modules" ]; then
        print_status "Installing frontend dependencies..."
        cd frontend
        npm install
        cd ..
    fi
    
    print_success "Development dependencies installed"
}

# Function to create development configuration
setup_dev_config() {
    print_status "Setting up development configuration..."
    
    # Create .env for development
    if [ ! -f ".env" ]; then
        print_status "Creating development .env file..."
        cat > .env << EOF
# Development Configuration
SERVER_PORT=8080
MAX_FILE_SIZE=10485760
TILES_DIR=tiles
LOG_LEVEL=debug
EOF
    fi
    
    # Create Air configuration for Go hot reloading
    if [ ! -f ".air.toml" ]; then
        print_status "Creating Air configuration..."
        cat > .air.toml << EOF
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ."
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "frontend", "dist", "scripts"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
EOF
    fi
    
    # Create tiles directory
    mkdir -p tiles
    
    print_success "Development configuration created"
}

# Function to start development environment
start_dev() {
    print_status "Starting development environment..."
    
    # Check if Air is installed
    if ! command_exists air; then
        print_error "Air is not installed. Run: ./scripts/dev.sh setup"
        exit 1
    fi
    
    # Check if .env exists
    if [ ! -f ".env" ]; then
        print_warning "No .env file found. Creating development config..."
        setup_dev_config
    fi
    
    # Start frontend development server
    print_status "Starting frontend development server..."
    cd frontend
    npm start &
    FRONTEND_PID=$!
    cd ..
    
    # Wait a moment for frontend to start
    sleep 3
    
    # Start backend with hot reloading
    print_status "Starting backend with hot reloading..."
    print_status "Frontend: http://localhost:3000"
    print_status "Backend:  http://localhost:8080"
    print_status "Press Ctrl+C to stop all services"
    
    # Trap to kill frontend when script exits
    trap 'kill $FRONTEND_PID 2>/dev/null; exit' INT TERM
    
    # Start Air for backend hot reloading
    air
}

# Function to run tests in development
run_tests() {
    print_status "Running tests in development mode..."
    
    # Run backend tests
    print_status "Running backend tests..."
    go test -v ./...
    
    # Run frontend tests
    print_status "Running frontend tests..."
    cd frontend
    npm test -- --watchAll=false --coverage
    cd ..
    
    print_success "All tests completed"
}

# Function to run linting
run_lint() {
    print_status "Running linting..."
    
    # Go linting
    if command_exists golangci-lint; then
        print_status "Running Go linting..."
        golangci-lint run
    else
        print_warning "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    fi
    
    # Frontend linting
    print_status "Running frontend linting..."
    cd frontend
    npm run lint
    cd ..
    
    print_success "Linting completed"
}

# Function to run type checking
run_type_check() {
    print_status "Running type checking..."
    
    # Frontend type checking
    print_status "Running frontend type checking..."
    cd frontend
    npm run type-check
    cd ..
    
    print_success "Type checking completed"
}

# Function to clean development artifacts
clean_dev() {
    print_status "Cleaning development artifacts..."
    
    # Remove temporary files
    if [ -d "tmp" ]; then
        rm -rf tmp
    fi
    
    # Remove build artifacts
    if [ -f "mosaic" ]; then
        rm mosaic
    fi
    
    if [ -d "frontend/build" ]; then
        rm -rf frontend/build
    fi
    
    # Remove node_modules (optional)
    if [ "$1" = "full" ]; then
        if [ -d "frontend/node_modules" ]; then
            rm -rf frontend/node_modules
        fi
    fi
    
    print_success "Development artifacts cleaned"
}

# Function to show development status
show_dev_status() {
    print_status "Development environment status..."
    
    # Check if Air is installed
    if command_exists air; then
        print_success "Air (Go hot reloading) is installed"
    else
        print_warning "Air is not installed"
    fi
    
    # Check if frontend dependencies are installed
    if [ -d "frontend/node_modules" ]; then
        print_success "Frontend dependencies are installed"
    else
        print_warning "Frontend dependencies are not installed"
    fi
    
    # Check if .env exists
    if [ -f ".env" ]; then
        print_success "Development .env file exists"
    else
        print_warning "Development .env file does not exist"
    fi
    
    # Check if tiles directory exists
    if [ -d "tiles" ]; then
        print_success "Tiles directory exists"
    else
        print_warning "Tiles directory does not exist"
    fi
    
    # Check if services are running
    if pgrep -f "npm start" >/dev/null 2>&1; then
        print_success "Frontend development server is running"
    else
        print_warning "Frontend development server is not running"
    fi
    
    if pgrep -f "air" >/dev/null 2>&1; then
        print_success "Backend development server (Air) is running"
    else
        print_warning "Backend development server is not running"
    fi
}

# Function to open development URLs
open_dev_urls() {
    print_status "Opening development URLs..."
    
    # Open frontend
    if command_exists open; then
        open http://localhost:3000
    elif command_exists xdg-open; then
        xdg-open http://localhost:3000
    else
        print_status "Frontend: http://localhost:3000"
    fi
    
    # Open backend API
    if command_exists open; then
        open http://localhost:8080/api/health
    elif command_exists xdg-open; then
        xdg-open http://localhost:8080/api/health
    else
        print_status "Backend API: http://localhost:8080/api/health"
    fi
}

# Main development process
main() {
    echo "🔧 Mosaic App Development Script"
    echo "================================="
    
    # Parse command line arguments
    case "${1:-help}" in
        "setup")
            install_dev_deps
            setup_dev_config
            print_success "Development environment setup completed!"
            ;;
        "start")
            start_dev
            ;;
        "test")
            run_tests
            ;;
        "lint")
            run_lint
            ;;
        "type-check")
            run_type_check
            ;;
        "clean")
            clean_dev "$2"
            ;;
        "status")
            show_dev_status
            ;;
        "open")
            open_dev_urls
            ;;
        "help"|*)
            echo "Usage: $0 [setup|start|test|lint|type-check|clean|status|open]"
            echo ""
            echo "Development Commands:"
            echo "  setup      - Install dependencies and setup development config"
            echo "  start      - Start development environment with hot reloading"
            echo "  test       - Run all tests"
            echo "  lint       - Run linting checks"
            echo "  type-check - Run TypeScript type checking"
            echo "  clean      - Clean development artifacts"
            echo "  status     - Show development environment status"
            echo "  open       - Open development URLs in browser"
            echo ""
            echo "Examples:"
            echo "  ./scripts/dev.sh setup    # First time setup"
            echo "  ./scripts/dev.sh start    # Start development servers"
            echo "  ./scripts/dev.sh test     # Run tests"
            echo "  ./scripts/dev.sh clean    # Clean artifacts"
            echo "  ./scripts/dev.sh clean full  # Full clean including node_modules"
            exit 1
            ;;
    esac
}

# Run main function
main "$@" 