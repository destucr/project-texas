# Use official Golang image as the builder stage
FROM golang:1.23.4-alpine3.21 as builder

# Set the working directory in the container
WORKDIR /app

# Copy go mod and sum files first to leverage Docker's cache
COPY go.mod go.sum ./

# Download Go dependencies (this will be cached if go.mod and go.sum don't change)
RUN go mod tidy

# Copy the entire project source code into the container
COPY . .

# Build the Go app (main.go is at the root)
RUN go build -o myapp main.go

# Use a lighter Golang image to run the app directly (instead of Debian)
FROM golang:1.23.4-alpine3.21

# Set the working directory for the app in the container
WORKDIR /app

# Copy the built Go app from the builder stage
COPY --from=builder /app/myapp .

# Expose the port that the app will run on (make sure your app listens on this port)
EXPOSE 8080

# Command to run the Go app
CMD ["./myapp"]
