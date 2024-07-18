# Use a lightweight base image
FROM golang:1.22

ENV GOPROXY=direct

# Set the working directory
WORKDIR /app

# Copy the source code
COPY . .

RUN go mod download

# Build the Go application
RUN go build -o main .

CMD ["./main"]
