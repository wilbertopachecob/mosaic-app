#!/bin/bash

# Mosaic App Deployment Script
# This script handles deployment to different environments

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

# Function to check if Docker is available
check_docker() {
    if ! command_exists docker; then
        print_error "Docker is not installed. Please install Docker first."
        return 1
    fi
    
    if ! docker info >/dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker first."
        return 1
    fi
    
    return 0
}

# Function to deploy locally
deploy_local() {
    print_status "Deploying locally..."
    
    # Check if binary exists
    if [ ! -f "mosaic" ]; then
        print_error "Binary not found. Please run build script first: ./scripts/build.sh"
        exit 1
    fi
    
    # Check if frontend build exists
    if [ ! -d "frontend/build" ]; then
        print_error "Frontend build not found. Please run build script first: ./scripts/build.sh"
        exit 1
    fi
    
    # Create .env if it doesn't exist
    if [ ! -f ".env" ]; then
        print_warning "No .env file found. Creating from example..."
        cp env.example .env
    fi
    
    # Create tiles directory if it doesn't exist
    mkdir -p tiles
    
    print_success "Starting Mosaic App locally..."
    print_status "The app will be available at: http://localhost:8080"
    print_status "Press Ctrl+C to stop"
    
    # Start the application
    ./mosaic
}

# Function to deploy with Docker
deploy_docker() {
    print_status "Deploying with Docker..."
    
    if ! check_docker; then
        exit 1
    fi
    
    # Check if dist directory exists
    if [ ! -d "dist" ]; then
        print_error "Dist directory not found. Please run build script first: ./scripts/build.sh"
        exit 1
    fi
    
    cd dist
    
    # Build Docker image
    print_status "Building Docker image..."
    docker build -t mosaic-app .
    
    # Stop existing container if running
    if docker ps -q -f name=mosaic-app >/dev/null 2>&1; then
        print_status "Stopping existing container..."
        docker stop mosaic-app
        docker rm mosaic-app
    fi
    
    # Run container
    print_status "Starting Docker container..."
    docker run -d \
        --name mosaic-app \
        -p 8080:8080 \
        -v "$(pwd)/tiles:/app/tiles" \
        --restart unless-stopped \
        mosaic-app
    
    print_success "Docker deployment completed!"
    print_status "The app is available at: http://localhost:8080"
    print_status "Container logs: docker logs mosaic-app"
    print_status "Stop container: docker stop mosaic-app"
    
    cd ..
}

# Function to deploy with Docker Compose
deploy_docker_compose() {
    print_status "Deploying with Docker Compose..."
    
    if ! check_docker; then
        exit 1
    fi
    
    if ! command_exists docker-compose; then
        print_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
    
    # Check if dist directory exists
    if [ ! -d "dist" ]; then
        print_error "Dist directory not found. Please run build script first: ./scripts/build.sh"
        exit 1
    fi
    
    cd dist
    
    # Create .env file for docker-compose
    if [ ! -f ".env" ]; then
        print_warning "Creating .env file for Docker Compose..."
        cat > .env << EOF
SERVER_PORT=8080
MAX_FILE_SIZE=10485760
TILES_DIR=tiles
LOG_LEVEL=info
EOF
    fi
    
    # Start services
    print_status "Starting services with Docker Compose..."
    docker-compose up -d
    
    print_success "Docker Compose deployment completed!"
    print_status "The app is available at: http://localhost:8080"
    print_status "View logs: docker-compose logs -f"
    print_status "Stop services: docker-compose down"
    
    cd ..
}

