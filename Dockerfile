FROM golang:1.22.2-alpine

# Install SQLite3 and other dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY /src/go.mod /src/go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the source to the Working Directory inside the container
COPY /src/. ./

# Build the Go app
RUN go build -o main game.go

# Expose port 8080 to the outside world
EXPOSE 1334

# Command to run the executable
CMD ["./main"]
