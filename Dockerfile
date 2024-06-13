# Use the official image to create a build environment
FROM golang:1.18

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . ./

# Build the Go app
RUN go build -v -o receipt-processor

# Expose port 3000 to the outside world
EXPOSE 3000

# Command to run the executable
CMD ["./receipt-processor"]