# Function to deploy to production
deploy_production() {
    print_status "Deploying to production..."
    
    # Check if dist directory exists
    if [ ! -d "dist" ]; then
        print_error "Dist directory not found. Please run build script first: ./scripts/build.sh"
        exit 1
    fi
    
    # Create production directory
    PROD_DIR="/opt/mosaic-app"
    
    print_status "Installing to $PROD_DIR..."
    
    # Create production directory
    sudo mkdir -p "$PROD_DIR"
    
    # Copy files
    sudo cp -r dist/* "$PROD_DIR/"
    
    # Set permissions
    sudo chown -R $USER:$USER "$PROD_DIR"
    chmod +x "$PROD_DIR/mosaic"
    chmod +x "$PROD_DIR/start.sh"
    
    # Create systemd service
    print_status "Creating systemd service..."
    sudo tee /etc/systemd/system/mosaic-app.service > /dev/null << EOF
[Unit]
Description=Mosaic App
After=network.target

[Service]
Type=simple
User=$USER
WorkingDirectory=$PROD_DIR
ExecStart=$PROD_DIR/start.sh
Restart=always
RestartSec=5
Environment=SERVER_PORT=8080
Environment=MAX_FILE_SIZE=10485760
Environment=TILES_DIR=tiles
Environment=LOG_LEVEL=info

[Install]
WantedBy=multi-user.target
EOF
    
    # Reload systemd and enable service
    sudo systemctl daemon-reload
    sudo systemctl enable mosaic-app
    sudo systemctl start mosaic-app
    
    print_success "Production deployment completed!"
    print_status "Service status: sudo systemctl status mosaic-app"
    print_status "View logs: sudo journalctl -u mosaic-app -f"
    print_status "Stop service: sudo systemctl stop mosaic-app"
}

# Function to deploy to cloud (example for AWS EC2)
deploy_cloud() {
    print_status "Deploying to cloud (AWS EC2 example)..."
    
    # Check if AWS CLI is available
    if ! command_exists aws; then
        print_error "AWS CLI is not installed. Please install AWS CLI first."
        exit 1
    fi
    
    # Check if dist directory exists
    if [ ! -d "dist" ]; then
        print_error "Dist directory not found. Please run build script first: ./scripts/build.sh"
        exit 1
    fi
    
    # Create deployment package
    print_status "Creating deployment package..."
    cd dist
    tar -czf mosaic-app.tar.gz *
    cd ..
    
    print_status "Deployment package created: dist/mosaic-app.tar.gz"
    print_warning "Please upload this package to your cloud instance and follow the installation instructions."
    
    # Example commands for AWS EC2
    cat > dist/cloud-deploy-instructions.md << 'EOF'
# Cloud Deployment Instructions

## AWS EC2 Example

1. Upload the deployment package to your EC2 instance:
   ```bash
   scp -i your-key.pem dist/mosaic-app.tar.gz ec2-user@your-instance-ip:~/
   ```

2. SSH into your EC2 instance:
   ```bash
   ssh -i your-key.pem ec2-user@your-instance-ip
   ```

3. Install Docker (if using Docker deployment):
   ```bash
   sudo yum update -y
   sudo yum install -y docker
   sudo service docker start
   sudo usermod -a -G docker ec2-user
   ```

4. Extract and deploy:
   ```bash
   tar -xzf mosaic-app.tar.gz
   sudo docker build -t mosaic-app .
   sudo docker run -d -p 80:8080 --name mosaic-app mosaic-app
   ```

5. Configure security group to allow HTTP traffic on port 80.

## Other Cloud Platforms

- **Google Cloud Run**: Use the provided Dockerfile
- **Azure Container Instances**: Use the provided Dockerfile
- **Heroku**: Create a Procfile with: `web: ./mosaic`
- **DigitalOcean App Platform**: Use the provided Dockerfile
EOF
    
    print_success "Cloud deployment package created!"
    print_status "See dist/cloud-deploy-instructions.md for detailed instructions."
}

# Function to show deployment status
show_status() {
    print_status "Checking deployment status..."
    
    # Check if app is running locally
    if pgrep -f "mosaic" >/dev/null 2>&1; then
        print_success "Mosaic App is running locally (PID: $(pgrep -f mosaic))"
    else
        print_warning "Mosaic App is not running locally"
    fi
    
    # Check if Docker container is running
    if command_exists docker; then
        if docker ps -q -f name=mosaic-app >/dev/null 2>&1; then
            print_success "Docker container is running"
            docker ps -f name=mosaic-app
        else
            print_warning "Docker container is not running"
        fi
    fi
    
    # Check if systemd service is running
    if command_exists systemctl; then
        if systemctl is-active --quiet mosaic-app 2>/dev/null; then
            print_success "Systemd service is running"
        else
            print_warning "Systemd service is not running"
        fi
    fi
    
    # Check if port is listening
    if netstat -tuln 2>/dev/null | grep -q ":8080 "; then
        print_success "Port 8080 is listening"
    else
        print_warning "Port 8080 is not listening"
    fi
}

# Main deployment process
main() {
    echo "🚀 Mosaic App Deployment Script"
    echo "================================"
    
    # Parse command line arguments
    case "${1:-help}" in
        "local")
            deploy_local
            ;;
        "docker")
            deploy_docker
            ;;
        "compose")
            deploy_docker_compose
            ;;
        "production")
            deploy_production
            ;;
        "cloud")
            deploy_cloud
            ;;
        "status")
            show_status
            ;;
        "help"|*)
            echo "Usage: $0 [local|docker|compose|production|cloud|status]"
            echo ""
            echo "Deployment Options:"
            echo "  local      - Deploy and run locally"
            echo "  docker     - Deploy using Docker"
            echo "  compose    - Deploy using Docker Compose"
            echo "  production - Deploy to production (systemd service)"
            echo "  cloud      - Create cloud deployment package"
            echo "  status     - Show deployment status"
            echo ""
            echo "Prerequisites:"
            echo "  - Run build script first: ./scripts/build.sh"
            echo "  - For Docker: Install Docker and Docker Compose"
            echo "  - For production: Run with sudo privileges"
            echo "  - For cloud: Follow platform-specific instructions"
            exit 1
            ;;
    esac
}

# Run main function
main "$@" 