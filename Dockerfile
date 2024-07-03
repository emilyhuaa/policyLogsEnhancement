# Use a lightweight base image
FROM golang:1.19-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o main .

# Use a minimal base image for the final image
FROM alpine:latest

# Create a non-root user
RUN adduser -D appuser

# Copy the built binary from the builder stage
COPY --from=builder /app/main /app/main

# Set the working directory
WORKDIR /app

# Switch to the non-root user
USER appuser

# Expose the port your application listens on (if applicable)
# EXPOSE 8080

# Set the entry point to run the application
ENTRYPOINT ["./main"]
