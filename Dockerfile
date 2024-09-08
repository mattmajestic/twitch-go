# Use the official Golang image
FROM golang:1.18-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to the working directory
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod tidy
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the Go app using the main.go inside cmd/
RUN go build -o main ./cmd/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
