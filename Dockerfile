# Use a lightweight base image
FROM public.ecr.aws/docker/library/golang:1.22.0

ENV GOPROXY=direct

# Set the working directory
WORKDIR /app

# Copy the source code
COPY . .

RUN go mod download

# Build the Go application
RUN go build -o main .

CMD ["./main"]
