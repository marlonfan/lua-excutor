# Stage 1: Build the frontend
FROM node:latest AS frontend-builder

WORKDIR /app

COPY . .

# Install dependencies
RUN npm install

# Build the frontend
RUN npm run build

# Stage 2: Build the backend
FROM golang:latest AS backend-builder

WORKDIR /app

# Copy the rest of the application code
COPY . .

# Download Go module dependencies
RUN go mod tidy

# Copy the build directory from the frontend build stage
COPY --from=frontend-builder /app/build ./build

# Build the Go application
RUN go build -o main .

# Stage 3: Run the application
FROM debian:latest

WORKDIR /root/

# Copy the built Go application from the backend build stage
COPY --from=backend-builder /app/main .

# Ensure the main file has execute permissions
RUN chmod +x ./main

# Set environment variables (if any)
ENV HTTP_AUTH_USERNAME=admin
ENV HTTP_AUTH_PASSWORD=admin
ENV HTTP_WEBHOOK_TOKEN=token

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]