# 🎨 Mosaic Generator

A modern, AI-powered mosaic generator built with React, TypeScript, and Go. Transform your images into beautiful mosaics using intelligent tile matching algorithms.

## ✨ Features

- **AI-Powered Tile Matching**: Advanced color matching algorithm for optimal mosaic generation
- **Multiple Tile Sizes**: Choose from 5px to 100px tile sizes for different detail levels
- **Real-time Preview**: See your original image before generating the mosaic
- **High-Quality Output**: Generate high-resolution mosaics with configurable quality
- **Modern UI**: Beautiful, responsive interface built with Bootstrap and Font Awesome
- **Error Handling**: Comprehensive error handling and user feedback
- **Performance Optimized**: Efficient image processing with Go backend

## 🏗️ Architecture

### Backend (Go)
- **Gorilla Mux**: HTTP router with middleware support
- **Image Processing**: Custom algorithms for color matching and image resizing
- **Configuration Management**: Environment-based configuration
- **Logging**: Structured logging with Logrus
- **Testing**: Comprehensive unit tests with Testify

### Frontend (React + TypeScript)
- **Modern React**: Functional components with hooks
- **TypeScript**: Full type safety and better developer experience
- **Bootstrap 5**: Responsive UI components
- **Font Awesome**: Beautiful icons throughout the interface
- **Error Boundaries**: Graceful error handling

## 🚀 Quick Start

### Prerequisites

- Go 1.22+ 
- Node.js 18+ and npm
- Image files for tiles (place in `tiles/` directory)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/your-username/mosaic-app.git
   cd mosaic-app
   ```

2. **Set up the backend**
   ```bash
   # Install Go dependencies
   go mod tidy
   
   # Create tiles directory and add some images
   mkdir tiles
   # Add your tile images to the tiles/ directory
   ```

3. **Set up the frontend**
   ```bash
   cd frontend
   npm install
   npm run build
   cd ..
   ```

4. **Configure environment variables**
   ```bash
   # Create .env file
   cat > .env << EOF
   SERVER_PORT=8080
   MAX_FILE_SIZE=10485760
   TILES_DIR=tiles
   LOG_LEVEL=info
   EOF
   ```

5. **Run the application**
   ```bash
   go run .
   ```

6. **Open your browser**
   Navigate to `http://localhost:8080`

## 📁 Project Structure

```
mosaic-app/
├── config/                 # Configuration management
│   └── config.go
├── lib/                    # Core libraries
│   ├── img/               # Image processing
│   │   ├── img.go
│   │   └── img_test.go
│   └── tiles_db/          # Tiles database
│       ├── tiles_db.go
│       └── tiles_db_test.go
├── models/                 # Data models
│   └── mosaic.go
├── frontend/              # React frontend
│   ├── src/
│   │   ├── components/    # React components
│   │   ├── App.tsx        # Main app component
│   │   └── index.tsx      # Entry point
│   ├── public/            # Static assets
│   └── package.json
├── tiles/                 # Tile images (create this)
├── main.go               # Application entry point
├── handlers.go           # HTTP handlers
├── routes.go             # Route definitions
├── go.mod               # Go dependencies
└── README.md            # This file
```

## 🧪 Testing

### Backend Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./lib/img -v
```

### Frontend Tests
```bash
cd frontend

# Run tests
npm test

# Run tests with coverage
npm run test:coverage

# Run type checking
npm run type-check
```

## 🔧 Configuration

The application can be configured using environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8080` | HTTP server port |
| `MAX_FILE_SIZE` | `10485760` | Maximum file size (10MB) |
| `TILES_DIR` | `tiles` | Directory containing tile images |
| `LOG_LEVEL` | `info` | Logging level (debug, info, warn, error) |

## 📊 API Endpoints

### Health Check
```
GET /api/health
```
Returns service health status.

### Generate Mosaic
```
POST /api/file/upload
Content-Type: multipart/form-data
```
Generates a mosaic from uploaded image.

**Parameters:**
- `imgUpload`: Image file (max 10MB)
- `tileSize`: Tile size in pixels (5-200)

**Response:**
```json
{
  "mosaicImg": "base64_encoded_image",
  "duration": 2.45
}
```

## 🎯 Usage

1. **Prepare Tiles**: Add small images to the `tiles/` directory
2. **Upload Image**: Select an image file (JPG, PNG, GIF, etc.)
3. **Choose Tile Size**: Select appropriate tile size for desired detail level
4. **Generate**: Click "Generate Mosaic" and wait for processing
5. **Download**: Download your generated mosaic

## 🛠️ Development

### Backend Development
```bash
# Run with hot reload (requires air)
go install github.com/cosmtrek/air@latest
air

# Run tests
go test ./...

# Build binary
go build -o mosaic-app .
```

### Frontend Development
```bash
cd frontend

# Start development server
npm start

# Build for production
npm run build

# Lint code
npm run lint
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Gorilla Mux](https://github.com/gorilla/mux) for HTTP routing
- [Logrus](https://github.com/sirupsen/logrus) for structured logging
- [Bootstrap](https://getbootstrap.com/) for UI components
- [Font Awesome](https://fontawesome.com/) for icons

## 📞 Support

If you encounter any issues or have questions:

1. Check the [Issues](https://github.com/your-username/mosaic-app/issues) page
2. Create a new issue with detailed information
3. Include error messages and steps to reproduce

---

**Happy Mosaic Creating! 🎨✨**