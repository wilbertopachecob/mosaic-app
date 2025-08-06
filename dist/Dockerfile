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
